package client

import (
	"context"
	"errors"

	"github.com/jalexanderII/zero_fintech/services/core/config"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoreClient struct {
}

var (
	coreGrpcService       = config.GetEnv("CORE_SERVER_PORT")
	coreGrpcServiceClient core.CoreClient
)

func prepareCoreGrpcClient(ctx context.Context) error {

	conn, err := grpc.DialContext(ctx, coreGrpcService, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock()}...)

	if err != nil {
		coreGrpcServiceClient = nil
		return errors.New("connection to core gRPC service failed")
	}

	if coreGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	coreGrpcServiceClient = core.NewCoreClient(conn)
	return nil
}
