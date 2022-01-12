package config

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

const Performance = 100

// ListGRPCResources is a helper function that lists all URLs that are registered on gRPC server.
// This makes it easy to register all the relevant routes in your HTTP router of choice.
func ListGRPCResources(server *grpc.Server) []string {
	var ret []string
	for serviceName, serviceInfo := range server.GetServiceInfo() {
		for _, methodInfo := range serviceInfo.Methods {
			fullResource := fmt.Sprintf("/%s/%s", serviceName, methodInfo.Name)
			ret = append(ret, fullResource)
		}
	}
	return ret
}

// NewDBContext returns a new Context according to app performance
func NewDBContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d*Performance/100)
}
