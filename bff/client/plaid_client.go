package client

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jalexanderII/zero_fintech/bff/models"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/plaid/plaid-go/plaid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

var environmentSecret = map[string]string{
	"sandbox":     "PLAID_SECRET_SANDBOX",
	"development": "PLAID_SECRET_DEV",
}

type PlaidClient struct {
	// Name of the service
	Name string
	// Client is the object that contains all database functionalities
	Client       *plaid.PlaidApiService
	RedirectURL  string
	Products     []plaid.Products
	CountryCodes []plaid.CountryCode
	// custom logger
	l *logrus.Logger
	// Database collection
	PlaidDB mongo.Collection
	// Grpc client
	CoreClient core.CoreClient
	// to pass tokens through methods
	LinkToken   *models.Token
	PublicToken *models.Token
}

func NewPlaidClient(l *logrus.Logger, pdb mongo.Collection, coreClient core.CoreClient) *PlaidClient {
	// set constants from env
	plaidEnv := utils.GetEnv("PLAID_ENV")
	plaidSecret := utils.GetEnv(environmentSecret[plaidEnv])
	plaidClient := utils.GetEnv("PLAID_CLIENT_ID")

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", plaidClient)
	configuration.AddDefaultHeader("PLAID-SECRET", plaidSecret)
	configuration.UseEnvironment(environments[plaidEnv])

	countryCodes := convertCountryCodes(strings.Split(utils.GetEnv("PLAID_COUNTRY_CODES"), ","))
	products := convertProducts(strings.Split(utils.GetEnv("PLAID_PRODUCTS"), ","))
	client := plaid.NewAPIClient(configuration)
	return &PlaidClient{
		Name:         "ZeroFintech",
		Client:       client.PlaidApi,
		RedirectURL:  utils.GetEnv("PLAID_REDIRECT_URI"),
		Products:     products,
		CountryCodes: countryCodes,
		l:            l,
		PlaidDB:      pdb,
		CoreClient:   coreClient,
		LinkToken:    nil,
		PublicToken:  nil,
	}
}

// LinkTokenCreate creates a link token using the specified parameters
func (p *PlaidClient) LinkTokenCreate(ctx context.Context, username string) (*models.CreateLinkTokenResponse, error) {
	DbUser, err := p.GetUser(ctx, username, "")
	if err != nil {
		p.l.Error("[DB Error] error fetching user")
		return nil, err
	}
	id := DbUser.ID.Hex()

	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: id,
	}
	request := plaid.NewLinkTokenCreateRequest(p.Name, "en", p.CountryCodes, user)
	request.SetRedirectUri(p.RedirectURL)

	userid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		p.l.Errorf("[Error] error setting Hex from Id %+v", err)
		return nil, err
	}
	token, err := p.GetUserToken(ctx, &models.User{ID: userid, Username: username})
	if err == nil {
		// An Item's access_token does not change when using Link in update mode,
		// so there is no need to repeat the exchange token process.
		request.SetAccessToken(token.Value)
	} else {
		// Update mode: must not provide any products
		request.SetProducts(p.Products)
	}

	p.l.Infof("Link token request %+v", request)
	linkTokenCreateResp, _, err := p.Client.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		p.l.Errorf("[Plaid Error] error creating link token %+v", renderError(err)["error"])
		return nil, err
	}

	p.l.Info("link token created: ", linkTokenCreateResp)
	return &models.CreateLinkTokenResponse{Token: linkTokenCreateResp.GetLinkToken(), UserId: id}, nil
}

// ExchangePublicToken this function takes care of creating the permanent access token
// that will be stored in the database for cross-platform connection to users' bank.
// If for whatever reason there is a problem with the client or public token, there
// are json responses and logs that will adequately reflect all issues
func (p *PlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (*models.Token, error) {
	// exchange the public_token for an access_token
	exchangePublicTokenResp, _, err := p.Client.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()
	if err != nil {
		p.l.Errorf("[Plaid Error] error getting exchangePublicTokenResp %+v", renderError(err)["error"])
		return nil, err
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	itemID := exchangePublicTokenResp.GetItemId()

	p.l.Info("public token: " + publicToken)
	p.l.Info("access token: " + accessToken)
	p.l.Info("item ID: " + itemID)
	return &models.Token{Value: accessToken, ItemId: itemID}, nil
}

func (p *PlaidClient) GetAccountDetails(ctx context.Context, token *models.Token) (*core.AccountDetailsResponse, error) {
	liabilitiesReq := plaid.NewLiabilitiesGetRequest(token.Value)
	liabilitiesResp, _, err := p.Client.LiabilitiesGet(ctx).LiabilitiesGetRequest(*liabilitiesReq).Execute()
	if err != nil {
		p.l.Error("Error getting Liabilities", "error", err)
		return nil, err
	}
	liabilitiesResponse := models.LiabilitiesResponse{Liabilities: liabilitiesResp.GetLiabilities().Credit}

	const iso8601TimeFormat = "2006-01-02"
	endDate := time.Now().Local().Format(iso8601TimeFormat)
	numMonths := time.Duration(-30 * 12 * 24)
	startDate := time.Now().Local().Add(numMonths * time.Hour).Format(iso8601TimeFormat)

	transactionsResp, _, err := p.Client.TransactionsGet(ctx).TransactionsGetRequest(
		*plaid.NewTransactionsGetRequest(token.Value, startDate, endDate),
	).Execute()
	if err != nil {
		p.l.Error("Error getting Transactions", "error", err)
		return nil, err
	}

	var creditAccounts []plaid.AccountBase
	var creditTransactions []plaid.Transaction
	accountIds := make(map[string]string)
	for _, account := range transactionsResp.GetAccounts() {
		if account.Type == "credit" {
			creditAccounts = append(creditAccounts, account)
			accountIds[account.AccountId] = account.Name
		}
	}

	for _, transaction := range transactionsResp.GetTransactions() {
		if _, ok := accountIds[transaction.AccountId]; ok {
			creditTransactions = append(creditTransactions, transaction)
		}
	}
	transactionsResponse := models.TransactionsResponse{Accounts: creditAccounts, Transactions: creditTransactions}

	response, err := p.PlaidResponseToPB(liabilitiesResponse, transactionsResponse, token.User)
	if err != nil {
		p.l.Error("Error converting PlaidResponse to PB", "error", err)
		return nil, err
	}
	return response, nil
}

func (p *PlaidClient) PlaidResponseToPB(lr models.LiabilitiesResponse, tr models.TransactionsResponse, user *models.User) (*core.AccountDetailsResponse, error) {
	UserId := user.ID.Hex()

	accountLiabilities := make(map[string]plaid.CreditCardLiability)
	for _, al := range lr.Liabilities {
		if al.AccountId.IsSet() {
			accId := al.AccountId.Get()
			accountLiabilities[*accId] = al
		} else {
			p.l.Error("Error isolating accountLiabilities")
			return nil, errors.New("error isolating accountLiabilities")
		}
	}

	accounts := make([]*core.Account, len(tr.Accounts))
	for idx, account := range tr.Accounts {
		if acc, ok := accountLiabilities[account.AccountId]; ok {
			var isOverdue = false
			var nextPaymentDueDate = ""
			var officialName = ""
			var subtype = ""
			var availableBalance = 0.0
			var currentBalance = 0.0
			var creditLimit = 0.0
			var isoCurrencyCode = ""

			if acc.IsOverdue.IsSet() {
				i := acc.IsOverdue.Get()
				if i != nil {
					isOverdue = *i
				}
			}
			if acc.NextPaymentDueDate.IsSet() {
				i := acc.NextPaymentDueDate.Get()
				if i != nil {
					nextPaymentDueDate = *i
				}
			}

			aprs := make([]*core.AnnualPercentageRates, len(acc.Aprs))
			for x, apr := range acc.Aprs {
				var balanceSubjectToApr = 0.0
				var interestChargeAmount = 0.0
				if apr.BalanceSubjectToApr.IsSet() {
					i := apr.BalanceSubjectToApr.Get()
					if i != nil {
						balanceSubjectToApr = float64(*i)
					}

				}
				if apr.InterestChargeAmount.IsSet() {
					i := apr.InterestChargeAmount.Get()
					if i != nil {
						interestChargeAmount = float64(*i)
					}
				}

				aprs[x] = &core.AnnualPercentageRates{
					AprPercentage:        float64(apr.AprPercentage),
					AprType:              apr.AprType,
					BalanceSubjectToApr:  balanceSubjectToApr,
					InterestChargeAmount: interestChargeAmount,
				}
			}

			if account.OfficialName.IsSet() {
				officialName = *account.OfficialName.Get()
			}
			if account.Subtype.IsSet() {
				subtype = string(*account.Subtype.Get())
			}
			if account.Balances.Available.IsSet() {
				availableBalance = float64(*account.Balances.Available.Get())
			}
			if account.Balances.Current.IsSet() {
				currentBalance = float64(*account.Balances.Current.Get())
			}
			if account.Balances.Limit.IsSet() {
				creditLimit = float64(*account.Balances.Limit.Get())
			}
			if account.Balances.IsoCurrencyCode.IsSet() {
				isoCurrencyCode = *account.Balances.IsoCurrencyCode.Get()
			}

			accounts[idx] = &core.Account{
				UserId:                 UserId,
				Name:                   account.Name,
				OfficialName:           officialName,
				Type:                   string(account.Type),
				Subtype:                subtype,
				AvailableBalance:       availableBalance,
				CurrentBalance:         currentBalance,
				CreditLimit:            creditLimit,
				IsoCurrencyCode:        isoCurrencyCode,
				AnnualPercentageRate:   aprs,
				IsOverdue:              isOverdue,
				LastPaymentAmount:      float64(acc.LastPaymentAmount),
				LastStatementIssueDate: acc.LastStatementIssueDate,
				LastStatementBalance:   float64(acc.LastStatementBalance),
				MinimumPaymentAmount:   float64(acc.MinimumPaymentAmount),
				NextPaymentDueDate:     nextPaymentDueDate,
				PlaidAccountId:         account.AccountId,
			}
		}
	}
	var transactions []*core.Transaction
	for _, transaction := range tr.Transactions {
		var pendingTransactionId = ""
		var categoryId = ""
		var address = ""
		var city = ""
		var state = ""
		var zipcode = ""
		var country = ""
		var storeNumber = ""
		var referenceNumber = ""
		var originalDescription = ""
		var isoCurrencyCode = ""
		var merchantName = ""
		var authorizedDate = ""
		var primaryCategory = ""
		var detailedCategory = ""

		if transaction.PendingTransactionId.IsSet() {
			i := transaction.PendingTransactionId.Get()
			if i != nil {
				pendingTransactionId = *i
			}

		}
		if transaction.CategoryId.IsSet() {
			i := transaction.CategoryId.Get()
			if i != nil {
				categoryId = *i
			}
		}
		if transaction.Location.Address.IsSet() && transaction.Location.Address.Get() != nil {
			address = *transaction.Location.Address.Get()
		}
		if transaction.Location.City.IsSet() && transaction.Location.City.Get() != nil {
			city = *transaction.Location.City.Get()
		}
		if transaction.Location.Region.IsSet() && transaction.Location.Region.Get() != nil {
			state = *transaction.Location.Region.Get()
		}
		if transaction.Location.PostalCode.IsSet() && transaction.Location.PostalCode.Get() != nil {
			zipcode = *transaction.Location.PostalCode.Get()
		}
		if transaction.Location.Country.IsSet() && transaction.Location.Country.Get() != nil {
			country = *transaction.Location.Country.Get()
		}
		if transaction.Location.StoreNumber.IsSet() && transaction.Location.StoreNumber.Get() != nil {
			storeNumber = *transaction.Location.StoreNumber.Get()
		}
		if transaction.PaymentMeta.ReferenceNumber.IsSet() && transaction.PaymentMeta.ReferenceNumber.Get() != nil {
			referenceNumber = *transaction.PaymentMeta.ReferenceNumber.Get()
		}
		if transaction.OriginalDescription.IsSet() && transaction.OriginalDescription.Get() != nil {
			originalDescription = *transaction.OriginalDescription.Get()
		}
		if transaction.IsoCurrencyCode.IsSet() && transaction.IsoCurrencyCode.Get() != nil {
			isoCurrencyCode = *transaction.IsoCurrencyCode.Get()
		}
		if transaction.MerchantName.IsSet() && transaction.MerchantName.Get() != nil {
			merchantName = *transaction.MerchantName.Get()
		}
		if transaction.AuthorizedDate.IsSet() && transaction.AuthorizedDate.Get() != nil {
			authorizedDate = *transaction.AuthorizedDate.Get()
		}
		if transaction.PersonalFinanceCategory.IsSet() && transaction.PersonalFinanceCategory.Get() != nil {
			primaryCategory = transaction.PersonalFinanceCategory.Get().Primary
			detailedCategory = transaction.PersonalFinanceCategory.Get().Detailed
		}

		transactions = append(transactions, &core.Transaction{
			UserId:               UserId,
			TransactionType:      *transaction.TransactionType,
			PendingTransactionId: pendingTransactionId,
			CategoryId:           categoryId,
			Category:             transaction.Category,
			TransactionDetails: &core.TransactionDetails{
				Address:         address,
				City:            city,
				State:           state,
				Zipcode:         zipcode,
				Country:         country,
				StoreNumber:     storeNumber,
				ReferenceNumber: referenceNumber,
			},
			Name:                transaction.Name,
			OriginalDescription: originalDescription,
			Amount:              float64(transaction.Amount),
			IsoCurrencyCode:     isoCurrencyCode,
			Date:                transaction.Date,
			Pending:             transaction.Pending,
			MerchantName:        merchantName,
			PaymentChannel:      transaction.PaymentChannel,
			AuthorizedDate:      authorizedDate,
			PrimaryCategory:     primaryCategory,
			DetailedCategory:    detailedCategory,
			PlaidAccountId:      transaction.AccountId,
			PlaidTransactionId:  transaction.TransactionId,
		})

	}
	return &core.AccountDetailsResponse{
		Accounts:     accounts,
		Transactions: transactions,
	}, nil
}

// SaveToken method adds the permanent plaid token and stores into the plaid tokens' table with the
// same id as the user.
func (p *PlaidClient) SaveToken(ctx context.Context, token *models.Token) error {
	token.ID = primitive.NewObjectID()
	_, err := p.PlaidDB.InsertOne(ctx, token)
	if err != nil {
		p.l.Info("Error inserting new Token", err)
		return err
	}
	return nil
}

func (p *PlaidClient) UpdateToken(ctx context.Context, TokenId primitive.ObjectID, value, itemId string) error {
	filter := bson.D{{Key: "_id", Value: TokenId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "value", Value: value}, {Key: "item_id", Value: itemId}}}}
	_, err := p.PlaidDB.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetTokens returns every token associated to the user in the form of a slice of Token pointers.
func (p *PlaidClient) GetTokens(ctx context.Context, username string) ([]*models.Token, error) {
	var results []models.Token
	cursor, err := p.PlaidDB.Find(ctx, bson.D{{Key: "user.username", Value: username}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		p.l.Error("[PlaidDb] Error getting all users tokens", "error", err)
		return nil, err
	}
	tokens := make([]*models.Token, len(results))
	for idx, token := range results {
		tokens[idx] = &token
	}
	return tokens, nil
}

// GetToken will get a token from the database and return it given the user's ID and the
// token id
func (p *PlaidClient) GetToken(ctx context.Context, accessToken, tokenId string) (*models.Token, error) {
	var token models.Token
	var filter []bson.M

	if tokenId != "" {
		id, err := primitive.ObjectIDFromHex(tokenId)
		if err != nil {
			return nil, err
		}
		filter = []bson.M{{"_id": id}, {"value": accessToken}}
	} else {
		filter = []bson.M{{"value": accessToken}}
	}

	err := p.PlaidDB.FindOne(ctx, bson.M{"$or": filter}).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (p *PlaidClient) GetUserToken(ctx context.Context, user *models.User) (*models.Token, error) {
	var token models.Token
	filter := []bson.M{{"user._id": user.ID}, {"user.username": user.Username}, {"user.email": user.Email}}
	err := p.PlaidDB.FindOne(ctx, bson.M{"$or": filter}).Decode(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (p *PlaidClient) GetUser(ctx context.Context, username, userId string) (*models.User, error) {
	userRequest := &core.GetUserRequest{
		Id:       userId,
		Username: username,
	}
	user, err := p.CoreClient.GetUser(ctx, userRequest)
	if err != nil {
		return nil, err
	}
	id, err := primitive.ObjectIDFromHex(user.GetId())
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       id,
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
	}, nil
}

func (p *PlaidClient) SetLinkToken(token *models.Token) {
	p.LinkToken = token
}

func (p *PlaidClient) SetPublicToken(token *models.Token) {
	p.PublicToken = token
}

func (p *PlaidClient) GetLinkToken() *models.Token {
	return p.LinkToken
}

func (p *PlaidClient) GetPublicToken() *models.Token {
	return p.PublicToken
}

func convertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	var countryCodes []plaid.CountryCode

	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func convertProducts(productStrs []string) []plaid.Products {
	var products []plaid.Products

	for _, productStr := range productStrs {
		products = append(products, plaid.Products(productStr))
	}

	return products
}

func renderError(originalErr error) map[string]interface{} {
	resp := make(map[string]interface{})
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		resp["error"] = plaidError
		return resp
	}
	resp["error"] = originalErr.Error()
	return resp
}
