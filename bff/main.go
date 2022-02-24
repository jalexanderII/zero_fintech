package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/jalexanderII/zero_fintech/bff/config"
	"github.com/jalexanderII/zero_fintech/bff/middleware"
	"github.com/jalexanderII/zero_fintech/bff/routes"
)

func main() {
	// Initialize standard Go views template engine
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	middleware.FiberMiddleware(app)
	routes.SetupRoutes(app)
	// Start server (with graceful shutdown).
	config.StartServerWithGracefulShutdown(app)
}
