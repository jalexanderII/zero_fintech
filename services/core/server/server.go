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

type CoreServer struct {
	core.UnimplementedCoreServer
	PaymentTaskDB mongo.Collection
	AccountDB     mongo.Collection
	TransactionDB mongo.Collection
	UserDB        mongo.Collection
	jwtm          *middleware.JWTManager
	l             hclog.Logger
}

func NewCoreServer(pdb mongo.Collection, adb mongo.Collection, tdb mongo.Collection, udb mongo.Collection, jwtm *middleware.JWTManager, l hclog.Logger) *CoreServer {
	return &CoreServer{PaymentTaskDB: pdb, AccountDB: adb, TransactionDB: tdb, UserDB: udb, jwtm: jwtm, l: l}
}

func (s CoreServer) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest) (*core.GetPaymentPlanResponse, error) {
	listOfIds := bson.A{}
	for _, id := range in.GetPaymentTasksIds() {
		hex, _ := primitive.ObjectIDFromHex(id)
		listOfIds = append(listOfIds, hex)
	}

	var results []database.PaymentTask
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
	// call Planning client with res , get response
	return &core.GetPaymentPlanResponse{}, nil
}
