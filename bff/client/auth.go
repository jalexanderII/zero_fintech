package client

import (
	"context"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"google.golang.org/grpc"
)

func (c *CoreClient) Login(ctx context.Context, in *core.LoginRequest, opts ...grpc.CallOption) (*core.AuthResponse, error) {
	return nil, nil
}
func (c *CoreClient) SignUp(ctx context.Context, in *core.SignupRequest, opts ...grpc.CallOption) (*core.AuthResponse, error) {
	return nil, nil
}
