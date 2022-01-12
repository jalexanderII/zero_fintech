package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoreServer struct {
	core.UnimplementedCoreServer
	DB mongo.Collection
	l  hclog.Logger
}

func NewCoreServer(db mongo.Collection, l hclog.Logger) *CoreServer {
	return &CoreServer{DB: db, l: l}
}

func (s CoreServer) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest) (*core.GetPaymentPlanResponse, error) {

	return nil, nil
}
