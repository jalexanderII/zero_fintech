package data

import (
	"reflect"

	"github.com/bxcodec/faker/v3"
	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var L = logrus.New()

func CustomGenerator() {
	_ = faker.AddProvider("preferred_plan_type", func(v reflect.Value) (interface{}, error) {
		return common.PlanType_PLAN_TYPE_OPTIM_CREDIT_SCORE, nil
	})
	_ = faker.AddProvider("preferred_payment_freq", func(v reflect.Value) (interface{}, error) {
		return common.PaymentFrequency_PAYMENT_FREQUENCY_MONTHLY, nil
	})
	_ = faker.AddProvider("account_id_to_token", func(v reflect.Value) (interface{}, error) {
		return map[string]string{"1": "2"}, nil
	})
}

func GenFakePaymentTask() (*common.PaymentTask, error) {
	CustomGenerator()
	var fake common.PaymentTask
	err := faker.FakeData(&fake)
	if err != nil {
		L.Error("[Error] Could not fake this object", "error", err)
		return nil, err
	}
	return &fake, nil
}

func GenFakeAccount() (*FakeAccount, error) {
	_ = faker.SetRandomMapAndSliceSize(2)
	CustomGenerator()
	var fake FakeAccount
	err := faker.FakeData(&fake)
	if err != nil {
		L.Error("[Error] Could not fake this object", "error", err)
		return nil, err
	}
	return &fake, nil
}

func GenFakeUser() (*FakeUser, error) {
	_ = faker.SetRandomMapAndSliceSize(2)
	CustomGenerator()
	var fake FakeUser
	err := faker.FakeData(&fake)
	if err != nil {
		L.Error("[Error] Could not fake this object", "error", err)
		return nil, err
	}
	return &fake, nil
}

func GenFakeTransaction() (*FakeTransaction, error) {
	_ = faker.SetRandomMapAndSliceSize(2)
	CustomGenerator()
	var fake FakeTransaction
	err := faker.FakeData(&fake)
	if err != nil {
		L.Error("[Error] Could not fake this object", "error", err)
		return nil, err
	}
	return &fake, nil
}

func FakeAccountToDB(account *FakeAccount) database.Account {
	aprs := make([]*database.AnnualPercentageRates, len(account.AnnualPercentageRate))
	for _, apr := range account.AnnualPercentageRate {
		aprs = append(aprs, &database.AnnualPercentageRates{
			AprPercentage:        apr.AprPercentage,
			AprType:              apr.AprType,
			BalanceSubjectToApr:  apr.BalanceSubjectToApr,
			InterestChargeAmount: apr.InterestChargeAmount,
		})
	}
	uid, _ := primitive.ObjectIDFromHex(account.UserId)

	return database.Account{
		ID:                     primitive.NewObjectID(),
		UserId:                 uid,
		PlaidAccountId:         account.PlaidAccountId,
		Name:                   account.Name,
		OfficialName:           account.OfficialName,
		Type:                   account.Type,
		Subtype:                account.Subtype,
		AvailableBalance:       account.AvailableBalance,
		CurrentBalance:         account.CurrentBalance,
		CreditLimit:            account.CreditLimit,
		IsoCurrencyCode:        account.IsoCurrencyCode,
		AnnualPercentageRate:   aprs,
		IsOverdue:              account.IsOverdue,
		LastPaymentAmount:      account.LastPaymentAmount,
		LastStatementIssueDate: account.LastStatementIssueDate,
		LastStatementBalance:   account.LastStatementBalance,
		MinimumPaymentAmount:   account.MinimumPaymentAmount,
		NextPaymentDueDate:     account.NextPaymentDueDate,
	}
}

func FakeTransactionToDB(transaction *FakeTransaction) database.Transaction {
	userId, _ := primitive.ObjectIDFromHex(transaction.UserId)
	accountId, _ := primitive.ObjectIDFromHex(transaction.AccountId)
	td := transaction.TransactionDetails
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
		ID:                   primitive.NewObjectID(),
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
