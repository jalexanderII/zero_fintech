package client

import (
	"context"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
)

func (c *CoreClient) CreateTransaction(ctx context.Context, in *core.CreateTransactionRequest, opts ...grpc.CallOption) (*core.Transaction, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) GetTransaction(ctx context.Context, in *core.GetTransactionRequest, opts ...grpc.CallOption) (*core.Transaction, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) ListTransactions(ctx context.Context, in *core.ListTransactionRequest, opts ...grpc.CallOption) (*core.ListTransactionResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) UpdateTransaction(ctx context.Context, in *core.UpdateTransactionRequest, opts ...grpc.CallOption) (*core.Transaction, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) DeleteTransaction(ctx context.Context, in *core.DeleteTransactionRequest, opts ...grpc.CallOption) (*core.DeleteTransactionResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
