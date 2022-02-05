package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jalexanderII/zero_fintech/gen/Go/payments"
	"github.com/sirupsen/logrus"
)

type PaymentsServer struct {
	payments.UnimplementedPlaidServer
	// custom logger
	l *logrus.Logger
}

func NewPaymentsServer(l *logrus.Logger) *PaymentsServer {
	return &PaymentsServer{l: l}
}

func (p PaymentsServer) GetLiabilities(ctx context.Context, in *payments.GetLiabilitiesRequest) (*payments.GetLiabilitiesResponse, error) {
	values := map[string]string{"access_token": in.GetAccessToken()}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}
	URL := "http://127.0.0.1:8000/api/liabilities/internal"
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// var obj LiabilitiesResponse
	// if err := json.Unmarshal(b, &obj); err != nil {
	// 	panic(err)
	// }
	// p.l.Info(obj)
	return &payments.GetLiabilitiesResponse{LiabilitiesGetResponse: b}, nil
}

type LiabilitiesResponse struct {
	Liabilities struct {
		Credit []struct {
			AccountId string `json:"account_id,omitempty"`
			Aprs      []struct {
				AprPercentage        float64     `json:"apr_percentage,omitempty"`
				AprType              string      `json:"apr_type,omitempty"`
				BalanceSubjectToApr  float64     `json:"balance_subject_to_apr,omitempty"`
				InterestChargeAmount interface{} `json:"interest_charge_amount,omitempty"`
			} `json:"aprs"`
			IsOverdue              interface{} `json:"is_overdue,omitempty"`
			LastPaymentAmount      int         `json:"last_payment_amount,omitempty"`
			LastPaymentDate        string      `json:"last_payment_date,omitempty"`
			LastStatementBalance   float64     `json:"last_statement_balance,omitempty"`
			LastStatementIssueDate string      `json:"last_statement_issue_date,omitempty"`
			MinimumPaymentAmount   int         `json:"minimum_payment_amount,omitempty"`
			NextPaymentDueDate     string      `json:"next_payment_due_date,omitempty"`
		} `json:"credit,omitempty"`
	} `json:"liabilities,omitempty"`
}
