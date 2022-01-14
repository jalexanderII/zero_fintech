package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/core/config"
	"github.com/jalexanderII/zero_fintech/services/core/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"github.com/jalexanderII/zero_fintech/services/core/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	TokenDuration = 15 * time.Minute
)

func main() {
	// establish default logger with log levels
	l := hclog.Default()
	l.Debug("Core Service")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", config.GetEnv("CORE_SERVER_PORT")))
	if err != nil {
		l.Error("failed to listen", "error", err)
		panic(err)
	}

	// jwtManger to manage user authentication using tokens
	jwtManager := middleware.NewJWTManager(config.GetEnv("JWTSecret"), TokenDuration)

	// Initiate MongoDB Database
	DB, err := database.InitiateMongoClient()
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}
	// Connect to the Collections inside the given DB
	coreCollection := *DB.Collection(config.GetEnv("CORE_COLLECTION"))
	accountCollection := *DB.Collection(config.GetEnv("ACCOUNT_COLLECTION"))
	transactionCollection := *DB.Collection(config.GetEnv("TRANSACTION_COLLECTION"))
	userCollection := *DB.Collection(config.GetEnv("USER_COLLECTION"))

	var serverOptions []grpc.ServerOption
	// Initiate grpcServer instance
	grpcServer := grpc.NewServer(serverOptions...)

	// Bind grpcServer to CoreService Server defined by proto
	core.RegisterCoreServer(grpcServer,
		server.NewCoreServer(coreCollection, accountCollection, transactionCollection, userCollection, jwtManager, l),
	)
	methods := config.ListGRPCResources(grpcServer)
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
