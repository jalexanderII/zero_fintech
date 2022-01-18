package routes

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jalexanderII/zero_fintech/bff/handlers"
	"github.com/jalexanderII/zero_fintech/bff/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	timeout = 10 * time.Second
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func SetupRoutes(app *fiber.App) {
	// ******CLIENTS*******
	// create client and context with timeout to reuse in all handlers
	ctx := context.Background()

	// set up universal dail options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	// AuthClient Connection
	authConn, err := grpc.Dial("localhost:9091", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	// create new auth client
	authClient := middleware.NewAuthClient(authConn, "", "", "")

	// add auth interceptor middleware to core client
	opts = append(opts, grpc.WithUnaryInterceptor(authClient.Interceptor.Unary()))
	coreConn, err := grpc.Dial("localhost:9090", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	// instantiate core client
	coreClient := core.NewCoreClient(coreConn)

	// ******HANDLERS*******
	// Set up handlers
	// welcome screen
	app.Get("/", welcome)
	api := app.Group("/api")

	// // monitoring api stats
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

// NewClientContext returns a new Context according to app performance
func NewClientContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
