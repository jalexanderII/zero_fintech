package client

import (
	"context"
	"log"
	"time"

	"github.com/jalexanderII/zero_fintech/bff/middleware"
	"github.com/jalexanderII/zero_fintech/gen/Go/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient is a client to call authentication RPC
type AuthClient struct {
	authClient  auth.AuthClient
	Interceptor *middleware.AuthInterceptor
	Username    string
	Email       string
	Password    string
	PhoneNumber string
}

// NewAuthClient returns a new auth client
func NewAuthClient(conn *grpc.ClientConn, username, email, password, phonenumber string) *AuthClient {
	a := auth.NewAuthClient(conn)
	return &AuthClient{a, middleware.NewAuthInterceptor(middleware.AccessibleRoles()), username, email, password, phonenumber}
}

// Login user and returns the access token
func (a *AuthClient) Login() (*auth.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &auth.LoginRequest{
		Username: a.Username,
		Password: a.Password,
	}

	res, err := a.authClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	a.Interceptor.SetToken(res.GetToken())

	return res, nil
}

// SignUp creates a new user and returns a new access-token
func (a *AuthClient) SignUp() (*auth.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &auth.SignupRequest{
		Username:    a.Username,
		Email:       a.Email,
		Password:    a.Password,
		PhoneNumber: a.PhoneNumber,
	}

	res, err := a.authClient.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}
	a.Interceptor.SetToken(res.GetToken())

	return res, nil
}

func SetUpAuthClient() (*AuthClient, []grpc.DialOption) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	authConn, err := grpc.Dial("localhost:9091", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return NewAuthClient(authConn, "", "", "", ""), opts
}
