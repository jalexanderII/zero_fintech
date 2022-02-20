package server

import (
	"context"
	"testing"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/payments"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"
)

var (
	L        = logrus.New()
	accounts = []*core.Account{
		{
			UserId:           "61df93c0ac601d1be8e64613",
			Name:             "CREDIT CARD",
			OfficialName:     "Chase Sapphire Preferred®",
			Type:             "credit",
			Subtype:          "credit card",
			AvailableBalance: 241.72000122070312,
			CurrentBalance:   15058.2802734375,
			CreditLimit:      15300,
			IsoCurrencyCode:  "USD",
			AnnualPercentageRate: []*core.AnnualPercentageRates{
				{AprPercentage: 24.99, AprType: "cash_apr"},
				{AprPercentage: 22.99, AprType: "purchase_apr", BalanceSubjectToApr: 14973.33},
			},
			LastPaymentAmount:      1000,
			LastStatementIssueDate: "2022-02-06",
			LastStatementBalance:   14973.330078125,
			MinimumPaymentAmount:   437,
			NextPaymentDueDate:     "2022-03-03",
			PlaidAccountId:         "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp",
		},
	}
	transactions = []*core.Transaction{
		{
			UserId:             "61df93c0ac601d1be8e64613",
			TransactionType:    "special",
			CategoryId:         "21010004",
			Category:           []string{"Transfer", "Third Party", "PayPal"},
			TransactionDetails: &core.TransactionDetails{},
			Name:               "FABLETICSLL",
			Amount:             49.95000076293945,
			IsoCurrencyCode:    "USD",
			Date:               "2022-02-08",
			MerchantName:       "Fableticsll",
			PaymentChannel:     "online",
			AuthorizedDate:     "2022-02-07",
			PlaidAccountId:     "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp",
			PlaidTransactionId: "JVP9Z0A9QriaaQJxLeaQTmpjNZbAwDubv3mv0",
		},
		{
			UserId:             "61df93c0ac601d1be8e64613",
			TransactionType:    "place",
			CategoryId:         "18000000",
			Category:           []string{"Service"},
			TransactionDetails: &core.TransactionDetails{},
			Name:               "LEETCODE.COM",
			Amount:             35,
			IsoCurrencyCode:    "USD",
			Date:               "2022-02-08",
			MerchantName:       "Leetcode",
			PaymentChannel:     "online",
			AuthorizedDate:     "2022-02-08",
			PlaidAccountId:     "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp",
			PlaidTransactionId: "xe4pZMBpLkCmm43qxzm4iB5ApL13jECMQpOQO",
		},
		{
			UserId:             "61df93c0ac601d1be8e64613",
			TransactionType:    "special",
			CategoryId:         "15002000",
			Category:           []string{"Interest", "Interest Charged"},
			TransactionDetails: &core.TransactionDetails{},
			Name:               "PURCHASE INTEREST CHARGE",
			Amount:             287.9800109863281,
			IsoCurrencyCode:    "USD",
			Date:               "2022-02-06",
			PaymentChannel:     "other",
			AuthorizedDate:     "2022-02-06",
			PlaidAccountId:     "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp",
			PlaidTransactionId: "8v8KzZ7KrdiyynV45JynuPaLX1oK93FyVkxVv",
		},
		{
			UserId:             "61df93c0ac601d1be8e64613",
			TransactionType:    "place",
			CategoryId:         "18018000",
			Category:           []string{"Service", "Entertainment"},
			TransactionDetails: &core.TransactionDetails{},
			Name:               "Prime Video*WP5AM61P3",
			Amount:             10.989999771118164,
			IsoCurrencyCode:    "USD",
			Date:               "2022-01-27",
			MerchantName:       "Amazon Prime Video",
			PaymentChannel:     "in store",
			AuthorizedDate:     "2022-01-26",
			PlaidAccountId:     "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp",
			PlaidTransactionId: "ZV8qrE3qkJiRRN1v0XRNcXgqw7R9ZpHRLPbLV",
		},
		{
			UserId:             "61df93c0ac601d1be8e64613",
			TransactionType:    "place",
			CategoryId:         "17018000",
			Category:           []string{"Recreation", "Gyms and Fitness Centers"},
			TransactionDetails: &core.TransactionDetails{StoreNumber: "124"},
			Name:               "EQUINOX MOTO #124",
			Amount:             313.5,
			IsoCurrencyCode:    "USD",
			Date:               "2022-01-24",
			MerchantName:       "Equinox Moto",
			PaymentChannel:     "in store",
			AuthorizedDate:     "2022-01-23",
			PlaidAccountId:     "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp",
			PlaidTransactionId: "5KaBgmeBvrhLLOx0rwLOT3ajJDyqoQUBAmwAb",
		},
	}
)

func GenServer() (*PaymentsServer, context.Context) {
	server := NewPaymentsServer(L)
	return server, context.TODO()
}

func TestPaymentsServer_GetAccountDetails(t *testing.T) {
	server, ctx := GenServer()

	accountDetailsResponse := &payments.AccountDetailsResponse{
		Accounts:     accounts,
		Transactions: transactions,
	}

	u := &core.User{
		Id:       "61df93c0ac601d1be8e64613",
		Username: "joel_admin",
		Email:    "fudoshin2596@gmail.com",
		AccountIdToToken: map[string]string{
			"6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp": utils.GetEnv("PLAID_JA_AT"),
		},
	}

	resp, err := server.GetAccountDetails(ctx,
		&payments.GetAccountDetailsRequest{
			User:        u,
			AccessToken: u.AccountIdToToken["6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp"],
			Months:      1,
		})
	if err != nil {
		t.Errorf("1: Error calling plaid API: %v", err)
	}

	// t.Logf("Response: %+v", resp.GetAccountDetailsResponse())
	accResp := resp.GetAccountDetailsResponse().GetAccounts()
	trxnResp := resp.GetAccountDetailsResponse().GetTransactions()

	for idx, acc := range accResp {
		expectedAcc := accountDetailsResponse.Accounts[idx]
		if acc.UserId != expectedAcc.UserId {
			t.Errorf("2: Wrong UserId, expected: %v, got %v", expectedAcc.UserId, acc.UserId)
		}
		if acc.Name != expectedAcc.Name {
			t.Errorf("2: Wrong Name, expected: %v, got %v", expectedAcc.Name, acc.Name)
		}
		if acc.OfficialName != expectedAcc.OfficialName {
			t.Errorf("2: Wrong OfficialName, expected: %v, got %v", expectedAcc.OfficialName, acc.OfficialName)
		}
		if acc.Type != expectedAcc.Type {
			t.Errorf("2: Wrong Type, expected: %v, got %v", expectedAcc.Type, acc.Type)
		}
		if acc.Subtype != expectedAcc.Subtype {
			t.Errorf("2: Wrong Subtype, expected: %v, got %v", expectedAcc.Subtype, acc.Subtype)
		}
		if acc.AvailableBalance != expectedAcc.AvailableBalance {
			t.Errorf("2: Wrong AvailableBalance, expected: %v, got %v", expectedAcc.AvailableBalance, acc.AvailableBalance)
		}
		if acc.CurrentBalance != expectedAcc.CurrentBalance {
			t.Errorf("2: Wrong CurrentBalance, expected: %v, got %v", expectedAcc.CurrentBalance, acc.CurrentBalance)
		}
		if acc.CreditLimit != expectedAcc.CreditLimit {
			t.Errorf("2: Wrong CreditLimit, expected: %v, got %v", expectedAcc.CreditLimit, acc.CreditLimit)
		}
		if acc.IsoCurrencyCode != expectedAcc.IsoCurrencyCode {
			t.Errorf("2: Wrong IsoCurrencyCode, expected: %v, got %v", expectedAcc.IsoCurrencyCode, acc.IsoCurrencyCode)
		}
		if acc.AnnualPercentageRate[0].AprPercentage != expectedAcc.AnnualPercentageRate[0].AprPercentage {
			t.Errorf("2: Wrong AprPercentage, expected: %v, got %v", expectedAcc.AnnualPercentageRate[0].AprPercentage, acc.AnnualPercentageRate[0].AprPercentage)
		}
		if acc.LastPaymentAmount != expectedAcc.LastPaymentAmount {
			t.Errorf("2: Wrong LastPaymentAmount, expected: %v, got %v", expectedAcc.LastPaymentAmount, acc.LastPaymentAmount)
		}
		if acc.LastStatementIssueDate != expectedAcc.LastStatementIssueDate {
			t.Errorf("2: Wrong LastStatementIssueDate, expected: %v, got %v", expectedAcc.LastStatementIssueDate, acc.LastStatementIssueDate)
		}
		if acc.LastStatementBalance != expectedAcc.LastStatementBalance {
			t.Errorf("2: Wrong LastStatementBalance, expected: %v, got %v", expectedAcc.LastStatementBalance, acc.LastStatementBalance)
		}
		if acc.MinimumPaymentAmount != expectedAcc.MinimumPaymentAmount {
			t.Errorf("2: Wrong MinimumPaymentAmount, expected: %v, got %v", expectedAcc.MinimumPaymentAmount, acc.MinimumPaymentAmount)
		}
		if acc.NextPaymentDueDate != expectedAcc.NextPaymentDueDate {
			t.Errorf("2: Wrong NextPaymentDueDate, expected: %v, got %v", expectedAcc.NextPaymentDueDate, acc.NextPaymentDueDate)
		}
		if acc.PlaidAccountId != expectedAcc.PlaidAccountId {
			t.Errorf("2: Wrong PlaidAccountId, expected: %v, got %v", expectedAcc.PlaidAccountId, acc.PlaidAccountId)
		}
	}
	for idx, trxn := range trxnResp {
		expectedTrxn := accountDetailsResponse.Transactions[idx]
		if trxn.UserId != expectedTrxn.UserId {
			t.Errorf("3: Wrong UserId, expected: %v, got %v", expectedTrxn.UserId, trxn.UserId)
		}
		if trxn.TransactionType != expectedTrxn.TransactionType {
			t.Errorf("3: Wrong TransactionType, expected: %v, got %v", expectedTrxn.TransactionType, trxn.TransactionType)
		}
		if trxn.CategoryId != expectedTrxn.CategoryId {
			t.Errorf("3: Wrong CategoryId, expected: %v, got %v", expectedTrxn.CategoryId, trxn.CategoryId)
		}
		if trxn.Category[0] != expectedTrxn.Category[0] {
			t.Errorf("3: Wrong Category, expected: %v, got %v", expectedTrxn.Category, trxn.Category)
		}
		if trxn.Name != expectedTrxn.Name {
			t.Errorf("3: Wrong Name, expected: %v, got %v", expectedTrxn.Name, trxn.Name)
		}
		if trxn.Amount != expectedTrxn.Amount {
			t.Errorf("3: Wrong Amount, expected: %v, got %v", expectedTrxn.Amount, trxn.Amount)
		}
		if trxn.IsoCurrencyCode != expectedTrxn.IsoCurrencyCode {
			t.Errorf("3: Wrong IsoCurrencyCode, expected: %v, got %v", expectedTrxn.IsoCurrencyCode, trxn.IsoCurrencyCode)
		}
		if trxn.Date != expectedTrxn.Date {
			t.Errorf("3: Wrong Date, expected: %v, got %v", expectedTrxn.Date, trxn.Date)
		}
		if trxn.MerchantName != expectedTrxn.MerchantName {
			t.Errorf("3: Wrong MerchantName, expected: %v, got %v", expectedTrxn.MerchantName, trxn.MerchantName)
		}
		if trxn.PaymentChannel != expectedTrxn.PaymentChannel {
			t.Errorf("3: Wrong PaymentChannel, expected: %v, got %v", expectedTrxn.PaymentChannel, trxn.PaymentChannel)
		}
		if trxn.AuthorizedDate != expectedTrxn.AuthorizedDate {
			t.Errorf("3: Wrong AuthorizedDate, expected: %v, got %v", expectedTrxn.AuthorizedDate, trxn.AuthorizedDate)
		}
		if trxn.PlaidAccountId != expectedTrxn.PlaidAccountId {
			t.Errorf("3: Wrong PlaidAccountId, expected: %v, got %v", expectedTrxn.PlaidAccountId, trxn.PlaidAccountId)
		}
		if trxn.PlaidTransactionId != expectedTrxn.PlaidTransactionId {
			t.Errorf("3: Wrong PlaidTransactionId, expected: %v, got %v", expectedTrxn.PlaidTransactionId, trxn.PlaidTransactionId)
		}
	}
}
