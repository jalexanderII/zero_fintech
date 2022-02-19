package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
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

func (p PaymentsServer) GetAccountDetails(ctx context.Context, in *payments.GetAccountDetailsRequest) (*payments.GetAccountDetailsResponse, error) {
	values1 := map[string]string{"access_token": in.GetAccessToken()}
	jsonData1, err := json.Marshal(values1)
	if err != nil {
		log.Fatal(err)
	}
	URL1 := "http://127.0.0.1:8000/api/liabilities/internal"
	resp1, err := http.Post(URL1, "application/json", bytes.NewBuffer(jsonData1))
	if err != nil {
		log.Fatal(err)
	}
	b1, err := io.ReadAll(resp1.Body)
	if err != nil {
		log.Fatalln(err)
	}

	values2 := map[string]string{
		"access_token": in.GetAccessToken(),
		"months":       strconv.FormatInt(in.GetMonths(), 10),
	}
	jsonData2, err := json.Marshal(values2)
	if err != nil {
		log.Fatal(err)
	}
	URL2 := "http://127.0.0.1:8000/api/transactions/internal"
	resp2, err := http.Post(URL2, "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		log.Fatal(err)
	}
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var obj1 LiabilitiesResponse
	var obj2 TransactionsResponse
	if err := json.Unmarshal(b1, &obj1); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b2, &obj2); err != nil {
		panic(err)
	}

	return &payments.GetAccountDetailsResponse{AccountDetailsResponse: PlaidResponseToPB(obj1, obj2, in.GetUser())}, nil
}

func PlaidResponseToPB(lr LiabilitiesResponse, tr TransactionsResponse, user *core.User) *payments.AccountDetailsResponse {
	accountLiabilities := make(map[string]CreditCardLiability)
	for _, al := range lr.Liabilities {
		accountLiabilities[al.AccountId] = al
	}
	accounts := make([]*core.Account, len(tr.Accounts))
	for _, account := range tr.Accounts {
		var acc CreditCardLiability
		if _, ok := accountLiabilities[account.AccountId]; ok {
			acc = accountLiabilities[account.AccountId]
			aprs := make([]*core.AnnualPercentageRates, len(acc.Aprs))
			for _, apr := range acc.Aprs {
				aprs = append(aprs, &core.AnnualPercentageRates{
					AprPercentage:        apr.AprPercentage,
					AprType:              apr.AprType,
					BalanceSubjectToApr:  apr.BalanceSubjectToApr,
					InterestChargeAmount: apr.InterestChargeAmount,
				})
			}
			accounts = append(accounts, &core.Account{
				UserId:                 user.Id,
				Name:                   account.Name,
				OfficialName:           account.OfficialName,
				Type:                   account.Type,
				Subtype:                account.Subtype,
				AvailableBalance:       float64(account.Balances.Available),
				CurrentBalance:         float64(account.Balances.Current),
				CreditLimit:            float64(account.Balances.Limit),
				IsoCurrencyCode:        account.Balances.IsoCurrencyCode,
				AnnualPercentageRate:   aprs,
				IsOverdue:              acc.IsOverdue,
				LastPaymentAmount:      float64(acc.LastPaymentAmount),
				LastStatementIssueDate: acc.LastStatementIssueDate,
				LastStatementBalance:   float64(acc.LastStatementBalance),
				MinimumPaymentAmount:   float64(acc.MinimumPaymentAmount),
				NextPaymentDueDate:     acc.NextPaymentDueDate,
				PlaidAccountId:         account.AccountId,
			})
		}
	}
	transactions := make([]*core.Transaction, len(tr.Transactions))
	for _, transaction := range tr.Transactions {
		transactions = append(transactions, &core.Transaction{
			UserId:               user.Id,
			TransactionType:      *transaction.TransactionType,
			PendingTransactionId: transaction.PendingTransactionId,
			CategoryId:           transaction.CategoryId,
			Category:             transaction.Category,
			TransactionDetails: &core.TransactionDetails{
				Address:         transaction.Location.Address,
				City:            transaction.Location.City,
				State:           transaction.Location.Region,
				Zipcode:         transaction.Location.PostalCode,
				Country:         transaction.Location.Country,
				StoreNumber:     transaction.Location.StoreNumber,
				ReferenceNumber: transaction.PaymentMeta.ReferenceNumber,
			},
			Name:                transaction.Name,
			OriginalDescription: transaction.OriginalDescription,
			Amount:              float64(transaction.Amount),
			IsoCurrencyCode:     transaction.IsoCurrencyCode,
			Date:                transaction.Date,
			Pending:             transaction.Pending,
			MerchantName:        transaction.MerchantName,
			PaymentChannel:      transaction.PaymentChannel,
			AuthorizedDate:      transaction.AuthorizedDate,
			PrimaryCategory:     transaction.PersonalFinanceCategory.Primary,
			DetailedCategory:    transaction.PersonalFinanceCategory.Detailed,
			PlaidAccountId:      transaction.AccountId,
			PlaidTransactionId:  transaction.TransactionId,
		})
	}
	return &payments.AccountDetailsResponse{
		Accounts:     accounts,
		Transactions: transactions,
	}
}
