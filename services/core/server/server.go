package server

import (
	"context"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/payments"
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
	plaidClient    payments.PlaidClient
	// custom logger
	l *logrus.Logger
}

func NewCoreServer(pdb mongo.Collection, adb mongo.Collection, tdb mongo.Collection, udb mongo.Collection, jwtm *middleware.JWTManager, pc planning.PlanningClient, plaid payments.PlaidClient, l *logrus.Logger) *CoreServer {
	return &CoreServer{PaymentTaskDB: pdb, AccountDB: adb, TransactionDB: tdb, UserDB: udb, jwtm: jwtm, planningClient: pc, plaidClient: plaid, l: l}
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
	res, err := s.planningClient.CreatePaymentPlan(ctx, &planning.CreatePaymentPlanRequest{PaymentTasks: paymentTasks, MetaData: in.GetMetaData()})
	if err != nil {
		return nil, err
	}
	return &common.PaymentPlanResponse{PaymentPlans: res.GetPaymentPlans()}, nil
}

func (s CoreServer) GetAccountDetails(ctx context.Context, in *payments.GetAccountDetailsRequest) (*payments.GetAccountDetailsResponse, error) {
	resp, err := s.plaidClient.GetAccountDetails(ctx, in)
	if err != nil {
		return nil, err
	}
	return &payments.GetAccountDetailsResponse{AccountDetailsResponse: resp.GetAccountDetailsResponse()}, nil
}
