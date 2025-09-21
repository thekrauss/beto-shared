package middleware

import (
	"context"

	"github.com/thekrauss/beto-shared/pkg/errors"
	"google.golang.org/grpc"
)

// attrape les erreurs et renvoie un gRPC status unifié
func GRPCErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			return nil, errors.ToGRPCError(err)
		}
		return resp, nil
	}
}
