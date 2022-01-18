package server

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/auth/database"
	"github.com/jalexanderII/zero_fintech/services/auth/gen/auth"
	"github.com/jalexanderII/zero_fintech/services/core/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var L = hclog.Default()

func GenServer() (*AuthServer, context.Context) {
	jwtManager := middleware.NewJWTManager(config.GetEnv("JWTSecret"), 15*time.Minute)
	DB, err := database.InitiateMongoClient()
	if err != nil {
		log.Fatal("MongoDB error: ", err)
	}
	userCollection := *DB.Collection(config.GetEnv("USER_COLLECTION"))

	server := NewAuthServer(userCollection, jwtManager, L)
	return server, context.TODO()
}

func Test_authServer_SignUp(t *testing.T) {
	server, ctx := GenServer()

	u := &auth.SignupRequest{
		Username: "joel_admin",
		Email:    "fudoshin2596@gmail.com",
		Password: "joel_admin",
	}

	_, err := server.SignUp(ctx, u)
	if err != nil {
		t.Errorf("1: Error creating new user: %v", err)
	}

	_, err = server.SignUp(ctx, &auth.SignupRequest{Username: "example", Email: "bad-email", Password: "example"})
	if err.Error() != "email validation failed" {
		t.Error("2: No or wrong error returned for email validation")
	}

	_, err = server.SignUp(ctx, &auth.SignupRequest{Username: u.Username, Email: "e@gmail.com", Password: "e"})
	if err.Error() != "username already taken" {
		t.Error("3: No or wrong error returned for username already taken")
	}

	_, err = server.SignUp(ctx, &auth.SignupRequest{Username: "e", Email: u.Email, Password: "e"})
	if err.Error() != "email already used" {
		t.Error("4: No or wrong error returned for email already taken")
	}
}

func Test_authServer_Login(t *testing.T) {
	server, ctx := GenServer()

	u := &auth.SignupRequest{
		Username: "guest",
		Email:    "guest@gmail.com",
		Password: "guest",
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	_, err := server.UserDB.InsertOne(ctx, database.AuthUser{ID: primitive.NewObjectID(), Username: u.Username, Email: u.Email, Password: string(pw)})
	if err != nil {
		t.Errorf("1: Error inserting new user into db: %v", err)
	}

	_, err = server.Login(ctx, &auth.LoginRequest{Username: u.Username, Password: u.Password})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}

	_, err = server.Login(ctx, &auth.LoginRequest{Username: u.Username, Password: "wrong"})
	if err == nil {
		t.Error("2: Error was nil")
	}
}
