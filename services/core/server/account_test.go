package server

import (
	"testing"
	"time"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExampleAccounts struct {
	Accounts []*core.Account
}

var exampleAccounts = ExampleAccounts{Accounts: []*core.Account{
	{
		UserId:    "61df93c0ac601d1be8e64613",
		Name:      "Chase",
		CreatedAt: timestamppb.New(time.Now()),
		AnnualPercentageRate: &core.AnnualPercentageRates{
			LowEnd:  0.099,
			HighEnd: 0.296,
		},
		PenaltyApr: &core.PenaltyAPR{
			PenaltyApr:    0.40,
			PenaltyReason: core.PenaltyAPR_PENALTY_REASON_LATE_PAYMENT,
		},
		DueDay:                5,
		MinimumInterestCharge: 1.5,
		AnnualAccountFee:      500,
		ForeignTransactionFee: 0.15,
		PromotionalRate: &core.PromotionalRate{
			TemporaryApr:   0.19,
			ExpirationDate: timestamppb.New(time.Now()),
		},
		MinimumPaymentDue:   250,
		CurrentBalance:      9000,
		PendingTransactions: 23.65,
		CreditLimit:         25000,
	},
	{
		UserId:    "61df93c0ac601d1be8e64613",
		Name:      "Amex",
		CreatedAt: timestamppb.New(time.Now()),
		AnnualPercentageRate: &core.AnnualPercentageRates{
			LowEnd:  0.004,
			HighEnd: 0.196,
		},
		PenaltyApr: &core.PenaltyAPR{
			PenaltyApr:    0.530,
			PenaltyReason: core.PenaltyAPR_PENALTY_REASON_LATE_PAYMENT,
		},
		DueDay:                1,
		MinimumInterestCharge: 1.5,
		AnnualAccountFee:      500,
		ForeignTransactionFee: 0.15,
		PromotionalRate: &core.PromotionalRate{
			TemporaryApr:   0.19,
			ExpirationDate: timestamppb.New(time.Now()),
		},
		MinimumPaymentDue:   300,
		CurrentBalance:      10000,
		PendingTransactions: 0,
		CreditLimit:         15000,
	},
	{
		UserId:    "61df93c0ac601d1be8e64613",
		Name:      "X1",
		CreatedAt: timestamppb.New(time.Now()),
		AnnualPercentageRate: &core.AnnualPercentageRates{
			LowEnd:  0.029,
			HighEnd: 0.29,
		},
		PenaltyApr: &core.PenaltyAPR{
			PenaltyApr:    0.0,
			PenaltyReason: core.PenaltyAPR_PENALTY_REASON_LATE_PAYMENT,
		},
		DueDay:                15,
		MinimumInterestCharge: 1,
		AnnualAccountFee:      0,
		ForeignTransactionFee: 0,
		PromotionalRate:       &core.PromotionalRate{},
		MinimumPaymentDue:     100,
		CurrentBalance:        0,
		PendingTransactions:   0,
		CreditLimit:           25000,
	},
}}

func TestCoreServer_CreateAccount(t *testing.T) {
	server, ctx := GenServer()
	account := exampleAccounts.Accounts[2]
	createdAccount, err := server.CreateAccount(ctx, &core.CreateAccountRequest{Account: account})
	if err != nil {
		t.Errorf("1: Error creating new account: %v", err)
	}
	account.AccountId = createdAccount.AccountId
	if createdAccount != account {
		t.Errorf("2: Error creating account with correct details: %v", createdAccount)
	}
}

func TestCoreServer_GetAccount(t *testing.T) {
	server, ctx := GenServer()

	account, err := server.GetAccount(ctx, &core.GetAccountRequest{Id: "61df9e1211f7c7c4f0a2bbcf"})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if account.Name != exampleAccounts.Accounts[2].Name {
		t.Errorf("2: Failed to fetch correct account: %+v", account)
	}
}

func TestCoreServer_ListAccounts(t *testing.T) {
	server, ctx := GenServer()

	accounts, err := server.ListAccounts(ctx, &core.ListAccountRequest{})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(accounts.Accounts) < 1 {
		t.Errorf("2: Failed to fetch realtors: %+v", accounts.Accounts[0])
	}
}

func TestCoreServer_UpdateAccount(t *testing.T) {
	server, ctx := GenServer()

	u := &core.Account{
		UserId:    "61df93c0ac601d1be8e64613",
		Name:      "Barclay",
		CreatedAt: timestamppb.New(time.Now()),
		AnnualPercentageRate: &core.AnnualPercentageRates{
			LowEnd:  0.23,
			HighEnd: 0.45,
		},
		PenaltyApr: &core.PenaltyAPR{
			PenaltyApr:    0.67,
			PenaltyReason: core.PenaltyAPR_PENALTY_REASON_LATE_PAYMENT,
		},
		DueDay:                2,
		MinimumInterestCharge: 10,
		AnnualAccountFee:      0,
		ForeignTransactionFee: 0.01,
		PromotionalRate:       &core.PromotionalRate{},
		MinimumPaymentDue:     10,
		CurrentBalance:        345.23,
		PendingTransactions:   123.99,
		CreditLimit:           5000,
	}

	account, err := server.UpdateAccount(ctx, &core.UpdateAccountRequest{Id: "61df9f3397fa9f3b7a9b67a8", Account: u})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if account.CreditLimit != u.CreditLimit {
		t.Errorf("2: Failed to fetch correct account: %+v", account)
	}
}

func TestCoreServer_DeleteAccount(t *testing.T) {
	server, ctx := GenServer()

	u, _ := server.ListAccounts(ctx, &core.ListAccountRequest{})
	originalLen := len(u.GetAccounts())

	newAccount := AccountPBToDB(
		&core.Account{
			UserId:    "61df93c0ac601d1be8e64613",
			Name:      "TO_DELETE",
			CreatedAt: timestamppb.New(time.Now()),
			AnnualPercentageRate: &core.AnnualPercentageRates{
				LowEnd:  0.23,
				HighEnd: 0.45,
			},
			PenaltyApr: &core.PenaltyAPR{
				PenaltyApr:    0.67,
				PenaltyReason: core.PenaltyAPR_PENALTY_REASON_LATE_PAYMENT,
			},
			DueDay:                2,
			MinimumInterestCharge: 10,
			AnnualAccountFee:      0,
			ForeignTransactionFee: 0.01,
			PromotionalRate:       &core.PromotionalRate{},
			MinimumPaymentDue:     10,
			CurrentBalance:        345.23,
			PendingTransactions:   123.99,
			CreditLimit:           5000,
		},
		primitive.NewObjectID(),
	)
	_, err := server.AccountDB.InsertOne(ctx, &newAccount)
	if err != nil {
		t.Errorf("1: Error creating new account:: %v", err)
	}

	accounts, err := server.ListAccounts(ctx, &core.ListAccountRequest{})
	if err != nil {
		t.Errorf("2: An error was returned: %v", err)
	}
	newLen := len(accounts.GetAccounts())
	if newLen != originalLen+1 {
		t.Errorf("3: An error adding a temp account, number of accounts in DB: %v", newLen)
	}

	deleted, err := server.DeleteAccount(ctx, &core.DeleteAccountRequest{Id: newAccount.ID.Hex()})
	if err != nil {
		t.Errorf("4: An error was returned: %v", err)
	}
	if deleted.Status != common.DELETE_STATUS_DELETE_STATUS_SUCCESS {
		t.Errorf("5: Failed to delete account: %+v\n, %+v", deleted.Status, deleted.GetAccount())
	}
}
