package database

import (
	"time"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
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
	PreferredPlanType         common.PlanType         `bson:"preferred_plan_type" faker:"preferred_plan_type" `
	PreferredTimelineInMonths float64                 `bson:"preferred_timeline_in_months" faker:"preferred_timeline_in_months" `
	PreferredPaymentFreq      common.PaymentFrequency `bson:"preferred_payment_freq" faker:"preferred_payment_freq" `
}

// AnnualPercentageRates is a DB Serialization of Proto AnnualPercentageRates
type AnnualPercentageRates struct {
	LowEnd  float64 `bson:"low_end"`
	HighEnd float64 `bson:"high_end"`
}

// PenaltyAPR is a DB Serialization of Proto PenaltyAPR
type PenaltyAPR struct {
	PenaltyAPR    float64                       `bson:"penalty_apr"`
	PenaltyReason core.PenaltyAPR_PenaltyReason `bson:"penalty_reason" faker:"penalty_reason"`
}

// PromotionalRate is a DB Serialization of Proto PromotionalRate
type PromotionalRate struct {
	TemporaryAPR   float64   `bson:"temporary_apr"`
	ExpirationDate time.Time `bson:"expiration_date"`
}

// Account is a DB Serialization of Proto Account
type Account struct {
	ID                    primitive.ObjectID    `bson:"_id"`
	UserId                primitive.ObjectID    `bson:"user_id"`
	Name                  string                `bson:"name"`
	CreatedAt             time.Time             `bson:"created_at"`
	AnnualPercentageRate  AnnualPercentageRates `bson:"annual_percentage_rate"`
	PenaltyAPR            PenaltyAPR            `bson:"penalty_apr"`
	DueDay                int32                 `bson:"due_day"`
	MinimumInterestCharge float64               `bson:"minimum_interest_charge"`
	AnnualAccountFee      float64               `bson:"annual_account_fee"`
	ForeignTransactionFee float64               `bson:"foreign_transaction_fee"`
	PromotionalRate       PromotionalRate       `bson:"promotional_rate"`
	MinimumPaymentDue     float64               `bson:"minimum_payment_due"`
	CurrentBalance        float64               `bson:"current_balance"`
	PendingTransactions   float64               `bson:"pending_transactions"`
	CreditLimit           float64               `bson:"credit_limit"`
}

// TransactionDetails is a DB Serialization of Proto TransactionDetails
type TransactionDetails struct {
	Address         string    `bson:"address"`
	DoingBusinessAs string    `bson:"doing_business_as"`
	DateProcessed   time.Time `bson:"date_processed"`
}

// Transaction is a DB Serialization of Proto Transaction
type Transaction struct {
	ID                 primitive.ObjectID `bson:"_id"`
	UserId             primitive.ObjectID `bson:"user_id"`
	AccountId          primitive.ObjectID `bson:"account_id"`
	Name               string             `bson:"name"`
	Amount             float64            `bson:"amount"`
	Date               time.Time          `bson:"date"`
	RewardsEarned      int32              `bson:"rewards_earned"`
	TransactionDetails TransactionDetails `bson:"transaction_details"`
}

// User is a DB Serialization of Proto User
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
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
