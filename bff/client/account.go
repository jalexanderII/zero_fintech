package client

import (
	"context"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
)

func (c *CoreClient) CreateAccount(ctx context.Context, in *core.CreateAccountRequest, opts ...grpc.CallOption) (*core.Account, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) GetAccount(ctx context.Context, in *core.GetAccountRequest, opts ...grpc.CallOption) (*core.Account, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) ListAccounts(ctx context.Context, in *core.ListAccountRequest, opts ...grpc.CallOption) (*core.ListAccountResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) UpdateAccount(ctx context.Context, in *core.UpdateAccountRequest, opts ...grpc.CallOption) (*core.Account, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) DeleteAccount(ctx context.Context, in *core.DeleteAccountRequest, opts ...grpc.CallOption) (*core.DeleteAccountResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
