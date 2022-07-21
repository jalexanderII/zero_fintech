package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthUser is a DB Serialization of Proto User
type AuthUser struct {
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
