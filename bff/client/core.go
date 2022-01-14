package client

import (
	"context"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
)

func (c *CoreClient) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest, opts ...grpc.CallOption) (*core.GetPaymentPlanResponse, error) {
	if err := prepareCoreGrpcClient(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
