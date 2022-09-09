package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/bff/config"
	"github.com/jalexanderII/zero_fintech/bff/middleware"
	"github.com/jalexanderII/zero_fintech/bff/routes"
	"github.com/jalexanderII/zero_fintech/utils"
)

func main() {
	// Initialize standard Go views template engine
	// engine := html.New("./views", ".html")

	// Initiate MongoDB Database
	DB, err := utils.InitiateMongoClient()
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}

	// app := fiber.New(fiber.Config{Views: engine})
	app := fiber.New()
	middleware.FiberMiddleware(app)
	routes.SetupRoutes(app, DB)
	// Start server (with graceful shutdown).
	config.StartServerWithGracefulShutdown(app)
}
