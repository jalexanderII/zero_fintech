package database

import (
	"fmt"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentTask is a DB Serialization of Proto PaymentTask
type PaymentTask struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	AccountId primitive.ObjectID `bson:"account_id"`
	Amount    float64            `bson:"amount"`
}

// MetaData is a DB Serialization of Proto MetaData
type MetaData struct {
	PreferredPlanType         common.PlanType         `bson:"preferred_plan_type"`
	PreferredTimelineInMonths float64                 `bson:"preferred_timeline_in_months"`
	PreferredPaymentFreq      common.PaymentFrequency `bson:"preferred_payment_freq"`
}

// AnnualPercentageRates is a DB Serialization of Proto AnnualPercentageRates
type AnnualPercentageRates struct {
	AprPercentage        float64 `bson:"apr_percentage"`
	AprType              string  `bson:"apr_type"`
	BalanceSubjectToApr  float64 `bson:"balance_subject_to_apr"`
	InterestChargeAmount float64 `bson:"interest_charge_amount"`
}

// Account is a DB Serialization of Proto Account
type Account struct {
	ID                     primitive.ObjectID       `bson:"_id"`
	PlaidAccountId         string                   `bson:"plaid_account_id"`
	UserId                 primitive.ObjectID       `bson:"user_id"`
	Name                   string                   `bson:"name"`
	OfficialName           string                   `bson:"official_name"`
	Type                   string                   `bson:"type"`
	Subtype                string                   `bson:"subtype"`
	AvailableBalance       float64                  `bson:"available_balance"`
	CurrentBalance         float64                  `bson:"current_balance"`
	CreditLimit            float64                  `bson:"credit_limit"`
	IsoCurrencyCode        string                   `bson:"iso_currency_code"`
	AnnualPercentageRate   []*AnnualPercentageRates `bson:"annual_percentage_rate"`
	IsOverdue              bool                     `bson:"is_overdue"`
	LastPaymentAmount      float64                  `bson:"last_payment_amount"`
	LastStatementIssueDate string                   `bson:"last_statement_issue_date"`
	LastStatementBalance   float64                  `bson:"last_statement_balance"`
	MinimumPaymentAmount   float64                  `bson:"minimum_payment_amount"`
	NextPaymentDueDate     string                   `bson:"next_payment_due_date"`
}

func (a *Account) NotNull() bool {
	return (a != nil && a.PlaidAccountId != "")
}

// Transaction is a DB Serialization of Proto Transaction
type Transaction struct {
	ID                   primitive.ObjectID  `bson:"_id"`
	PlaidTransactionId   string              `bson:"plaid_transaction_id"`
	AccountId            primitive.ObjectID  `bson:"account_id"`
	PlaidAccountId       string              `bson:"plaid_account_id"`
	UserId               primitive.ObjectID  `bson:"user_id"`
	TransactionType      string              `bson:"transaction_type"`
	PendingTransactionId string              `bson:"pending_transaction_id"`
	CategoryId           string              `bson:"category_id"`
	Category             []string            `bson:"category"`
	TransactionDetails   *TransactionDetails `bson:"transaction_details"`
	Name                 string              `bson:"name"`
	OriginalDescription  string              `bson:"original_description"`
	Amount               float64             `bson:"amount"`
	IsoCurrencyCode      string              `bson:"iso_currency_code"`
	Date                 string              `bson:"date"`
	Pending              bool                `bson:"pending"`
	MerchantName         string              `bson:"merchant_name"`
	PaymentChannel       string              `bson:"payment_channel"`
	AuthorizedDate       string              `bson:"authorized_date"`
	PrimaryCategory      string              `bson:"primary_category"`
	DetailedCategory     string              `bson:"detailed_category"`
}

// TransactionDetails is a DB Serialization of Proto TransactionDetails
type TransactionDetails struct {
	Address         string `bson:"address"`
	City            string `bson:"city"`
	State           string `bson:"state"`
	Zipcode         string `bson:"zipcode"`
	Country         string `bson:"country"`
	StoreNumber     string `bson:"store_number"`
	ReferenceNumber string `bson:"reference_number"`
}

// User is a DB Serialization of Proto User
type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Username    string             `bson:"username"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	PhoneNumber string             `bson:"phone_number"`
}

func FormatPhoneNumber(pn string) string {
	// if the first char isn't a plus, add it
	if pn[0:1] != "+" {
		return fmt.Sprintf("+%s", pn)
	}
	return pn
}
