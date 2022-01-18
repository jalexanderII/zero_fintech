package interceptor

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	jwtManager      *middleware.JWTManager
	accessibleRoles map[string]bool
	l               hclog.Logger
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(jwtManager *middleware.JWTManager, accessibleRoles map[string]bool, l hclog.Logger) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles, l}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		interceptor.l.Debug("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	_, ok := interceptor.accessibleRoles[method]
	if ok {
		// everyone can access
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	interceptor.l.Debug("Meta data from client: ", md)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided", md)
	}

	accessToken := values[0]
	_, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return nil
}

func AccessibleRoles() map[string]bool {
	const authServicePath = "/auth.Auth/"
	return map[string]bool{
		// Auth paths not Protected since they are needed to generate the tokens
		authServicePath + "Login":  true,
		authServicePath + "SignUp": true,
	}
}
