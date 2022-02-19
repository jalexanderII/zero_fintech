package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/payments"
	"github.com/jalexanderII/zero_fintech/utils"
	"github.com/sirupsen/logrus"
)

var L = logrus.New()

func GenServer() (*PaymentsServer, context.Context) {
	server := NewPaymentsServer(L)
	return server, context.TODO()
}

func TestPaymentsServer_GetLiabilities(t *testing.T) {
	server, ctx := GenServer()
	u := &core.User{
		Id:       "61df93c0ac601d1be8e64613",
		Username: "joel_admin",
		Email:    "fudoshin2596@gmail.com",
	}

	resp, err := server.GetLiabilities(ctx,
		&payments.GetLiabilitiesRequest{
			User:        u,
			AccessToken: utils.GetEnv("PLAID_JA_AT"),
		})
	if err != nil {
		t.Errorf("1: Error calling plaid API: %v", err)
	}

	// t.Logf(string(resp.GetLiabilitiesGetResponse()))

	var obj LiabilitiesResponse
	if err := json.Unmarshal(resp.GetLiabilitiesGetResponse(), &obj); err != nil {
		panic(err)
	}

	if obj.Liabilities == nil {
		t.Errorf("2: Error calling plaid API : %v", obj.Liabilities)
	}
	credit := obj.Liabilities[0]
	if credit.AccountId != "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp" {
		t.Errorf("3: Error calling plaid API : %v", credit)
	}
	if len(credit.Aprs) != 2 {
		t.Errorf("4: Error calling plaid API : %v", credit)
	}
}

func TestPaymentsServer_GetTransactions(t *testing.T) {
	server, ctx := GenServer()
	u := &core.User{
		Id:       "61df93c0ac601d1be8e64613",
		Username: "joel_admin",
		Email:    "fudoshin2596@gmail.com",
	}

	resp, err := server.GetTransactions(ctx,
		&payments.GetTransactionsRequest{
			User:        u,
			AccessToken: utils.GetEnv("PLAID_JA_AT"),
			Months:      6,
		})
	if err != nil {
		t.Errorf("1: Error calling plaid API: %v", err)
	}

	// t.Logf(string(resp.GetTransactionsGetResponse()))

	var obj TransactionsResponse
	if err := json.Unmarshal(resp.GetTransactionsGetResponse(), &obj); err != nil {
		panic(err)
	}

	if obj.Accounts == nil || obj.Transactions == nil {
		t.Errorf("2: Error calling plaid API : %v", obj)
	}

	account := obj.Accounts[0]
	if account.Type != "credit" {
		t.Errorf("3: Account has wrong type : %v", account)
	}

	for _, trxn := range obj.Transactions {
		if trxn.AccountId != "6kx3QJN30rcPP6wjeZP6s4KB1PmMp9UaymBKp" {
			t.Errorf("4: Trasactions are for wrong account : %v", trxn)
		}
	}
}
