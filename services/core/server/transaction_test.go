package server

import (
	"testing"
	"time"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExampleTransactions struct {
	Transactions []*core.Transaction
}

var exampleTransactions = ExampleTransactions{Transactions: []*core.Transaction{
	{
		UserId:        "61df93c0ac601d1be8e64613",
		AccountId:     "61df9b621d2c2b15a6e53ec9",
		Name:          "Equinox",
		Amount:        325,
		Date:          timestamppb.New(time.Now()),
		RewardsEarned: 3,
		TransactionDetails: &core.TransactionDetails{
			Address:         "Prince Street, New York NY",
			DoingBusinessAs: "Equinox Prince Street",
			DateProcessed:   timestamppb.New(time.Now()),
		},
	},
	{
		UserId:        "61df93c0ac601d1be8e64613",
		AccountId:     "61df9b621d2c2b15a6e53ec9",
		Name:          "TST* ACME 00033545 NEW YORK NY",
		Amount:        87,
		Date:          timestamppb.New(time.Now()),
		RewardsEarned: 87,
		TransactionDetails: &core.TransactionDetails{
			Address:         "9 GREAT JONES ST NEW YORK, NY 10012 USA",
			DoingBusinessAs: "TST* ACME",
			DateProcessed:   timestamppb.New(time.Now()),
		},
	},
}}

func TestCoreServer_CreateTransaction(t *testing.T) {
	server, ctx := GenServer()

	amexTransaction := exampleTransactions.Transactions[0]

	transaction, err := server.CreateTransaction(ctx, &core.CreateTransactionRequest{Transaction: amexTransaction})
	if err != nil {
		t.Errorf("1: Error creating new transaction: %v", err)
	}
	amexTransaction.TransactionId = transaction.TransactionId
	if transaction != amexTransaction {
		t.Errorf("2: Error creating account with correct details: %v", transaction)
	}
}

func TestCoreServer_GetTransaction(t *testing.T) {
	server, ctx := GenServer()

	transaction, err := server.GetTransaction(ctx, &core.GetTransactionRequest{Id: "61dfa20adebb9d4fb62b9703"})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if transaction.Name != "Equinox" {
		t.Errorf("2: Failed to fetch correct transaction: %+v", transaction)
	}
}

func TestCoreServer_ListTransactions(t *testing.T) {
	server, ctx := GenServer()

	transactions, err := server.ListTransactions(ctx, &core.ListTransactionRequest{})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(transactions.Transactions) < 1 {
		t.Errorf("2: Failed to fetch realtors: %+v", transactions.Transactions[0])
	}

}

func TestCoreServer_UpdateTransaction(t *testing.T) {
	server, ctx := GenServer()

	u := &core.Transaction{
		UserId:        "61df93c0ac601d1be8e64613",
		AccountId:     "61df9b621d2c2b15a6e53ec9",
		Name:          "Acme",
		Amount:        87,
		Date:          timestamppb.New(time.Now()),
		RewardsEarned: 87,
		TransactionDetails: &core.TransactionDetails{
			Address:         "9 GREAT JONES ST NEW YORK, NY 10012 USA",
			DoingBusinessAs: "TST* ACME",
			DateProcessed:   timestamppb.New(time.Now()),
		},
	}

	transaction, err := server.UpdateTransaction(ctx, &core.UpdateTransactionRequest{Id: "61dfa3ed434e4e409f054717", Transaction: u})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if transaction.Name != u.Name {
		t.Errorf("2: Failed to fetch correct transaction: %+v", transaction)
	}

}

func TestCoreServer_DeleteTransaction(t *testing.T) {
	server, ctx := GenServer()

	u, err := server.ListTransactions(ctx, &core.ListTransactionRequest{})
	originalLen := len(u.GetTransactions())

	newTransaction := TransactionPBToDB(
		&core.Transaction{
			UserId:        "61df93c0ac601d1be8e64613",
			AccountId:     "61df9b621d2c2b15a6e53ec9",
			Name:          "TO_DELETE",
			Amount:        87,
			Date:          timestamppb.New(time.Now()),
			RewardsEarned: 87,
			TransactionDetails: &core.TransactionDetails{
				Address:         "9 GREAT JONES ST NEW YORK, NY 10012 USA",
				DoingBusinessAs: "TST* ACME",
				DateProcessed:   timestamppb.New(time.Now()),
			},
		},
		primitive.NewObjectID(),
	)
	_, err = server.TransactionDB.InsertOne(ctx, &newTransaction)
	if err != nil {
		t.Errorf("1: Error creating new transaction:: %v", err)
	}

	transactions, err := server.ListTransactions(ctx, &core.ListTransactionRequest{})
	if err != nil {
		t.Errorf("2: An error was returned: %v", err)
	}
	newLen := len(transactions.GetTransactions())
	if newLen != originalLen+1 {
		t.Errorf("3: An error adding a temp transaction, number of transactions in DB: %v", newLen)
	}

	deleted, err := server.DeleteTransaction(ctx, &core.DeleteTransactionRequest{Id: newTransaction.ID.Hex()})
	if err != nil {
		t.Errorf("4: An error was returned: %v", err)
	}
	if deleted.Status != core.DELETE_STATUS_DELETE_STATUS_SUCCESS {
		t.Errorf("5: Failed to delete transaction: %+v\n, %+v", deleted.Status, deleted.GetTransaction())
	}
}
