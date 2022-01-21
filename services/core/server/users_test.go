package server

import (
	"testing"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/core/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

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

	u, _ := server.ListUsers(ctx, &core.ListUserRequest{})
	originalLen := len(u.GetUsers())

	pw, _ := bcrypt.GenerateFromPassword([]byte("to_delete"), bcrypt.DefaultCost)
	newUser := database.User{ID: primitive.NewObjectID(), Email: "to_delete@gmail.com", Username: "to_delete", Password: string(pw)}
	_, err := server.UserDB.InsertOne(ctx, &newUser)
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
	if deleted.Status != common.DELETE_STATUS_DELETE_STATUS_SUCCESS {
		t.Errorf("5: Failed to delete user: %+v\n, %+v", deleted.Status, deleted.GetUser())
	}
}
