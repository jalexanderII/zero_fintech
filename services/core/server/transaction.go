package server

import (
	"context"
	"log"
	"time"

	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s CoreServer) CreateTransaction(ctx context.Context, in *core.CreateTransactionRequest) (*core.Transaction, error) {
	transaction := in.GetTransaction()
	newTransaction := TransactionPBToDB(transaction, primitive.NewObjectID())

	_, err := s.TransactionDB.InsertOne(ctx, newTransaction)
	if err != nil {
		log.Printf("Error inserting new Transaction: %v\n", err)
		return nil, err
	}
	return transaction, nil
}

func (s CoreServer) GetTransaction(ctx context.Context, in *core.GetTransactionRequest) (*core.Transaction, error) {
	var transaction database.Transaction
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	err = s.TransactionDB.FindOne(ctx, filter).Decode(&transaction)
	if err != nil {
		return nil, err
	}
	return TransactionDBToPB(transaction), nil
}

func (s CoreServer) ListTransactions(ctx context.Context, in *core.ListTransactionRequest) (*core.ListTransactionResponse, error) {
	var results []database.Transaction
	cursor, err := s.TransactionDB.Find(ctx, bson.D{})
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[TransactionDB] Error getting all users", "error", err)
		return nil, err
	}
	res := make([]*core.Transaction, len(results))
	for idx, transaction := range results {
		res[idx] = TransactionDBToPB(transaction)
	}
	return &core.ListTransactionResponse{Transactions: res}, nil
}

func (s CoreServer) UpdateTransaction(ctx context.Context, in *core.UpdateTransactionRequest) (*core.Transaction, error) {
	transaction := in.GetTransaction()
	name, amount, rewardsEarned := transaction.Name, transaction.Amount, transaction.RewardsEarned
	date := primitive.Timestamp{T: uint32(transaction.Date.AsTime().Unix()), I: 0}
	td := database.TransactionDetails{
		Address:         transaction.GetTransactionDetails().Address,
		DoingBusinessAs: transaction.GetTransactionDetails().DoingBusinessAs,
		DateProcessed:   primitive.Timestamp{T: uint32(transaction.GetTransactionDetails().DateProcessed.AsTime().Unix()), I: 0},
	}

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set",
			bson.D{
				{"name", name}, {"amount", amount},
				{"date", date}, {"rewards_earned", rewardsEarned},
				{"transaction_details", td},
			},
		},
	}
	_, err = s.TransactionDB.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	var tt database.Transaction
	err = s.TransactionDB.FindOne(ctx, filter).Decode(&tt)
	return TransactionDBToPB(tt), nil
}

func (s CoreServer) DeleteTransaction(ctx context.Context, in *core.DeleteTransactionRequest) (*core.DeleteTransactionResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", id}}
	_, err = s.TransactionDB.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	var transaction database.Transaction
	err = s.TransactionDB.FindOne(ctx, filter).Decode(&transaction)
	return &core.DeleteTransactionResponse{Status: core.DELETE_STATUS_DELETE_STATUS_SUCCESS, Transaction: TransactionDBToPB(transaction)}, nil
}

// TransactionPBToDB converts a Transaction proto object to its serialized DB object
func TransactionPBToDB(transaction *core.Transaction, id primitive.ObjectID) database.Transaction {
	userId, _ := primitive.ObjectIDFromHex(transaction.GetUserId())
	accountId, _ := primitive.ObjectIDFromHex(transaction.GetAccountId())

	return database.Transaction{
		ID:            id,
		UserId:        userId,
		AccountId:     accountId,
		Name:          transaction.Name,
		Amount:        transaction.Amount,
		Date:          primitive.Timestamp{T: uint32(transaction.Date.AsTime().Unix()), I: 0},
		RewardsEarned: transaction.RewardsEarned,
		TransactionDetails: database.TransactionDetails{
			Address:         transaction.GetTransactionDetails().Address,
			DoingBusinessAs: transaction.GetTransactionDetails().DoingBusinessAs,
			DateProcessed:   primitive.Timestamp{T: uint32(transaction.GetTransactionDetails().DateProcessed.AsTime().Unix()), I: 0},
		},
	}
}

// TransactionDBToPB converts a Transaction DB object to its proto object
func TransactionDBToPB(transaction database.Transaction) *core.Transaction {
	return &core.Transaction{
		TransactionId:      transaction.ID.Hex(),
		UserId:             transaction.UserId.Hex(),
		AccountId:          transaction.AccountId.Hex(),
		Name:               transaction.Name,
		Amount:             transaction.Amount,
		Date:               timestamppb.New(time.Unix(int64(transaction.Date.T), 0)),
		RewardsEarned:      transaction.RewardsEarned,
		TransactionDetails: &core.TransactionDetails{},
	}
}
