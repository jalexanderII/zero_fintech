package database

import (
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentTask struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserId        primitive.ObjectID `bson:"user_id"`
	TransactionId primitive.ObjectID `bson:"transaction_id"`
	AccountId     primitive.ObjectID `bson:"account_id"`
	MetaData      MetaData           `bson:"meta_data"`
}

type MetaData struct {
	PreferredPlanType    core.PlanType         `bson:"preferred_plan_type"`
	PreferredTimeline    float64               `bson:"preferred_timeline"`
	PreferredPaymentFreq core.PaymentFrequency `bson:"preferred_payment_freq"`
}

type AnnualPercentageRates struct {
	LowEnd  float64 `bson:"low_end"`
	HighEnd float64 `bson:"high_end"`
}

type PenaltyAPR struct {
	PenaltyAPR    float64                       `bson:"penalty_apr"`
	PenaltyReason core.PenaltyAPR_PenaltyReason `bson:"penalty_reason"`
}

type PromotionalRate struct {
	TemporaryAPR   float64             `bson:"temporary_apr"`
	ExpirationDate primitive.Timestamp `bson:"expiration_date"`
}

type Account struct {
	ID                    primitive.ObjectID    `bson:"_id"`
	UserId                primitive.ObjectID    `bson:"user_id"`
	Name                  string                `bson:"name"`
	CreatedAt             primitive.Timestamp   `bson:"created_at"`
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

type TransactionDetails struct {
	Address         string              `bson:"address"`
	DoingBusinessAs string              `bson:"doing_business_as"`
	DateProcessed   primitive.Timestamp `bson:"date_processed"`
}
type Transaction struct {
	ID                 primitive.ObjectID  `bson:"_id"`
	UserId             primitive.ObjectID  `bson:"user_id"`
	AccountId          primitive.ObjectID  `bson:"account_id"`
	Name               string              `bson:"name"`
	Amount             float64             `bson:"amount"`
	Date               primitive.Timestamp `bson:"date"`
	RewardsEarned      int32               `bson:"rewards_earned"`
	TransactionDetails TransactionDetails  `bson:"transaction_details"`
}
