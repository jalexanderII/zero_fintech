package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/Core/config"
	"github.com/jalexanderII/zero_fintech/services/Core/database"
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
	"github.com/jalexanderII/zero_fintech/services/Core/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	l := hclog.Default()
	l.Debug("Core Service")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", config.GetEnv("CORE_SERVER_PORT")))
	if err != nil {
		l.Error("failed to listen", "error", err)
		panic(err)
	}

	DB := database.InitiateMongoClient()
	coreCollection := *DB.Collection(config.GetEnv("CORE_COLLECTION"))
	accountCollection := *DB.Collection(config.GetEnv("ACCOUNT_COLLECTION"))
	transactionCollection := *DB.Collection(config.GetEnv("TRANSACTION_COLLECTION"))

	var serverOptions []grpc.ServerOption
	grpcServer := grpc.NewServer(serverOptions...)

	core.RegisterCoreServer(grpcServer, server.NewCoreServer(coreCollection, accountCollection, transactionCollection, l))
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
