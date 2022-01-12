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

type Account struct {
	ID primitive.ObjectID `bson:"_id"`
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
