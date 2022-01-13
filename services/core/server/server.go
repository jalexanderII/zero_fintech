package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/Core/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/Core/database"
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
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
	res := make([]*core.PaymentTask, len(results))
	for idx, paymentTask := range results {
		res[idx] = PaymentTaskDBToPB(paymentTask)
	}

	// TODO:
	// call Planning client with res and get back a response of Payment Plans
	return &core.GetPaymentPlanResponse{}, nil
}
