package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/jalexanderII/zero_fintech/services/Core/database"
	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	EmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

func (s CoreServer) Login(ctx context.Context, in *core.LoginRequest) (*core.AuthResponse, error) {
	username, password := in.GetUsername(), in.GetPassword()
	var user database.User
	err := s.UserDB.FindOne(ctx, bson.M{"$or": []bson.M{{"username": username}, {"email": username}}}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %v", err)
	}
	if user.ID.IsZero() || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("wrong login credentials provided")
	}

	token, err := s.jwtm.Generate(&user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	return &core.AuthResponse{Token: token}, nil
}

func (s CoreServer) SignUp(ctx context.Context, in *core.SignupRequest) (*core.AuthResponse, error) {
	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	match, _ := regexp.MatchString(EmailRegex, email)
	if !match {
		return nil, errors.New("email validation failed")
	}

	res, err := s.UsernameUsed(ctx, username)
	if err != nil {
		log.Printf("Error returned from UsernameUsed: %v\n", err)
		return nil, err
	}
	if res {
		return nil, errors.New("username already taken")
	}
	res, err = s.EmailUsed(ctx, email)
	if err != nil {
		log.Printf("Error returned from EmailUsed: %v\n", err)
		return nil, err
	}
	if res {
		return nil, errors.New("email already used")
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := database.User{ID: primitive.NewObjectID(), Email: email, Username: username, Password: string(pw)}

	_, err = s.UserDB.InsertOne(ctx, newUser)
	if err != nil {
		log.Printf("Error inserting new user: %v\n", err)
		return nil, err
	}

	token, err := s.jwtm.Generate(&newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	return &core.AuthResponse{Token: token}, nil
}

func (s CoreServer) EmailUsed(ctx context.Context, email string) (bool, error) {
	var user database.User
	filter := bson.D{{"email", email}}
	err := s.UserDB.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, fmt.Errorf("error fetching email: %v", err)
	}
	s.l.Info("email already exists", "user", user.ID)
	return true, nil
}

func (s CoreServer) UsernameUsed(ctx context.Context, username string) (bool, error) {
	var user database.User
	filter := bson.D{{"username", username}}
	err := s.UserDB.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, fmt.Errorf("error fetching username: %v", err)
	}
	s.l.Info("username already exists", "user", user.ID)
	return true, nil
}

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
	return &core.DeleteUserResponse{Status: core.DELETE_STATUS_DELETE_STATUS_SUCCESS, User: UserDBToPB(&user)}, nil
}

func UserDBToPB(user *database.User) *core.User {
	return &core.User{
		Id:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}
