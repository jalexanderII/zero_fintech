package server

import (
	"context"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/planning"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	cursor, err := s.PaymentTaskDB.Find(ctx, bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: listOfIds}}}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[PaymentTaskDB] Error getting all PaymentTasks", "error", err)
		return nil, err
	}
	req := make([]*common.PaymentTask, len(results))
	for idx, paymentTask := range results {
		req[idx] = PaymentTaskDBToPB(paymentTask)
	}

	res, err := s.planningClient.CreatePaymentPlan(ctx, &planning.CreatePaymentPlanRequest{PaymentTasks: req})
	if err != nil {
		return nil, err
	}
	s.l.Info("[Payment Plans] Response", "PaymentPlans", res)
	return &core.GetPaymentPlanResponse{PaymentPlans: res.GetPaymentPlans()}, nil
}
