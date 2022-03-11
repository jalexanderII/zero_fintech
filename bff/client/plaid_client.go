package client

import (
	"context"
	"strings"
	"time"

	"github.com/jalexanderII/zero_fintech/bff/models"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/payments"
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
	"sandbox":     utils.GetEnv("PLAID_SECRET_SANDBOX"),
	"development": utils.GetEnv("PLAID_SECRET_DEV"),
}

type PlaidClient struct {
	// Name of the service
	Name string
	// Client is the object that contains all database functionalities
	Client       *plaid.APIClient
	RedirectURL  string
	Products     []plaid.Products
	CountryCodes []plaid.CountryCode
	// custom logger
	l *logrus.Logger
	// Database collection
	PlaidDB mongo.Collection
	UserDB  mongo.Collection
}

func NewPlaidClient(l *logrus.Logger, pdb mongo.Collection, udb mongo.Collection) *PlaidClient {
	// set constants from env
	PlaidEnv := utils.GetEnv("PLAID_ENV")
	PlaidSecret := utils.GetEnv(environmentSecret[PlaidEnv])

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", utils.GetEnv("PLAID_CLIENT_ID"))
	configuration.AddDefaultHeader("PLAID-SECRET", PlaidSecret)
	configuration.UseEnvironment(environments[PlaidEnv])

	countryCodes := convertCountryCodes(strings.Split(utils.GetEnv("PLAID_COUNTRY_CODES"), ","))
	products := convertProducts(strings.Split(utils.GetEnv("PLAID_PRODUCTS"), ","))

	return &PlaidClient{
		Name:         "ZeroFintech",
		Client:       plaid.NewAPIClient(configuration),
		RedirectURL:  utils.GetEnv("PLAID_REDIRECT_URI"),
		Products:     products,
		CountryCodes: countryCodes,
		l:            l,
		PlaidDB:      pdb,
		UserDB:       udb,
	}
}

// LinkTokenCreate creates a link token using the specified parameters
func (p *PlaidClient) LinkTokenCreate(ctx context.Context, username string) (string, error) {
	DbUser, err := p.GetUser(ctx, username, "")
	if err != nil {
		return "", err
	}

	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: DbUser.ID.Hex(),
	}

	request := plaid.NewLinkTokenCreateRequest(p.Name, "en", p.CountryCodes, user)
	request.SetProducts(p.Products)
	request.SetRedirectUri(p.RedirectURL)

	linkTokenCreateResp, _, err := p.Client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	p.l.Info("new link token created for user", linkTokenCreateResp)
	return linkTokenCreateResp.GetLinkToken(), nil
}

// ExchangePublicToken this function takes care of creating the permanent access token
// that will be stored in the database for cross-platform connection to users' bank.
// If for whatever reason there is a problem with the client or public token, there
// are json responses and logs that will adequately reflect all issues
func (p *PlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (*models.Token, error) {
	// exchange the public_token for an access_token
	exchangePublicTokenResp, _, err := p.Client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()
	if err != nil {
		return nil, err
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	itemID := exchangePublicTokenResp.GetItemId()
	institutionRequest := plaid.NewInstitutionsGetByIdRequest(itemID, p.CountryCodes)
	institution, _, err := p.Client.PlaidApi.InstitutionsGetById(ctx).InstitutionsGetByIdRequest(*institutionRequest).Execute()
	if err != nil {
		return nil, err
	}

	p.l.Info("public token: " + publicToken)
	p.l.Info("access token: " + accessToken)
	p.l.Info("item ID: " + itemID)
	return &models.Token{Value: accessToken, ItemId: itemID, Institution: institution.Institution.GetName()}, nil
}

func (p *PlaidClient) GetAccountDetails(ctx context.Context, token *models.Token) (*payments.AccountDetailsResponse, error) {
	liabilitiesReq := plaid.NewLiabilitiesGetRequest(token.Value)
	liabilitiesResp, _, err := p.Client.PlaidApi.LiabilitiesGet(ctx).LiabilitiesGetRequest(*liabilitiesReq).Execute()
	if err != nil {
		return nil, err
	}
	liabilitiesResponse := models.LiabilitiesResponse{Liabilities: liabilitiesResp.GetLiabilities().Credit}

	const iso8601TimeFormat = "2006-01-02"
	endDate := time.Now().Local().Format(iso8601TimeFormat)
	numMonths := time.Duration(-30 * 12 * 24)
	startDate := time.Now().Local().Add(numMonths * time.Hour).Format(iso8601TimeFormat)

	transactionsResp, _, err := p.Client.PlaidApi.TransactionsGet(ctx).TransactionsGetRequest(
		*plaid.NewTransactionsGetRequest(token.Value, startDate, endDate),
	).Execute()
	if err != nil {
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

	return PlaidResponseToPB(liabilitiesResponse, transactionsResponse, &token.User), nil
}

func PlaidResponseToPB(lr models.LiabilitiesResponse, tr models.TransactionsResponse, user *models.User) *payments.AccountDetailsResponse {
	UserId := user.ID.Hex()
	accountLiabilities := make(map[string]plaid.CreditCardLiability)
	for _, al := range lr.Liabilities {
		accountLiabilities[*al.AccountId.Get()] = al
	}
	var accounts []*core.Account
	for _, account := range tr.Accounts {
		var acc plaid.CreditCardLiability
		if _, ok := accountLiabilities[account.AccountId]; ok {
			acc = accountLiabilities[account.AccountId]
			var aprs []*core.AnnualPercentageRates
			for _, apr := range acc.Aprs {

				aprs = append(aprs, &core.AnnualPercentageRates{
					AprPercentage:        float64(apr.AprPercentage),
					AprType:              apr.AprType,
					BalanceSubjectToApr:  float64(*apr.BalanceSubjectToApr.Get()),
					InterestChargeAmount: float64(*apr.InterestChargeAmount.Get()),
				})
			}
			accounts = append(accounts, &core.Account{
				UserId:                 UserId,
				Name:                   account.Name,
				OfficialName:           *account.OfficialName.Get(),
				Type:                   string(account.Type),
				Subtype:                string(*account.Subtype.Get()),
				AvailableBalance:       float64(*account.Balances.Available.Get()),
				CurrentBalance:         float64(*account.Balances.Current.Get()),
				CreditLimit:            float64(*account.Balances.Limit.Get()),
				IsoCurrencyCode:        *account.Balances.IsoCurrencyCode.Get(),
				AnnualPercentageRate:   aprs,
				IsOverdue:              *acc.IsOverdue.Get(),
				LastPaymentAmount:      float64(acc.LastPaymentAmount),
				LastStatementIssueDate: acc.LastStatementIssueDate,
				LastStatementBalance:   float64(acc.LastStatementBalance),
				MinimumPaymentAmount:   float64(acc.MinimumPaymentAmount),
				NextPaymentDueDate:     *acc.NextPaymentDueDate.Get(),
				PlaidAccountId:         account.AccountId,
			})
		}
	}
	var transactions []*core.Transaction
	for _, transaction := range tr.Transactions {
		transactions = append(transactions, &core.Transaction{
			UserId:               UserId,
			TransactionType:      *transaction.TransactionType,
			PendingTransactionId: *transaction.PendingTransactionId.Get(),
			CategoryId:           *transaction.CategoryId.Get(),
			Category:             transaction.Category,
			TransactionDetails: &core.TransactionDetails{
				Address:         *transaction.Location.Address.Get(),
				City:            *transaction.Location.City.Get(),
				State:           *transaction.Location.Region.Get(),
				Zipcode:         *transaction.Location.PostalCode.Get(),
				Country:         *transaction.Location.Country.Get(),
				StoreNumber:     *transaction.Location.StoreNumber.Get(),
				ReferenceNumber: *transaction.PaymentMeta.ReferenceNumber.Get(),
			},
			Name:                transaction.Name,
			OriginalDescription: *transaction.OriginalDescription.Get(),
			Amount:              float64(transaction.Amount),
			IsoCurrencyCode:     *transaction.IsoCurrencyCode.Get(),
			Date:                transaction.Date,
			Pending:             transaction.Pending,
			MerchantName:        *transaction.MerchantName.Get(),
			PaymentChannel:      transaction.PaymentChannel,
			AuthorizedDate:      *transaction.AuthorizedDate.Get(),
			PrimaryCategory:     transaction.PersonalFinanceCategory.Get().Primary,
			DetailedCategory:    transaction.PersonalFinanceCategory.Get().Detailed,
			PlaidAccountId:      transaction.AccountId,
			PlaidTransactionId:  transaction.TransactionId,
		})
	}
	return &payments.AccountDetailsResponse{
		Accounts:     accounts,
		Transactions: transactions,
	}
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

func (p *PlaidClient) UpdateToken(ctx context.Context, TokenId, value, itemId string) error {
	DbToken, err := p.GetToken(ctx, "", TokenId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: DbToken.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "value", Value: value}, {Key: "item_id", Value: itemId}}}}
	_, err = p.PlaidDB.UpdateOne(ctx, filter, update)
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

func (p *PlaidClient) GetUser(ctx context.Context, username, userId string) (*models.User, error) {
	var user models.User
	var filter []bson.M

	if userId != "" {
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return nil, err
		}
		filter = []bson.M{{"_id": id}, {"username": username}}
	} else {
		filter = []bson.M{{"username": username}}
	}

	err := p.UserDB.FindOne(ctx, bson.M{"$or": filter}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
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
