package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/client"
	"github.com/jalexanderII/zero_fintech/services/core/config/interceptor"
	"github.com/jalexanderII/zero_fintech/services/core/server"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	TokenDuration = 15 * time.Minute
)

// Create a new instance of the logger.
var l = logrus.New()

func main() {
	// establish default logger with log levels
	l.Debug("Core Service")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", utils.GetEnv("CORE_SERVER_PORT")))
	if err != nil {
		l.Error("failed to listen", "error", err)
		panic(err)
	}

	planningClient := client.SetUpPlanningClient()
	// jwtManger to manage user authentication using tokens
	jwtManager := middleware.NewJWTManager(utils.GetEnv("JWTSecret"), TokenDuration)
	authInterceptor := interceptor.NewAuthInterceptor(jwtManager, interceptor.AccessibleRoles(), l)

	// Initiate MongoDB Database
	DB, err := utils.InitiateMongoClient()
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}

	// Connect to the Collections inside the given DB
	coreCollection := *DB.Collection(utils.GetEnv("CORE_COLLECTION"))
	accountCollection := *DB.Collection(utils.GetEnv("ACCOUNT_COLLECTION"))
	transactionCollection := *DB.Collection(utils.GetEnv("TRANSACTION_COLLECTION"))
	userCollection := *DB.Collection(utils.GetEnv("USER_COLLECTION"))

	// Initiate grpcServer instance
	serverOptions := []grpc.ServerOption{grpc.UnaryInterceptor(authInterceptor.Unary())}
	grpcServer := grpc.NewServer(serverOptions...)

	// Bind grpcServer to CoreService Server defined by proto
	core.RegisterCoreServer(grpcServer,
		server.NewCoreServer(coreCollection, accountCollection, transactionCollection, userCollection, jwtManager, planningClient, l),
	)
	methods := utils.ListGRPCResources(grpcServer)
	l.Info("Methods on this server", "methods", methods)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)

	l.Info("Server started", "port", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Serving gRPC: ", err)
	}
}
