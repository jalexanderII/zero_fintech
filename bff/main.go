package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/bff/config"
	"github.com/jalexanderII/zero_fintech/bff/middleware"
	"github.com/jalexanderII/zero_fintech/bff/routes"
)

func main() {
	app := fiber.New()
	middleware.FiberMiddleware(app)
	routes.SetupRoutes(app)

	// Start server (with graceful shutdown).
	config.StartServerWithGracefulShutdown(app)
}
