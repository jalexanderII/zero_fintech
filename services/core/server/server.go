package server

import (
	"context"
	"time"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/planning"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func NewCoreServer(pdb mongo.Collection, adb mongo.Collection, tdb mongo.Collection, udb mongo.Collection, jwtm *middleware.JWTManager, pc planning.PlanningClient, l *logrus.Logger) *CoreServer {
	return &CoreServer{PaymentTaskDB: pdb, AccountDB: adb, TransactionDB: tdb, UserDB: udb, jwtm: jwtm, planningClient: pc, l: l}
}

func (s CoreServer) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest) (*core.GetPaymentPlanResponse, error) {
	// fetch payment plans from the database
	var results []database.PaymentTask
	listOfIds := bson.A{}
	// convert strings to their representative MongoDB primitive ID
	for _, id := range in.GetPaymentTasksIds() {
		hex, _ := primitive.ObjectIDFromHex(id)
		listOfIds = append(listOfIds, hex)
	}
	cursor, err := s.PaymentTaskDB.Find(ctx, bson.D{{"_id", bson.D{{"$in", listOfIds}}}})
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[PaymentTaskDB] Error getting all PaymentTasks", "error", err)
		return nil, err
	}
	req := make([]*core.PaymentTask, len(results))
	for idx, paymentTask := range results {
		req[idx] = PaymentTaskDBToPB(paymentTask)
	}

	res, err := s.planningClient.CreatePaymentPlan(ctx, &planning.CreatePaymentPlanRequest{PaymentTasks: req})
	s.l.Info("[Payment Plans] Response", "PaymentPlans", res)
	return &core.GetPaymentPlanResponse{PaymentPlans: res.GetPaymentPlans()}, nil
}

func MockClientCall(tasks []*core.PaymentTask) []*planning.PaymentPlan {
	var plans []*planning.PaymentPlan
	var actions []*planning.PaymentAction
	var ids []string
	var total float64
	for _, task := range tasks {
		ids = append(ids, task.PaymentTaskId)
		total += task.Amount
		a := &planning.PaymentAction{
			AccountId:       task.AccountId,
			Amount:          float32(task.GetAmount()),
			TransactionDate: timestamppb.New(time.Now()),
			Status:          planning.PaymentActionStatus_PAYMENT_ACTION_STATUS_PENDING,
		}
		actions = append(actions, a)
	}
	plan := &planning.PaymentPlan{
		PaymentPlanId:    primitive.NewObjectID().Hex(),
		UserId:           tasks[0].UserId,
		PaymentTaskId:    ids,
		Timeline:         12,
		PaymentFreq:      common.PaymentFrequency_PAYMENT_FREQUENCY_MONTHLY,
		AmountPerPayment: float32(total / 12),
		PlanType:         common.PlanType_PLAN_TYPE_OPTIM_CREDIT_SCORE,
		EndDate:          timestamppb.New(time.Now()),
		Active:           true,
		Status:           planning.PaymentStatus_PAYMENT_STATUS_CURRENT,
		PaymentAction:    actions,
	}
	plans = append(plans, plan)
	return plans
}
