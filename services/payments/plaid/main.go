package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/services/payments/plaid/config"
	"github.com/jalexanderII/zero_fintech/services/payments/plaid/middleware"
	"github.com/jalexanderII/zero_fintech/services/payments/plaid/routes"
)

func main() {
	app := fiber.New()
	middleware.FiberMiddleware(app)
	routes.SetupRoutes(app)
	// Start server (with graceful shutdown).
	config.StartServerWithGracefulShutdown(app)
}
