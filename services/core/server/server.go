package server

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/core/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
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
	// custom logger
	l hclog.Logger
}

func NewCoreServer(pdb mongo.Collection, adb mongo.Collection, tdb mongo.Collection, udb mongo.Collection, jwtm *middleware.JWTManager, l hclog.Logger) *CoreServer {
	return &CoreServer{PaymentTaskDB: pdb, AccountDB: adb, TransactionDB: tdb, UserDB: udb, jwtm: jwtm, l: l}
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

	// TODO:
	// call Planning client with res and get back a response of Payment Plans
	res := MockClientCall(req)
	s.l.Info("[Payment Plans] Response", "PaymentPlans", res)
	return &core.GetPaymentPlanResponse{PaymentPlans: res}, nil
}

func MockClientCall(tasks []*core.PaymentTask) []*core.PaymentPlan {
	var plans []*core.PaymentPlan
	var actions []*core.PaymentAction
	var ids []string
	var total float64
	for _, task := range tasks {
		ids = append(ids, task.PaymentTaskId)
		total += task.Amount
		a := &core.PaymentAction{
			AccountId:       task.AccountId,
			Amount:          float32(task.GetAmount()),
			TransactionDate: timestamppb.New(time.Now()),
			Status:          core.PaymentActionStatus_PAYMENT_ACTION_STATUS_PENDING,
		}
		actions = append(actions, a)
	}
	plan := &core.PaymentPlan{
		PaymentPlanId:    primitive.NewObjectID().Hex(),
		UserId:           tasks[0].UserId,
		PaymentTaskId:    ids,
		Timeline:         12,
		PaymentFreq:      core.PaymentFrequency_PAYMENTFREQ_MONTHLY,
		AmountPerPayment: float32(total / 12),
		PlanType:         core.PlanType_PLANTYPE_OPTIM_CREDIT_SCORE,
		EndDate:          timestamppb.New(time.Now()),
		Active:           true,
		Status:           core.PaymentStatus_PAYMENT_STATUS_CURRENT,
		PaymentAction:    actions,
	}
	plans = append(plans, plan)
	return plans
}
