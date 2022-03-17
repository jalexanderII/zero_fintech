package models

import (
	"github.com/plaid/plaid-go/plaid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token for use of plaid public token retrieval
type Token struct {
	ID            primitive.ObjectID `bson:"_id"`
	User          *User              `bson:"user"`
	Value         string             `bson:"value"`
	ItemId        string             `bson:"item_id"`
	Institution   string             `bson:"institution"`
	InstitutionID string             `bson:"institution_id"`
}

type CreateLinkTokenResponse struct {
	UserId string
	Token  string
}

// User is a DB Serialization of Proto User
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
}

type LiabilitiesResponse struct {
	Liabilities []plaid.CreditCardLiability `json:"liabilities"`
}

type TransactionsResponse struct {
	Accounts     []plaid.AccountBase `json:"accounts,omitempty"`
	Transactions []plaid.Transaction `json:"transactions,omitempty"`
}

type PlaidMetaData struct {
	Institution struct {
		Name          string `json:"name"`
		InstitutionId string `json:"institution_id"`
	} `json:"institution"`
	Accounts []struct {
		Id                 string `json:"id"`
		Name               string `json:"name"`
		Mask               string `json:"mask"`
		Type               string `json:"type"`
		Subtype            string `json:"subtype"`
		VerificationStatus string `json:"verification_status,omitempty"`
	} `json:"accounts"`
	LinkSessionId string `json:"link_session_id"`
}
