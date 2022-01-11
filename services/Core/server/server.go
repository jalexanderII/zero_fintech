package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
)

type CoreServer struct {
	core.UnimplementedCoreServer
	l hclog.Logger
}

func NewCoreServer(l hclog.Logger) *CoreServer {
	return &CoreServer{l: l}
}

func (s CoreServer) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest) (*core.GetPaymentPlanResponse, error) {
	return nil, nil
}
