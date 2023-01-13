package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor is a client interceptor for authentication
type AuthInterceptor struct {
	AuthMethods map[string]bool
	AccessToken string
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(authMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{AuthMethods: authMethods}
}

// Unary returns a client interceptor to authenticate unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> Auth unary interceptor: %s", method)

		if !interceptor.AuthMethods[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	log.Println("token AppendToOutgoingContext")
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.AccessToken)
}

func (interceptor *AuthInterceptor) SetToken(accessToken string) {
	interceptor.AccessToken = accessToken
	log.Println("token set", interceptor.AccessToken)
}

func AccessibleRoles() map[string]bool {
	const authServicePath = "/auth.Auth/"
	const coreServicePath = "/core.Core/"
	return map[string]bool{
		// Auth paths are not Protected since they are needed to generate the tokens
		authServicePath + "Login":   true,
		authServicePath + "SignUp":  true,
		coreServicePath + "GetUser": true,
	}
}
