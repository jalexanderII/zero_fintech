package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jalexanderII/zero_fintech/bff/client"
	"github.com/jalexanderII/zero_fintech/bff/handlers"
)

func SetupRoutes(app *fiber.App) {
	// Create a client and context to reuse in all handlers
	ctx := context.Background()
	authClient, grpcOpts := client.SetUpAuthClient()
	coreClient := client.SetUpCoreClient(authClient, grpcOpts)

	// Create handlers for bff server
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Hello, World!") })
	api := app.Group("/api")

	// Monitoring api stats
	api.Get("/dashboard", monitor.New())

	// Auth endpoints
	authEndpoints := api.Group("/auth")
	authEndpoints.Post("/login", handlers.Login(authClient))
	authEndpoints.Post("/signup", handlers.SignUp(authClient))

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
	coreEndpoints.Get("/paymenttask", handlers.ListPaymentTasks(coreClient, ctx))
	coreEndpoints.Get("/paymenttask/:id", handlers.GetPaymentTask(coreClient, ctx))
	coreEndpoints.Patch("/paymenttask/:id", handlers.UpdatePaymentTask(coreClient, ctx))
	coreEndpoints.Delete("/paymenttask/:id", handlers.DeletePaymentTask(coreClient, ctx))
}
