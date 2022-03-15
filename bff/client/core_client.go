package client

import (
	"log"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"google.golang.org/grpc"
)

// SetUpCoreClient creates a new CoreClient and with the AuthClient Interceptor for token authorization
func SetUpCoreClient(authClient *AuthClient, opts []grpc.DialOption) core.CoreClient {
	opts = append(opts, grpc.WithUnaryInterceptor(authClient.Interceptor.Unary()))
	coreConn, err := grpc.Dial("localhost:9090", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return core.NewCoreClient(coreConn)
}
