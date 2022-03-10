package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jalexanderII/zero_fintech/gen/Go/payments"
	"github.com/jalexanderII/zero_fintech/services/payments/plaidExternal/server"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Create a new instance of the logger.
var l = logrus.New()

func main() {
	// establish default logger with log levels
	l.Debug("Payments Service")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", utils.GetEnv("PAYMENTS_SERVER_PORT")))
	if err != nil {
		l.Error("failed to listen", "error", err)
		panic(err)
	}

	// Initiate grpcServer instance
	var serverOptions []grpc.ServerOption
	grpcServer := grpc.NewServer(serverOptions...)

	// Bind grpcServer to CoreService Server defined by proto
	payments.RegisterPlaidServer(grpcServer, server.NewPaymentsServer(l))
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
