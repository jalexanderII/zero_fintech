package config

import (
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/plaid/plaid-go/plaid"
)

var (
	PLAID_CLIENT_ID                      = ""
	PLAID_SECRET                         = ""
	PLAID_ENV                            = ""
	PLAID_PRODUCTS                       = ""
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	APP_PORT                             = ""
	client              *plaid.APIClient = nil
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

func init() {
	// set constants from env
	PLAID_CLIENT_ID = utils.GetEnv("PLAID_CLIENT_ID")
	PLAID_SECRET = utils.GetEnv("PLAID_SECRET")
	PLAID_ENV = utils.GetEnv("PLAID_ENV")
	PLAID_PRODUCTS = utils.GetEnv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = utils.GetEnv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = utils.GetEnv("PLAID_REDIRECT_URI")
	APP_PORT = utils.GetEnv("PLAID_APP_PORT")

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", PLAID_SECRET)
	configuration.UseEnvironment(environments[PLAID_ENV])
	client = plaid.NewAPIClient(configuration)
}
