package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/gen/Go/auth"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/jalexanderII/zero_fintech/services/auth/database"

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

// AuthServer is the server for the AuthService, it will connect to its own mongodb database and will be reachable via
// grpc from microservices and via grpc proxy for clients
type AuthServer struct {
	auth.UnimplementedAuthServer
	// Database collections
	UserDB mongo.Collection
	// authentication manager
	jwtm *middleware.JWTManager
	// custom logger
	l hclog.Logger
}

func NewAuthServer(udb mongo.Collection, jwtm *middleware.JWTManager, l hclog.Logger) *AuthServer {
	return &AuthServer{UserDB: udb, jwtm: jwtm, l: l}
}

func (s AuthServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.AuthResponse, error) {
	username, password := in.GetUsername(), in.GetPassword()
	var user database.AuthUser
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

	return &auth.AuthResponse{Token: token}, nil
}

func (s AuthServer) SignUp(ctx context.Context, in *auth.SignupRequest) (*auth.AuthResponse, error) {
	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	// Email regex validation
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

	// hashed passwords are saved in the DB
	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := database.AuthUser{ID: primitive.NewObjectID(), Email: email, Username: username, Password: string(pw)}

	_, err = s.UserDB.InsertOne(ctx, newUser)
	if err != nil {
		log.Printf("Error inserting new user: %v\n", err)
		return nil, err
	}

	token, err := s.jwtm.Generate(&newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	return &auth.AuthResponse{Token: token}, nil
}

// EmailUsed checks if the email is already present in the DB
func (s AuthServer) EmailUsed(ctx context.Context, email string) (bool, error) {
	var user database.AuthUser
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

// UsernameUsed checks if the username is already present in the DB
func (s AuthServer) UsernameUsed(ctx context.Context, username string) (bool, error) {
	var user database.AuthUser
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
