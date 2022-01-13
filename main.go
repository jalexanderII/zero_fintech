package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jalexanderII/zero_fintech/services/core/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	l := hclog.Default()
	l.Debug("Listings Service")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", config.GetEnv("CORE_SERVER_PORT")))
	if err != nil {
		l.Error("failed to listen", "error", err)
		panic(err)
	}

	var serverOptions []grpc.ServerOption
	grpcServer := grpc.NewServer(serverOptions...)

	methods := config.ListGRPCResources(grpcServer)
	l.Info("Methods on this server", "methods", methods)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)

	l.Info("Server started", "port", lis.Addr().String())
	go func() {
		log.Fatal("Serving gRPC: ", grpcServer.Serve(lis).Error())
	}()

	// From https://rogchap.com/2019/07/26/in-process-grpc-web-proxy/
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	httpServer := &http.Server{
		Addr: config.GetEnv("WEBPROXYPORT"),
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}
	log.Fatal("Serving Proxy: ", httpServer.ListenAndServe().Error())
}
