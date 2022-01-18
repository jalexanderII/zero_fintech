package middleware

import (
	"context"
	"time"

	"github.com/jalexanderII/zero_fintech/services/auth/gen/auth"
	"google.golang.org/grpc"
)

// AuthClient is a client to call authentication RPC
type AuthClient struct {
	authClient  auth.AuthClient
	Interceptor *AuthInterceptor
	Username    string
	Email       string
	Password    string
}

// NewAuthClient returns a new auth client
func NewAuthClient(conn *grpc.ClientConn, username, email, password string) *AuthClient {
	a := auth.NewAuthClient(conn)
	return &AuthClient{a, NewAuthInterceptor(AccessibleRoles()), username, email, password}
}

// Login user and returns the access token
func (a *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &auth.LoginRequest{
		Username: a.Username,
		Password: a.Password,
	}

	res, err := a.authClient.Login(ctx, req)
	if err != nil {
		return "", err
	}
	a.Interceptor.SetToken(res.GetToken())

	return res.GetToken(), nil
}

// SignUp creates a new user and returns a new access-token
func (a *AuthClient) SignUp() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &auth.SignupRequest{
		Username: a.Username,
		Email:    a.Email,
		Password: a.Password,
	}

	res, err := a.authClient.SignUp(ctx, req)
	if err != nil {
		return "", err
	}
	a.Interceptor.SetToken(res.GetToken())

	return res.GetToken(), nil
}
