package server

import (
	"context"
	"log"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s CoreServer) CreatePaymentTask(ctx context.Context, in *core.CreatePaymentTaskRequest) (*common.PaymentTask, error) {
	paymentTask := in.GetPaymentTask()
	newPaymentTask := PaymentTaskPBToDB(paymentTask, primitive.NewObjectID())

	_, err := s.PaymentTaskDB.InsertOne(ctx, newPaymentTask)
	if err != nil {
		log.Printf("Error inserting new PaymentTask: %v\n", err)
		return nil, err
	}
	return paymentTask, nil
}

func (s CoreServer) GetPaymentTask(ctx context.Context, in *core.GetPaymentTaskRequest) (*common.PaymentTask, error) {
	var paymentTask database.PaymentTask
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	err = s.PaymentTaskDB.FindOne(ctx, filter).Decode(&paymentTask)
	if err != nil {
		return nil, err
	}
	return PaymentTaskDBToPB(paymentTask), nil
}

func (s CoreServer) ListPaymentTasks(ctx context.Context, in *core.ListPaymentTaskRequest) (*core.ListPaymentTaskResponse, error) {
	var results []database.PaymentTask
	cursor, err := s.PaymentTaskDB.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[PaymentTaskDB] Error getting all users", "error", err)
		return nil, err
	}
	res := make([]*common.PaymentTask, len(results))
	for idx, paymentTask := range results {
		res[idx] = PaymentTaskDBToPB(paymentTask)
	}
	return &core.ListPaymentTaskResponse{PaymentTasks: res}, nil
}

func (s CoreServer) UpdatePaymentTask(ctx context.Context, in *core.UpdatePaymentTaskRequest) (*common.PaymentTask, error) {
	paymentTask := in.GetPaymentTask()
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "amount", Value: paymentTask.Amount}}}}
	_, err = s.PaymentTaskDB.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	var payment_task database.PaymentTask
	err = s.PaymentTaskDB.FindOne(ctx, filter).Decode(&payment_task)
	if err != nil {
		return nil, err
	}
	return PaymentTaskDBToPB(payment_task), nil
}
func (s CoreServer) DeletePaymentTask(ctx context.Context, in *core.DeletePaymentTaskRequest) (*core.DeletePaymentTaskResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	_, err = s.PaymentTaskDB.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	var paymentTask database.PaymentTask
	err = s.PaymentTaskDB.FindOne(ctx, filter).Decode(&paymentTask)
	if err != nil {
		return nil, err
	}
	return &core.DeletePaymentTaskResponse{Status: common.DELETE_STATUS_DELETE_STATUS_SUCCESS, PaymentTask: PaymentTaskDBToPB(paymentTask)}, nil
}

// CreateManyPaymentTask - Insert multiple documents at once in the collection.
func (s CoreServer) CreateManyPaymentTask(ctx context.Context, in *core.CreateManyPaymentTaskRequest) (*core.CreateManyPaymentTaskResponse, error) {
	// Map struct slice to interface slice as InsertMany accepts interface slice as parameter
	insertableList := make([]interface{}, len(in.GetPaymentTasks()))
	for i, v := range in.GetPaymentTasks() {
		insertableList[i] = PaymentTaskPBToDB(v, primitive.NewObjectID())
	}

	// Perform InsertMany operation & validate against the error.
	insertManyResult, err := s.PaymentTaskDB.InsertMany(ctx, insertableList)
	if err != nil {
		return nil, err
	}

	resp := make([]string, len(insertManyResult.InsertedIDs))
	for idx, id := range insertManyResult.InsertedIDs {
		ido := id.(primitive.ObjectID)
		resp[idx] = ido.Hex()
	}

	// Return success without any error.
	return &core.CreateManyPaymentTaskResponse{PaymentTaskIds: resp}, nil
}

// PaymentTaskPBToDB converts a PaymentTask proto object to its serialized DB object
func PaymentTaskPBToDB(paymentTask *common.PaymentTask, id primitive.ObjectID) database.PaymentTask {
	userId, _ := primitive.ObjectIDFromHex(paymentTask.GetUserId())
	accountId, _ := primitive.ObjectIDFromHex(paymentTask.GetAccountId())

	return database.PaymentTask{
		ID:        id,
		UserId:    userId,
		AccountId: accountId,
		Amount:    paymentTask.Amount,
	}
}

// PaymentTaskDBToPB converts a PaymentTask DB object to its proto object
func PaymentTaskDBToPB(paymentTask database.PaymentTask) *common.PaymentTask {
	return &common.PaymentTask{
		PaymentTaskId: paymentTask.ID.Hex(),
		UserId:        paymentTask.UserId.Hex(),
		AccountId:     paymentTask.AccountId.Hex(),
		Amount:        paymentTask.Amount,
	}
}
