package client

import (
	"context"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
)

func (c *CoreClient) CreatePaymentTask(ctx context.Context, in *core.CreatePaymentTaskRequest, opts ...grpc.CallOption) (*core.PaymentTask, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) GetPaymentTask(ctx context.Context, in *core.GetPaymentTaskRequest, opts ...grpc.CallOption) (*core.PaymentTask, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) ListPaymentTasks(ctx context.Context, in *core.ListPaymentTaskRequest, opts ...grpc.CallOption) (*core.ListPaymentTaskResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) UpdatePaymentTask(ctx context.Context, in *core.UpdatePaymentTaskRequest, opts ...grpc.CallOption) (*core.PaymentTask, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
func (c *CoreClient) DeletePaymentTask(ctx context.Context, in *core.DeletePaymentTaskRequest, opts ...grpc.CallOption) (*core.DeletePaymentTaskResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
