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
		p.l.Errorf("[Plaid Error] getting Liabilities %+v", renderError(err)["error"])
		return nil, err
	}
	liabilitiesResponse := models.LiabilitiesResponse{Liabilities: liabilitiesResp.GetLiabilities().Credit}
	time.Sleep(2 * time.Second)

	const iso8601TimeFormat = "2006-01-02"
	endDate := time.Now().Local().Format(iso8601TimeFormat)
	numMonths := time.Duration(-30 * 12 * 24)
	startDate := time.Now().Local().Add(numMonths * time.Hour).Format(iso8601TimeFormat)

	transactionsResp, _, err := p.Client.TransactionsGet(ctx).TransactionsGetRequest(
		*plaid.NewTransactionsGetRequest(token.Value, startDate, endDate),
	).Execute()
	if err != nil {
		p.l.Errorf("[Plaid Error] getting Transactions %+v", renderError(err)["error"])
		return nil, err
	}
	time.Sleep(2 * time.Second)

	var creditAccounts []plaid.AccountBase
	var creditTransactions []plaid.Transaction
	accountIds := make(map[string]string)
	for _, account := range transactionsResp.GetAccounts() {
		if account.Type == plaid.ACCOUNTTYPE_CREDIT {
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
			aprs := make([]*core.AnnualPercentageRates, len(acc.Aprs))
			for x, apr := range acc.Aprs {
				aprs[x] = &core.AnnualPercentageRates{
					AprPercentage:        float64(apr.AprPercentage),
					AprType:              apr.AprType,
					BalanceSubjectToApr:  float64(apr.GetBalanceSubjectToApr()),
					InterestChargeAmount: float64(apr.GetInterestChargeAmount()),
				}
			}

			accounts[idx] = &core.Account{
				UserId:                 UserId,
				Name:                   account.Name,
				OfficialName:           account.GetOfficialName(),
				Type:                   string(account.Type),
				Subtype:                string(account.GetSubtype()),
				AvailableBalance:       float64(account.Balances.GetAvailable()),
				CurrentBalance:         float64(account.Balances.GetCurrent()),
				CreditLimit:            float64(account.Balances.GetLimit()),
				IsoCurrencyCode:        account.Balances.GetIsoCurrencyCode(),
				AnnualPercentageRate:   aprs,
				IsOverdue:              acc.GetIsOverdue(),
				LastPaymentAmount:      float64(acc.LastPaymentAmount),
				LastStatementIssueDate: acc.LastStatementIssueDate,
				LastStatementBalance:   float64(acc.LastStatementBalance),
				MinimumPaymentAmount:   float64(acc.MinimumPaymentAmount),
				NextPaymentDueDate:     acc.GetNextPaymentDueDate(),
				PlaidAccountId:         account.AccountId,
			}
		}
	}
	var transactions []*core.Transaction
	for _, transaction := range tr.Transactions {
		transactions = append(transactions, &core.Transaction{
			UserId:               UserId,
			TransactionType:      transaction.GetTransactionType(),
			PendingTransactionId: transaction.GetPendingTransactionId(),
			CategoryId:           transaction.GetCategoryId(),
			Category:             transaction.Category,
			TransactionDetails: &core.TransactionDetails{
				Address:         transaction.Location.GetAddress(),
				City:            transaction.Location.GetCity(),
				State:           transaction.Location.GetRegion(),
				Zipcode:         transaction.Location.GetPostalCode(),
				Country:         transaction.Location.GetCountry(),
				StoreNumber:     transaction.Location.GetStoreNumber(),
				ReferenceNumber: transaction.PaymentMeta.GetReferenceNumber(),
			},
			Name:                transaction.Name,
			OriginalDescription: transaction.GetOriginalDescription(),
			Amount:              float64(transaction.Amount),
			IsoCurrencyCode:     transaction.GetIsoCurrencyCode(),
			Date:                transaction.Date,
			Pending:             transaction.Pending,
			MerchantName:        transaction.GetMerchantName(),
			PaymentChannel:      transaction.PaymentChannel,
			AuthorizedDate:      transaction.GetAuthorizedDate(),
			PrimaryCategory:     transaction.GetPersonalFinanceCategory().Primary,
			DetailedCategory:    transaction.GetPersonalFinanceCategory().Detailed,
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
