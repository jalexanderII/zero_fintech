package models

import (
	"github.com/plaid/plaid-go/plaid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token for use of plaid public token retrieval
type Token struct {
	ID          primitive.ObjectID `bson:"_id"`
	User        User               `bson:"user"`
	Value       string             `bson:"value"`
	ItemId      string             `bson:"item_id"`
	Institution string             `bson:"institution"`
}

// User is a DB Serialization of Proto User
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type LiabilitiesResponse struct {
	Liabilities []plaid.CreditCardLiability `json:"liabilities"`
}

type TransactionsResponse struct {
	Accounts     []plaid.AccountBase `json:"accounts,omitempty"`
	Transactions []plaid.Transaction `json:"transactions,omitempty"`
}
