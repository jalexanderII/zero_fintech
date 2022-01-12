package database

import (
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
}
