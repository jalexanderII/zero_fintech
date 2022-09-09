package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jalexanderII/zero_fintech/bff/client"
	"github.com/jalexanderII/zero_fintech/bff/handlers"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a new instance of the logger.
var l = logrus.New()

func SetupRoutes(app *fiber.App, DB mongo.Database) {
	// Create a client and context to reuse in all handlers
	ctx := context.Background()
	authClient, grpcOpts := client.SetUpAuthClient()
	coreClient := client.SetUpCoreClient(authClient, grpcOpts)
	// Connect to the Collections inside the given DB
	plaidCollection := *DB.Collection(utils.GetEnv("PLAID_COLLECTION"))
	plaidClient := client.NewPlaidClient(l, plaidCollection, coreClient)

	// Create handlers for bff server
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Hello, World!") })

	api := app.Group("/api")

	// Monitoring api stats
	api.Get("/dashboard", monitor.New())

	// Auth endpoints
	authEndpoints := api.Group("/auth")
	authEndpoints.Post("/login", handlers.Login(authClient))
	authEndpoints.Post("/signup", handlers.SignUp(authClient))
	authEndpoints.Post("/logout", handlers.Logout(authClient))
	// Plaid endpoints within Auth
	plaidEndpoints := authEndpoints.Group("/plaid")
	plaidEndpoints.Post("/info", handlers.Info(plaidClient))
	plaidEndpoints.Get("/link/:username/:purpose", handlers.Link)
	plaidEndpoints.Post("/create_link", handlers.CreateLinkToken(plaidClient, ctx))
	plaidEndpoints.Post("/exchange", handlers.ExchangePublicToken(plaidClient, ctx))
	// plaidEndpoints.Patch("/exchange", handlers.ExchangePublicToken(plaidClient, ctx))

	// User endpoints
	userEndpoints := api.Group("/users")
	userEndpoints.Get("/", handlers.ListUsers(coreClient, ctx))
	userEndpoints.Get("/:id", handlers.GetUser(coreClient, ctx))
	userEndpoints.Patch("/:id", handlers.UpdateUser(coreClient, ctx))
	userEndpoints.Delete("/:id", handlers.DeleteUser(coreClient, ctx))

	// Core endpoints
	coreEndpoints := api.Group("/core")
	coreEndpoints.Post("/paymenttask", handlers.CreatePaymentTask(coreClient, ctx))
	coreEndpoints.Post("/paymentplan", handlers.GetPaymentPlan(coreClient, ctx))
	coreEndpoints.Get("/paymentplan/:id", handlers.GetUserPaymentPlans(coreClient, ctx))
	coreEndpoints.Get("/paymenttask", handlers.ListPaymentTasks(coreClient, ctx))
	coreEndpoints.Get("/paymenttask/:id", handlers.GetPaymentTask(coreClient, ctx))
	coreEndpoints.Patch("/paymenttask/:id", handlers.UpdatePaymentTask(coreClient, ctx))
	coreEndpoints.Delete("/paymenttask/:id", handlers.DeletePaymentTask(coreClient, ctx))
	coreEndpoints.Get("/accounts/:id", handlers.GetUserAccounts(coreClient, ctx))
	coreEndpoints.Get("/accounts/debit/balance/:id", handlers.GetUserDebitAccountBalance(coreClient, ctx))
	coreEndpoints.Get("/transactions/:id", handlers.GetUserTransactions(coreClient, ctx))
	coreEndpoints.Get("/accounts/credit/balance/:id", handlers.GetUserTotalCreditAccountBalance(coreClient, ctx))

	coreDashboardEndpoints := coreEndpoints.Group("/dashboard")
	coreDashboardEndpoints.Get("/waterfall/:id", handlers.GetWaterfallOverview(coreClient, ctx))
	coreDashboardEndpoints.Get("/amount_paid/:id", handlers.GetAmountPaidPercentage(coreClient, ctx))
	coreDashboardEndpoints.Get("/covered/:id", handlers.GetPercentageCoveredByPlans(coreClient, ctx))
}
