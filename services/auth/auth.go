package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jalexanderII/zero_fintech/gen/Go/auth"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/auth/database"
	"github.com/jalexanderII/zero_fintech/services/auth/server"
	"github.com/jalexanderII/zero_fintech/services/core/config"
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
	l.Debug("Auth Service")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", utils.GetEnv("AUTH_SERVER_PORT")))
	if err != nil {
		l.Error("failed to listen", "error", err)
		panic(err)
	}

	// jwtManger to manage user authentication using tokens
	jwtManager := middleware.NewJWTManager(utils.GetEnv("JWTSecret"), TokenDuration)

	// Initiate MongoDB Database
	DB, err := database.InitiateMongoClient()
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}

	// Connect to the Collections inside the given DB
	userCollection := *DB.Collection(utils.GetEnv("USER_COLLECTION"))

	// Initiate grpcServer instance
	var serverOptions []grpc.ServerOption
	grpcServer := grpc.NewServer(serverOptions...)

	// Bind grpcServer to AuthService Server defined by proto
	auth.RegisterAuthServer(grpcServer, server.NewAuthServer(userCollection, jwtManager, l))
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
