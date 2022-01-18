package server

import (
	"context"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/services/core/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s CoreServer) GetUser(ctx context.Context, in *core.GetUserRequest) (*core.User, error) {
	var user database.User
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	err = s.UserDB.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return UserDBToPB(&user), nil
}

func (s CoreServer) ListUsers(ctx context.Context, in *core.ListUserRequest) (*core.ListUserResponse, error) {
	var results []database.User
	cursor, err := s.UserDB.Find(ctx, bson.D{})
	if err = cursor.All(ctx, &results); err != nil {
		s.l.Error("[DB] Error getting all users", "error", err)
		return nil, err
	}
	res := make([]*core.User, len(results))
	for idx, user := range results {
		res[idx] = UserDBToPB(&user)
	}
	return &core.ListUserResponse{Users: res}, nil
}

func (s CoreServer) UpdateUser(ctx context.Context, in *core.UpdateUserRequest) (*core.User, error) {
	username, email := in.GetUser().GetUsername(), in.GetUser().GetEmail()
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"username", username}, {"email", email}}}}
	_, err = s.UserDB.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	var user database.User
	err = s.UserDB.FindOne(ctx, filter).Decode(&user)
	return UserDBToPB(&user), nil
}

func (s CoreServer) DeleteUser(ctx context.Context, in *core.DeleteUserRequest) (*core.DeleteUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", id}}
	_, err = s.UserDB.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	var user database.User
	err = s.UserDB.FindOne(ctx, filter).Decode(&user)
	return &core.DeleteUserResponse{Status: common.DELETE_STATUS_DELETE_STATUS_SUCCESS, User: UserDBToPB(&user)}, nil
}

// UserDBToPB converts a User DB object to its proto object
func UserDBToPB(user *database.User) *core.User {
	return &core.User{
		Id:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}
