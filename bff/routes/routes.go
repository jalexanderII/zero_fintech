package routes

import (
	"context"

	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jalexanderII/zero_fintech/bff/client"
	"github.com/jalexanderII/zero_fintech/bff/handlers"
	"github.com/jalexanderII/zero_fintech/bff/middleware"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a new instance of the logger.
var l = logrus.New()

func protect(f fiber.Handler) fiber.Handler {
	return adaptor.HTTPHandler(middleware.EnsureValidToken()(adaptor.FiberHandlerFunc(f)))
}

func greet(msg string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": msg})
	}
}

func scopedGreet(c *fiber.Ctx) error {
	token := c.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope("read:admin-messages") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Insufficient scope."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Hello from a private endpoint! You need to be authenticated to see this."})
}

func SetupRoutes(app *fiber.App, DB mongo.Database) {
	// Create a client and context to reuse in all handlers
	ctx := context.Background()
	authClient, grpcOpts := client.SetUpAuthClient()
	coreClient := client.SetUpCoreClient(authClient, grpcOpts)
	// Connect to the Collections inside the given DB
	plaidCollection := *DB.Collection(utils.GetEnv("PLAID_COLLECTION"))
	plaidClient := client.NewPlaidClient(l, plaidCollection, coreClient)

	// Create handlers for bff server
	app.Get("/", greet("Hello, World!"))

	api := app.Group("/api")

	// Monitoring api stats
	api.Get("/dashboard", monitor.New())

	authTest := api.Group("/messages")
	authTest.Get("/admin", greet("Admin Hi"))
	authTest.Get("/public", greet("Hello from a public endpoint! You don't need to be authenticated to see this."))
	authTest.Get("/protected", protect(greet("Hello from a private endpoint! You need to be authenticated to see this.")))
	authTest.Get("/admin", protect(scopedGreet))

	// Auth endpoints
	authEndpoints := api.Group("/auth")
	authEndpoints.Post("/login", handlers.Login(authClient))
	authEndpoints.Post("/signup", handlers.SignUp(authClient))
	authEndpoints.Post("/logout", handlers.Logout(authClient))
	// Plaid endpoints within Auth
	plaidEndpoints := authEndpoints.Group("/plaid")
	plaidEndpoints.Post("/info", handlers.Info(plaidClient))
	plaidEndpoints.Get("/link/:email/:purpose", handlers.Link)
	plaidEndpoints.Post("/create_link", protect(handlers.CreateLinkToken(plaidClient, ctx)))
	plaidEndpoints.Post("/exchange", protect(handlers.ExchangePublicToken(plaidClient, ctx)))
	plaidEndpoints.Post("/create_link/internal", handlers.CreateLinkToken(plaidClient, ctx))
	plaidEndpoints.Post("/exchange/internal", handlers.ExchangePublicToken(plaidClient, ctx))
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
	coreEndpoints.Post("/paymentplan/:email", handlers.GetPaymentPlan(coreClient, ctx))
	coreEndpoints.Get("/paymentplan/:email", handlers.GetUserPaymentPlans(coreClient, ctx))
	coreEndpoints.Get("/paymenttask", handlers.ListPaymentTasks(coreClient, ctx))
	coreEndpoints.Get("/paymenttask/:id", handlers.GetPaymentTask(coreClient, ctx))
	coreEndpoints.Patch("/paymenttask/:id", handlers.UpdatePaymentTask(coreClient, ctx))
	coreEndpoints.Delete("/paymenttask/:id", handlers.DeletePaymentTask(coreClient, ctx))
	coreEndpoints.Get("/accounts/:email", handlers.GetUserAccounts(coreClient, ctx))
	coreEndpoints.Get("/accounts/debit/balance/:email", handlers.GetUserDebitAccountBalance(coreClient, ctx))
	coreEndpoints.Get("/transactions/:email", handlers.GetUserTransactions(coreClient, ctx))
	coreEndpoints.Get("/accounts/credit/balance/:email", handlers.GetUserTotalCreditAccountBalance(coreClient, ctx))
	coreEndpoints.Get("/accounts/credit/exist/:email", handlers.IsCreditAccountLinked(coreClient, ctx))
	coreEndpoints.Get("/accounts/debit/exist/:email", handlers.IsDebitAccountLinked(coreClient, ctx))
	coreEndpoints.Get("/accounts/plaid/buttons/exist/:email", handlers.ArePlaidAccountsLinked(coreClient, ctx))
	coreEndpoints.Get("/kpi/:email", handlers.GetUserKPIs(coreClient, ctx))

	coreDashboardEndpoints := coreEndpoints.Group("/dashboard")
	coreDashboardEndpoints.Get("/waterfall/:email", handlers.GetWaterfallOverview(coreClient, ctx))
	coreDashboardEndpoints.Get("/amount_paid/:id", handlers.GetAmountPaidPercentage(coreClient, ctx))
	coreDashboardEndpoints.Get("/covered/:id", handlers.GetPercentageCoveredByPlans(coreClient, ctx))
}
