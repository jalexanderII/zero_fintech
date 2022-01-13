package server

import (
	"testing"

	"github.com/jalexanderII/zero_fintech/services/core/database"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Test_authServer_SignUp(t *testing.T) {
	server, ctx := GenServer()

	u := &core.SignupRequest{
		Username: "joel_admin",
		Email:    "fudoshin2596@gmail.com",
		Password: "joel_admin",
	}

	_, err := server.SignUp(ctx, u)
	if err != nil {
		t.Errorf("1: Error creating new user: %v", err)
	}

	_, err = server.SignUp(ctx, &core.SignupRequest{Username: "example", Email: "bad-email", Password: "example"})
	if err.Error() != "email validation failed" {
		t.Error("2: No or wrong error returned for email validation")
	}

	_, err = server.SignUp(ctx, &core.SignupRequest{Username: u.Username, Email: "e@gmail.com", Password: "e"})
	if err.Error() != "username already taken" {
		t.Error("3: No or wrong error returned for username already taken")
	}

	_, err = server.SignUp(ctx, &core.SignupRequest{Username: "e", Email: u.Email, Password: "e"})
	if err.Error() != "email already used" {
		t.Error("4: No or wrong error returned for email already taken")
	}
}

func TestAuthServer_GetUser(t *testing.T) {
	server, ctx := GenServer()

	user, err := server.GetUser(ctx, &core.GetUserRequest{Id: "61df93c0ac601d1be8e64613"})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if user.Username != "joel_admin" {
		t.Errorf("2: Failed to fetch correct user: %+v", user)
	}
}

func TestAuthServer_ListUsers(t *testing.T) {
	server, ctx := GenServer()

	users, err := server.ListUsers(ctx, &core.ListUserRequest{})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(users.Users) < 1 {
		t.Errorf("2: Failed to fetch realtors: %+v", users.Users[0])
	}
}

func TestAuthServer_UpdateUser(t *testing.T) {
	server, ctx := GenServer()

	u := &core.User{
		Username: "exampleUpdated",
		Email:    "exampleUpdated@gmail.com",
		Password: "exampleUpdated",
	}

	user, err := server.UpdateUser(ctx, &core.UpdateUserRequest{Id: "61df94304f6d9090c7be7b55", User: u})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if user.Email != u.Email {
		t.Errorf("2: Failed to fetch correct user: %+v", user)
	}
}

func TestAuthServer_DeleteUser(t *testing.T) {
	server, ctx := GenServer()

	u, err := server.ListUsers(ctx, &core.ListUserRequest{})
	originalLen := len(u.GetUsers())

	pw, _ := bcrypt.GenerateFromPassword([]byte("to_delete"), bcrypt.DefaultCost)
	newUser := database.User{ID: primitive.NewObjectID(), Email: "to_delete@gmail.com", Username: "to_delete", Password: string(pw)}
	_, err = server.UserDB.InsertOne(ctx, &newUser)
	if err != nil {
		t.Errorf("1: Error creating new user:: %v", err)
	}

	users, err := server.ListUsers(ctx, &core.ListUserRequest{})
	if err != nil {
		t.Errorf("2: An error was returned: %v", err)
	}
	newLen := len(users.GetUsers())
	if newLen != originalLen+1 {
		t.Errorf("3: An error adding a temp user, number of users in DB: %v", newLen)
	}

	deleted, err := server.DeleteUser(ctx, &core.DeleteUserRequest{Id: newUser.ID.Hex()})
	if err != nil {
		t.Errorf("4: An error was returned: %v", err)
	}
	if deleted.Status != core.DELETE_STATUS_DELETE_STATUS_SUCCESS {
		t.Errorf("5: Failed to delete user: %+v\n, %+v", deleted.Status, deleted.GetUser())
	}
}

func Test_authServer_Login(t *testing.T) {
	server, ctx := GenServer()

	u := &core.SignupRequest{
		Username: "guest",
		Email:    "guest@gmail.com",
		Password: "guest",
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	_, err := server.UserDB.InsertOne(ctx, database.User{ID: primitive.NewObjectID(), Username: u.Username, Email: u.Email, Password: string(pw)})
	if err != nil {
		t.Errorf("1: Error inserting new user into db: %v", err)
	}

	_, err = server.Login(ctx, &core.LoginRequest{Username: u.Username, Password: u.Password})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}

	_, err = server.Login(ctx, &core.LoginRequest{Username: u.Username, Password: "wrong"})
	if err == nil {
		t.Error("2: Error was nil")
	}
}
