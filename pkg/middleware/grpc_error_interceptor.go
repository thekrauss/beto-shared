package middleware

import (
	"context"

	"github.com/thekrauss/beto-shared/pkg/errors"
	"google.golang.org/grpc"
)

// catch errors and return a unified gRPC status
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
