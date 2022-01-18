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

func (s CoreServer) CreatePaymentTask(ctx context.Context, in *core.CreatePaymentTaskRequest) (*core.PaymentTask, error) {
	paymentTask := in.GetPaymentTask()
	newPaymentTask := PaymentTaskPBToDB(paymentTask, primitive.NewObjectID())

	_, err := s.PaymentTaskDB.InsertOne(ctx, newPaymentTask)
	if err != nil {
		log.Printf("Error inserting new PaymentTask: %v\n", err)
		return nil, err
	}
	return paymentTask, nil
}

func (s CoreServer) GetPaymentTask(ctx context.Context, in *core.GetPaymentTaskRequest) (*core.PaymentTask, error) {
	var paymentTask database.PaymentTask
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	err = s.PaymentTaskDB.FindOne(ctx, filter).Decode(&paymentTask)
	if err != nil {
		return nil, err
	}
	return PaymentTaskDBToPB(paymentTask), nil
}

func (s CoreServer) ListPaymentTasks(ctx context.Context, in *core.ListPaymentTaskRequest) (*core.ListPaymentTaskResponse, error) {
	var results []database.PaymentTask
	cursor, err := s.PaymentTaskDB.Find(ctx, bson.D{})
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[PaymentTaskDB] Error getting all users", "error", err)
		return nil, err
	}
	res := make([]*core.PaymentTask, len(results))
	for idx, paymentTask := range results {
		res[idx] = PaymentTaskDBToPB(paymentTask)
	}
	return &core.ListPaymentTaskResponse{PaymentTasks: res}, nil
}

func (s CoreServer) UpdatePaymentTask(ctx context.Context, in *core.UpdatePaymentTaskRequest) (*core.PaymentTask, error) {
	paymentTask := in.GetPaymentTask()
	metaData := database.MetaData{
		PreferredPlanType:    paymentTask.GetMetaData().GetPreferredPlanType(),
		PreferredPaymentFreq: paymentTask.GetMetaData().GetPreferredPaymentFreq(),
	}
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"amount", paymentTask.Amount}, {"meta_data", metaData}}}}
	_, err = s.PaymentTaskDB.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	var payment_task database.PaymentTask
	err = s.PaymentTaskDB.FindOne(ctx, filter).Decode(&payment_task)
	return PaymentTaskDBToPB(payment_task), nil
}
func (s CoreServer) DeletePaymentTask(ctx context.Context, in *core.DeletePaymentTaskRequest) (*core.DeletePaymentTaskResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", id}}
	_, err = s.PaymentTaskDB.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	var paymentTask database.PaymentTask
	err = s.PaymentTaskDB.FindOne(ctx, filter).Decode(&paymentTask)
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
	_, err := s.PaymentTaskDB.InsertMany(ctx, insertableList)
	if err != nil {
		return nil, err
	}
	// Return success without any error.
	return &core.CreateManyPaymentTaskResponse{}, nil
}

// PaymentTaskPBToDB converts a PaymentTask proto object to its serialized DB object
func PaymentTaskPBToDB(paymentTask *core.PaymentTask, id primitive.ObjectID) database.PaymentTask {
	userId, _ := primitive.ObjectIDFromHex(paymentTask.GetUserId())
	accountId, _ := primitive.ObjectIDFromHex(paymentTask.GetAccountId())
	transactionId, _ := primitive.ObjectIDFromHex(paymentTask.GetTransactionId())

	return database.PaymentTask{
		ID:            id,
		UserId:        userId,
		TransactionId: transactionId,
		AccountId:     accountId,
		Amount:        paymentTask.Amount,
		MetaData: database.MetaData{
			PreferredPlanType:    paymentTask.GetMetaData().GetPreferredPlanType(),
			PreferredPaymentFreq: paymentTask.GetMetaData().GetPreferredPaymentFreq(),
		},
	}
}

// PaymentTaskDBToPB converts a PaymentTask DB object to its proto object
func PaymentTaskDBToPB(paymentTask database.PaymentTask) *core.PaymentTask {
	return &core.PaymentTask{
		PaymentTaskId: paymentTask.ID.Hex(),
		UserId:        paymentTask.UserId.Hex(),
		TransactionId: paymentTask.TransactionId.Hex(),
		AccountId:     paymentTask.AccountId.Hex(),
		Amount:        paymentTask.Amount,
		MetaData: &core.MetaData{
			PreferredPlanType:    paymentTask.MetaData.PreferredPlanType,
			PreferredPaymentFreq: paymentTask.MetaData.PreferredPaymentFreq,
		},
	}
}
