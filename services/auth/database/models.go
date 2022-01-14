package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthUser is a DB Serialization of Proto User
type AuthUser struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}
