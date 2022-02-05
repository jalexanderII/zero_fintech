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
	var obj LiabilitiesResponse
	if err := json.Unmarshal(resp.GetLiabilitiesGetResponse(), &obj); err != nil {
		panic(err)
	}

	if obj.Liabilities.Credit == nil {
		t.Errorf("2: Error calling plaid API : %v", obj.Liabilities)
	}
	credit := obj.Liabilities.Credit[0]
	if credit.AccountId != "dROgn7DjN0hMjE8qOpedfvjV9w3KrvHb8LwzY" {
		t.Errorf("3: Error calling plaid API : %v", credit)
	}
	if len(credit.Aprs) != 2 {
		t.Errorf("4: Error calling plaid API : %v", credit)
	}
}
