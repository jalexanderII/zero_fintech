package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/Core/database"
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
)

type CoreServer struct {
	core.UnimplementedCoreServer
	DB *database.CoreDB
	l  hclog.Logger
}

func NewCoreServer(db *database.CoreDB, l hclog.Logger) *CoreServer {
	return &CoreServer{DB: db, l: l}
}

func (s CoreServer) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest) (*core.GetPaymentPlanResponse, error) {
	return nil, nil
}
