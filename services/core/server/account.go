package server

import (
	"context"
	"log"

	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s CoreServer) CreateAccount(ctx context.Context, in *core.CreateAccountRequest) (*core.Account, error) {
	account := in.GetAccount()
	newAccount := AccountPBToDB(account, primitive.NewObjectID())

	dbAccount, err := s.AccountDB.InsertOne(ctx, newAccount)
	if err != nil {
		log.Printf("Error inserting new account: %v\n", err)
		return nil, err
	}
	if oid, ok := dbAccount.InsertedID.(primitive.ObjectID); ok {
		account.AccountId = oid.Hex()
	}

	return account, nil
}

func (s CoreServer) GetAccount(ctx context.Context, in *core.GetAccountRequest) (*core.Account, error) {
	var account database.Account
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	err = s.AccountDB.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return AccountDBToPB(account), nil
}

func (s CoreServer) ListAccounts(ctx context.Context, in *core.ListAccountRequest) (*core.ListAccountResponse, error) {
	var results []database.Account
	cursor, err := s.AccountDB.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[AccountDB] Error getting all accounts", "error", err)
		return nil, err
	}
	res := make([]*core.Account, len(results))
	for idx, account := range results {
		res[idx] = AccountDBToPB(account)
	}
	return &core.ListAccountResponse{Accounts: res}, nil
}

func (s CoreServer) ListUserAccounts(ctx context.Context, in *core.ListUserAccountsRequest) (*core.ListAccountResponse, error) {
	var results []database.Account
	id, err := primitive.ObjectIDFromHex(in.GetUserId())
	if err != nil {
		return nil, err
	}
	s.l.Info("id:", id.Hex())

	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, err := s.AccountDB.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	s.l.Info("Cursor:", cursor)
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[AccountDB] Error getting all accounts for user", "error", err)
		return nil, err
	}
	s.l.Info("Cursor results:", results)
	res := make([]*core.Account, len(results))
	for idx, account := range results {
		res[idx] = AccountDBToPB(account)
	}
	return &core.ListAccountResponse{Accounts: res}, nil
}

// AccountPBToDB converts an Account proto object to its serialized DB object
func AccountPBToDB(account *core.Account, id primitive.ObjectID) database.Account {
	userId, _ := primitive.ObjectIDFromHex(account.GetUserId())
	var aprs []*database.AnnualPercentageRates
	for _, apr := range account.AnnualPercentageRate {
		aprs = append(aprs, &database.AnnualPercentageRates{
			AprPercentage:        apr.AprPercentage,
			AprType:              apr.AprType,
			BalanceSubjectToApr:  apr.BalanceSubjectToApr,
			InterestChargeAmount: apr.InterestChargeAmount,
		})
	}

	return database.Account{
		ID:                     id,
		UserId:                 userId,
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

// AccountDBToPB converts an Account DB object to its proto object
func AccountDBToPB(account database.Account) *core.Account {
	aprs := make([]*core.AnnualPercentageRates, len(account.AnnualPercentageRate))
	for _, apr := range account.AnnualPercentageRate {

		aprs = append(aprs, &core.AnnualPercentageRates{
			AprPercentage:        apr.AprPercentage,
			AprType:              apr.AprType,
			BalanceSubjectToApr:  apr.BalanceSubjectToApr,
			InterestChargeAmount: apr.InterestChargeAmount,
		})
	}
	return &core.Account{
		AccountId:              account.ID.Hex(),
		PlaidAccountId:         account.PlaidAccountId,
		UserId:                 account.UserId.Hex(),
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
