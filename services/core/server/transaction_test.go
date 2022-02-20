package server

import (
	"testing"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/core/data"
)

func TestCoreServer_CreateManyTransactions(t *testing.T) {
	server, ctx := GenServer()
	limit := 50
	c := 0
	for ok := true; ok; ok = c < limit {
		fake, _ := data.GenFakeTransaction()
		_, err := server.TransactionDB.InsertOne(ctx, data.FakeTransactionToDB(fake))
		if err != nil {
			t.Errorf("1: Error creating new account: %v", err)
		}
		c++
	}
}

func TestCoreServer_GetTransaction(t *testing.T) {
	server, ctx := GenServer()

	transaction, err := server.GetTransaction(ctx, &core.GetTransactionRequest{Id: "6212a2f867c13199e5a58412"})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if transaction.MerchantName != "Miss Eudora Parker" {
		t.Errorf("2: Failed to fetch correct transaction: %+v", transaction)
	}
}

func TestCoreServer_ListTransactions(t *testing.T) {
	server, ctx := GenServer()

	transactions, err := server.ListTransactions(ctx, &core.ListTransactionRequest{})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(transactions.Transactions) != 50 {
		t.Errorf("2: Failed to fetch all transactions: %+v", len(transactions.Transactions))
	}
}
