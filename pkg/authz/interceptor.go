package authz

import (
	"context"
	"strings"

	"github.com/thekrauss/beto-shared/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// protects gRPC calls
func KeystoneAuthInterceptor(validator *KeystoneValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.ToGRPCError(errors.New(errors.CodeUnauthorized, "missing metadata"))
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 || !strings.HasPrefix(authHeaders[0], "Bearer ") {
			return nil, errors.ToGRPCError(errors.New(errors.CodeUnauthorized, "missing bearer token"))
		}

		token := strings.TrimPrefix(authHeaders[0], "Bearer ")
		claims, verr := validator.ValidateToken(ctx, token)
		if verr != nil {
			return nil, errors.ToGRPCError(verr)
		}

		ctx = WithClaims(ctx, claims)
		return handler(ctx, req)
	}
}
