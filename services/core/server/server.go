package server

import (
	"context"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/planning"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// CoreServer is the server for the CoreService, it will connect to its own mongodb database and will be reachable via
// grpc from microservices and via grpc proxy for clients
type CoreServer struct {
	core.UnimplementedCoreServer
	// Database collections
	PaymentTaskDB mongo.Collection
	AccountDB     mongo.Collection
	TransactionDB mongo.Collection
	UserDB        mongo.Collection
	// authentication manager
	jwtm *middleware.JWTManager
	// clients
	planningClient planning.PlanningClient
	// custom logger
	l *logrus.Logger
}

func NewCoreServer(pdb mongo.Collection, adb mongo.Collection, tdb mongo.Collection,
	udb mongo.Collection, jwtm *middleware.JWTManager, pc planning.PlanningClient, l *logrus.Logger,
) *CoreServer {
	return &CoreServer{
		PaymentTaskDB: pdb, AccountDB: adb, TransactionDB: tdb,
		UserDB: udb, jwtm: jwtm, planningClient: pc, l: l,
	}
}

func (s CoreServer) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest) (*common.PaymentPlanResponse, error) {
	// create payment task from user inputs
	paymentTasks := make([]*common.PaymentTask, len(in.GetAccountInfo()))

	for idx, item := range in.GetAccountInfo() {
		task := &common.PaymentTask{
			UserId:    in.UserId,
			AccountId: item.AccountId,
			Amount:    item.Amount,
		}
		paymentTasks[idx] = task
	}

	// save payment tasks to DB
	listOfIds, err := s.CreateManyPaymentTask(ctx, &common.CreateManyPaymentTaskRequest{PaymentTasks: paymentTasks})
	if err != nil {
		s.l.Error("[PaymentTask] Error creating PaymentTasks", "error", err)
		return nil, err
	}

	for idx, id := range listOfIds.GetPaymentTaskIds() {
		pt, _ := s.GetPaymentTask(ctx, &common.GetPaymentTaskRequest{Id: id})
		paymentTasks[idx] = pt
	}

	// send payment tasks to planning to get payment plans
	res, err := s.planningClient.CreatePaymentPlan(ctx, &planning.CreatePaymentPlanRequest{PaymentTasks: paymentTasks, MetaData: in.GetMetaData(), SavePlan: in.GetSavePlan()})
	if err != nil {
		return nil, err
	}
	return &common.PaymentPlanResponse{PaymentPlans: res.GetPaymentPlans()}, nil
}

func (s CoreServer) GetWaterfallOverview(ctx context.Context, in *planning.GetUserOverviewRequest) (*planning.WaterfallOverviewResponse, error) {
	return s.planningClient.GetWaterfallOverview(ctx, in)
}

func (s CoreServer) GetAmountPaidPercentage(ctx context.Context, in *planning.GetUserOverviewRequest) (*planning.GetAmountPaidPercentageResponse, error) {
	return s.planningClient.GetAmountPaidPercentage(ctx, in)
}

func (s CoreServer) GetPercentageCoveredByPlans(ctx context.Context, in *planning.GetUserOverviewRequest) (*planning.GetPercentageCoveredByPlansResponse, error) {
	return s.planningClient.GetPercentageCoveredByPlans(ctx, in)
}
