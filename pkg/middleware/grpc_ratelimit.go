package middleware

import (
	"context"
	"time"

	"github.com/thekrauss/beto-shared/pkg/errors"
	"github.com/thekrauss/beto-shared/pkg/redis"
	"google.golang.org/grpc"
)

// applique un rate limit distribué via Redis
func GRPCRateLimitInterceptor(keyFunc func(ctx context.Context, req interface{}) string, limit int, window time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		key := keyFunc(ctx, req)

		allowed, rerr := redis.AllowRequest(ctx, key, limit, window)
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
nexxt step chaque microservice doit lire ses limites depuis la config et applique automatiquement les middlewares
*/
