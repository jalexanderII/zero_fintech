package server

import (
	"context"
	"log"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/core/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	filter := bson.D{{Key: "_id", Value: id}}
	err = s.TransactionDB.FindOne(ctx, filter).Decode(&transaction)
	if err != nil {
		return nil, err
	}
	return TransactionDBToPB(transaction), nil
}

func (s CoreServer) ListTransactions(ctx context.Context, in *core.ListTransactionRequest) (*core.ListTransactionResponse, error) {
	var results []database.Transaction
	cursor, err := s.TransactionDB.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[TransactionDB] Error getting all transactions", "error", err)
		return nil, err
	}
	res := make([]*core.Transaction, len(results))
	for idx, transaction := range results {
		res[idx] = TransactionDBToPB(transaction)
	}
	return &core.ListTransactionResponse{Transactions: res}, nil
}

func (s CoreServer) ListUserTransactions(ctx context.Context, in *core.ListUserTransactionsRequest) (*core.ListTransactionResponse, error) {
	var results []database.Transaction
	id, err := primitive.ObjectIDFromHex(in.GetUserId())
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, err := s.TransactionDB.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[TransactionDB] Error getting all transactions for users", "error", err)
		return nil, err
	}
	res := make([]*core.Transaction, len(results))
	for idx, transaction := range results {
		res[idx] = TransactionDBToPB(transaction)
	}
	return &core.ListTransactionResponse{Transactions: res}, nil
}

// TransactionPBToDB converts a Transaction proto object to its serialized DB object
func TransactionPBToDB(transaction *core.Transaction, id primitive.ObjectID) database.Transaction {
	userId, _ := primitive.ObjectIDFromHex(transaction.GetUserId())
	accountId, _ := primitive.ObjectIDFromHex(transaction.GetAccountId())
	td := transaction.GetTransactionDetails()
	transactionDetails := &database.TransactionDetails{
		Address:         td.Address,
		City:            td.City,
		State:           td.State,
		Zipcode:         td.Zipcode,
		Country:         td.Country,
		StoreNumber:     td.StoreNumber,
		ReferenceNumber: td.ReferenceNumber,
	}
	return database.Transaction{
		ID:                   id,
		UserId:               userId,
		AccountId:            accountId,
		PlaidAccountId:       transaction.PlaidAccountId,
		Name:                 transaction.Name,
		TransactionType:      transaction.TransactionType,
		PendingTransactionId: transaction.PendingTransactionId,
		CategoryId:           transaction.CategoryId,
		Category:             transaction.Category,
		TransactionDetails:   transactionDetails,
		OriginalDescription:  transaction.OriginalDescription,
		Amount:               transaction.Amount,
		IsoCurrencyCode:      transaction.IsoCurrencyCode,
		Date:                 transaction.Date,
		Pending:              transaction.Pending,
		PlaidTransactionId:   transaction.PlaidTransactionId,
		MerchantName:         transaction.MerchantName,
		PaymentChannel:       transaction.PaymentChannel,
		AuthorizedDate:       transaction.AuthorizedDate,
		PrimaryCategory:      transaction.PrimaryCategory,
		DetailedCategory:     transaction.DetailedCategory,
	}
}

// TransactionDBToPB converts a Transaction DB object to its proto object
func TransactionDBToPB(transaction database.Transaction) *core.Transaction {
	td := transaction.TransactionDetails
	transactionDetails := &core.TransactionDetails{
		Address:         td.Address,
		City:            td.City,
		State:           td.State,
		Zipcode:         td.Zipcode,
		Country:         td.Country,
		StoreNumber:     td.StoreNumber,
		ReferenceNumber: td.ReferenceNumber,
	}
	return &core.Transaction{
		TransactionId:        transaction.ID.Hex(),
		UserId:               transaction.UserId.Hex(),
		AccountId:            transaction.AccountId.Hex(),
		PlaidAccountId:       transaction.PlaidAccountId,
		PlaidTransactionId:   transaction.PlaidTransactionId,
		TransactionType:      transaction.TransactionType,
		PendingTransactionId: transaction.PendingTransactionId,
		CategoryId:           transaction.CategoryId,
		Category:             transaction.Category,
		TransactionDetails:   transactionDetails,
		Name:                 transaction.Name,
		OriginalDescription:  transaction.OriginalDescription,
		Amount:               transaction.Amount,
		IsoCurrencyCode:      transaction.IsoCurrencyCode,
		Date:                 transaction.Date,
		Pending:              transaction.Pending,
		MerchantName:         transaction.MerchantName,
		PaymentChannel:       transaction.PaymentChannel,
		AuthorizedDate:       transaction.AuthorizedDate,
		PrimaryCategory:      transaction.PrimaryCategory,
		DetailedCategory:     transaction.DetailedCategory,
	}
}
