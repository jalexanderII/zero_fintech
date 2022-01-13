package server

import (
	"context"
	"log"
	"time"

	"github.com/jalexanderII/zero_fintech/services/Core/database"
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s CoreServer) CreateAccount(ctx context.Context, in *core.CreateAccountRequest) (*core.Account, error) {
	account := in.GetAccount()
	newAccount := AccountPBToDB(account, primitive.NewObjectID())

	_, err := s.AccountDB.InsertOne(ctx, newAccount)
	if err != nil {
		log.Printf("Error inserting new account: %v\n", err)
		return nil, err
	}
	return account, nil
}

func (s CoreServer) GetAccount(ctx context.Context, in *core.GetAccountRequest) (*core.Account, error) {
	var account database.Account
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	err = s.AccountDB.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return AccountDBToPB(account), nil
}

func (s CoreServer) ListAccounts(ctx context.Context, in *core.ListAccountRequest) (*core.ListAccountResponse, error) {
	var results []database.Account
	cursor, err := s.AccountDB.Find(ctx, bson.D{})
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[AccountDB] Error getting all users", "error", err)
		return nil, err
	}
	res := make([]*core.Account, len(results))
	for idx, account := range results {
		res[idx] = AccountDBToPB(account)
	}
	return &core.ListAccountResponse{Accounts: res}, nil
}

func (s CoreServer) UpdateAccount(ctx context.Context, in *core.UpdateAccountRequest) (*core.Account, error) {
	account := in.GetAccount()
	annualPercentageRate := database.AnnualPercentageRates{
		LowEnd:  account.GetAnnualPercentageRate().LowEnd,
		HighEnd: account.GetAnnualPercentageRate().HighEnd,
	}
	penaltyApr := database.PenaltyAPR{
		PenaltyAPR:    account.GetPenaltyApr().PenaltyApr,
		PenaltyReason: account.GetPenaltyApr().GetPenaltyReason(),
	}
	promotionalRate := database.PromotionalRate{
		TemporaryAPR:   account.GetPromotionalRate().TemporaryApr,
		ExpirationDate: primitive.Timestamp{T: uint32(account.GetPromotionalRate().ExpirationDate.AsTime().Unix()), I: 0},
	}

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set",
			bson.D{
				{"annual_percentage_rate", annualPercentageRate}, {"penalty_apr", penaltyApr},
				{"due_day", account.GetDueDay()}, {"minimum_interest_charge", account.GetMinimumInterestCharge()},
				{"annual_account_fee", account.GetAnnualAccountFee()}, {"foreign_transaction_fee", account.GetForeignTransactionFee()},
				{"promotional_rate", promotionalRate}, {"minimum_payment_due", account.GetMinimumPaymentDue()},
				{"current_balance", account.GetCurrentBalance()}, {"pending_transactions", account.GetPendingTransactions()},
				{"credit_limit", account.GetCreditLimit()},
			},
		},
	}
	_, err = s.AccountDB.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	var aa database.Account
	err = s.AccountDB.FindOne(ctx, filter).Decode(&aa)
	return AccountDBToPB(aa), nil
}

func (s CoreServer) DeleteAccount(ctx context.Context, in *core.DeleteAccountRequest) (*core.DeleteAccountResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", id}}
	_, err = s.AccountDB.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	var account database.Account
	err = s.AccountDB.FindOne(ctx, filter).Decode(&account)
	return &core.DeleteAccountResponse{Status: core.DELETE_STATUS_DELETE_STATUS_SUCCESS, Account: AccountDBToPB(account)}, nil
}

func AccountPBToDB(account *core.Account, id primitive.ObjectID) database.Account {
	userId, _ := primitive.ObjectIDFromHex(account.GetUserId())

	return database.Account{
		ID:        id,
		UserId:    userId,
		Name:      account.Name,
		CreatedAt: primitive.Timestamp{T: uint32(account.CreatedAt.AsTime().Unix()), I: 0},
		AnnualPercentageRate: database.AnnualPercentageRates{
			LowEnd:  account.GetAnnualPercentageRate().LowEnd,
			HighEnd: account.GetAnnualPercentageRate().HighEnd,
		},
		PenaltyAPR: database.PenaltyAPR{
			PenaltyAPR:    account.GetPenaltyApr().PenaltyApr,
			PenaltyReason: account.GetPenaltyApr().GetPenaltyReason(),
		},
		DueDay:                account.DueDay,
		MinimumInterestCharge: account.MinimumInterestCharge,
		AnnualAccountFee:      account.AnnualAccountFee,
		ForeignTransactionFee: account.ForeignTransactionFee,
		PromotionalRate: database.PromotionalRate{
			TemporaryAPR:   account.GetPromotionalRate().TemporaryApr,
			ExpirationDate: primitive.Timestamp{T: uint32(account.GetPromotionalRate().ExpirationDate.AsTime().Unix()), I: 0},
		},
		MinimumPaymentDue:   account.MinimumPaymentDue,
		CurrentBalance:      account.CurrentBalance,
		PendingTransactions: account.PendingTransactions,
		CreditLimit:         account.CreditLimit,
	}
}

func AccountDBToPB(account database.Account) *core.Account {
	return &core.Account{
		AccountId: account.ID.Hex(),
		UserId:    account.UserId.Hex(),
		Name:      account.Name,
		CreatedAt: timestamppb.New(time.Unix(int64(account.CreatedAt.T), 0)),
		AnnualPercentageRate: &core.AnnualPercentageRates{
			LowEnd:  account.AnnualPercentageRate.LowEnd,
			HighEnd: account.AnnualPercentageRate.HighEnd,
		},
		PenaltyApr: &core.PenaltyAPR{
			PenaltyApr:    account.PenaltyAPR.PenaltyAPR,
			PenaltyReason: account.PenaltyAPR.PenaltyReason,
		},
		DueDay:                account.DueDay,
		MinimumInterestCharge: account.MinimumInterestCharge,
		AnnualAccountFee:      account.AnnualAccountFee,
		ForeignTransactionFee: account.ForeignTransactionFee,
		PromotionalRate: &core.PromotionalRate{
			TemporaryApr:   account.PromotionalRate.TemporaryAPR,
			ExpirationDate: timestamppb.New(time.Unix(int64(account.PromotionalRate.ExpirationDate.T), 0)),
		},
		MinimumPaymentDue:   account.MinimumPaymentDue,
		CurrentBalance:      account.CurrentBalance,
		PendingTransactions: account.PendingTransactions,
	}
}
