package client

import (
	"context"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
)

func (c *CoreClient) GetUser(ctx context.Context, in *core.GetUserRequest, opts ...grpc.CallOption) (*core.User, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) ListUsers(ctx context.Context, in *core.ListUserRequest, opts ...grpc.CallOption) (*core.ListUserResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) UpdateUser(ctx context.Context, in *core.UpdateUserRequest, opts ...grpc.CallOption) (*core.User, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) DeleteUser(ctx context.Context, in *core.DeleteUserRequest, opts ...grpc.CallOption) (*core.DeleteUserResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
