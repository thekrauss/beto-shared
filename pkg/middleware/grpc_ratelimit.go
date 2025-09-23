package middleware

import (
	"context"

	"github.com/thekrauss/beto-shared/pkg/config"
	"github.com/thekrauss/beto-shared/pkg/errors"
	"github.com/thekrauss/beto-shared/pkg/redis"
	"google.golang.org/grpc"
)

// applies a rate limit distributed via Redis
func GRPCRateLimitInterceptor(keyFunc func(ctx context.Context, req interface{}) string, cfg config.RateLimitConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		key := keyFunc(ctx, req)

		allowed, rerr := redis.AllowRequest(ctx, key, cfg.Limit, cfg.Window)
		if rerr != nil {
			return nil, errors.ToGRPCError(errors.New(errors.CodeInternal, "rate limit check failed"))
		}
		if !allowed {
			return nil, errors.ToGRPCError(errors.New(errors.CodeTimeout, "rate limit exceeded"))
		}

		return handler(ctx, req)
	}
}

/*
 */
