package client

import (
	"context"
	"strings"

	"github.com/jalexanderII/zero_fintech/bff/models"
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
	// Client is the object that contains all database functionalities
	Client       *plaid.APIClient
	RedirectURL  string
	Products     []plaid.Products
	CountryCodes []plaid.CountryCode
	// custom logger
	l *logrus.Logger
	// Database collection
	PlaidDB mongo.Collection
}

func NewPlaidClient(l *logrus.Logger, pdb mongo.Collection) *PlaidClient {
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

	return &PlaidClient{Client: plaid.NewAPIClient(configuration), RedirectURL: utils.GetEnv("PLAID_REDIRECT_URI"), Products: products, CountryCodes: countryCodes, l: l, PlaidDB: pdb}
}

// LinkTokenCreate creates a link token using the specified parameters
func (p *PlaidClient) LinkTokenCreate(ctx context.Context, username, token string) (string, error) {
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: token,
	}

	request := plaid.NewLinkTokenCreateRequest(username, "en", p.CountryCodes, user)
	request.SetProducts(p.Products)
	request.SetRedirectUri(p.RedirectURL)

	linkTokenCreateResp, _, err := p.Client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	if err != nil {
		return "", err
	}
	p.l.Info("new link token created for user", linkTokenCreateResp)
	return linkTokenCreateResp.GetLinkToken(), nil
}

// UpdatePlaidToken returns a link token in update mode to allow user re-authentication
func (p *PlaidClient) UpdatePlaidToken(ctx context.Context, username string, token models.Token) (string, error) {
	return "", nil
}

func (p *PlaidClient) GetAccountDetails(ctx context.Context, username string, token models.Token) (string, error) {
	return "", nil
}

// ExchangePublicToken this function takes care of creating the permanent access token
// that will be stored in the database for cross-platform connection to users' bank.
// If for whatever reason there is a problem with the client or public token, there
// are json responses and logs that will adequately reflect all issues
func (p *PlaidClient) ExchangePublicToken(ctx context.Context, username string, token models.Token) (string, error) {
	return "", nil
}

// SaveToken method adds the permanent plaid token and stores into the plaid tokens' table with the
// same id as the user.
func (p *PlaidClient) SaveToken(ctx context.Context, token models.Token) error {
	token.ID = primitive.NewObjectID()
	_, err := p.PlaidDB.InsertOne(ctx, token)
	if err != nil {
		p.l.Info("Error inserting new Token", err)
		return err
	}
	return nil
}

func (p *PlaidClient) UpdateToken(ctx context.Context, tokenId, value, itemId string) error {
	id, err := primitive.ObjectIDFromHex(tokenId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: id}}
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
	cursor, err := p.PlaidDB.Find(ctx, bson.D{{Key: "username", Value: username}})
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
func (p *PlaidClient) GetToken(ctx context.Context, tokenId string) (*models.Token, error) {
	var Token models.Token
	id, err := primitive.ObjectIDFromHex(tokenId)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	err = p.PlaidDB.FindOne(ctx, filter).Decode(&Token)
	if err != nil {
		return nil, err
	}
	return &Token, nil
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
