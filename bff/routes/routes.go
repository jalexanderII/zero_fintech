package routes

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jalexanderII/zero_fintech/bff/handlers"
	"github.com/jalexanderII/zero_fintech/bff/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/config"
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
	ctx, cancel := NewClientContext(timeout)
	defer cancel()

	// set up universal dail options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	// AuthClient Connection
	authConn, err := grpc.Dial(config.GetEnv("AUTH_SERVER_PORT"), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer authConn.Close()
	// create new auth client
	authClient := middleware.NewAuthClient(authConn, "", "", "")

	// add auth interceptor middleware to core client
	opts = append(opts, grpc.WithUnaryInterceptor(handlers.Interceptor.Unary()))
	coreConn, err := grpc.Dial(config.GetEnv("CORE_SERVER_PORT"), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer coreConn.Close()
	// instantiate core client
	coreClient := core.NewCoreClient(coreConn)

	// ******HANDLERS*******
	// Set up handlers
	// welcome screen
	app.Get("/", welcome)
	api := app.Group("/api")

	// monitoring api stats
	api.Get("/dashboard", monitor.New())

	// Auth endpoints
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login(authClient))
	auth.Post("/signup", handlers.SignUp(authClient))

	// User endpoints
	// users := api.Group("/users")
	// users.Get("/", , handlers.ListUsers(coreClient, ctx))
	// users.Get("/:id", , handlers.GetUser(coreClient, ctx))
	// users.Patch("/:id", , handlers.UpdateUser(coreClient, ctx))
	// users.Delete("/:id", , handlers.DeleteUser(coreClient, ctx))

	// Core endpoints
	coreH := api.Group("/core")
	coreH.Post("/paymenttask", handlers.CreatePaymentTask(coreClient, ctx))
	// coreH.Post("/paymentplan", handlers.GetPaymentPlan(coreClient, ctx))
	// coreH.Get("/paymenttask", handlers.ListPaymentTasks(coreClient, ctx))
	// coreH.Get("/paymenttask/:id", handlers.GetPaymentTask(coreClient, ctx))
	// coreH.Patch("/paymenttask/:id", handlers.UpdatePaymentTask(coreClient, ctx))
	// coreH.Delete("/paymenttask/:id", handlers.DeletePaymentTask(coreClient, ctx))
}

// NewClientContext returns a new Context according to app performance
func NewClientContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
