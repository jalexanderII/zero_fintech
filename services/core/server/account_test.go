package server

import (
	"testing"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/core/data"
)

func TestCoreServer_CreateManyAccounts(t *testing.T) {
	server, ctx := GenServer()
	limit := 10
	c := 0
	for ok := true; ok; ok = c < limit {
		fake, _ := data.GenFakeAccount()
		_, err := server.AccountDB.InsertOne(ctx, data.FakeAccountToDB(fake))
		if err != nil {
			t.Errorf("1: Error creating new account: %v", err)
		}
		c++
	}
}

func TestCoreServer_CreateAccount(t *testing.T) {
	server, ctx := GenServer()

	acc := &core.Account{
		UserId:           "6212a0101fca9390a37a32d2",
		PlaidAccountId:   "2",
		Name:             "X1",
		OfficialName:     "X1",
		Type:             "credit",
		Subtype:          "credit_card",
		AvailableBalance: 23000,
		CurrentBalance:   2000,
		CreditLimit:      25000,
		IsoCurrencyCode:  "USD",
		AnnualPercentageRate: []*core.AnnualPercentageRates{
			{
				AprPercentage:        22.99,
				AprType:              "cash",
				BalanceSubjectToApr:  25000,
				InterestChargeAmount: 250,
			},
			{
				AprPercentage:        2.49,
				AprType:              "penalty",
				BalanceSubjectToApr:  2000,
				InterestChargeAmount: 100,
			},
		},
		IsOverdue:              false,
		LastPaymentAmount:      0,
		LastStatementIssueDate: "2022-02-20",
		LastStatementBalance:   2000,
		MinimumPaymentAmount:   500,
		NextPaymentDueDate:     "2022-03-01",
	}
	account, err := server.CreateAccount(ctx, &core.CreateAccountRequest{Account: acc})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(account.AnnualPercentageRate) != 2 {
		t.Errorf("2: Failed to fetch correct account: %+v", account)
	}
}

func TestCoreServer_GetAccount(t *testing.T) {
	server, ctx := GenServer()

	account, err := server.GetAccount(ctx, &core.GetAccountRequest{Id: "6212a29794c88ffb3de9d762"})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if account.Name != "Crate-and-Barrel" {
		t.Errorf("2: Failed to fetch correct account: %+v", account)
	}
}

func TestCoreServer_ListAccounts(t *testing.T) {
	server, ctx := GenServer()

	accounts, err := server.ListAccounts(ctx, &core.ListAccountRequest{})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(accounts.Accounts) != 11 {
		t.Errorf("2: Failed to fetch all accounts: %+v", len(accounts.Accounts))
	}
}
