package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// records gRPC metrics
func GRPCMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		duration := time.Since(start).Seconds()
		st, _ := status.FromError(err)

		GRPCRequestsTotal.WithLabelValues(info.FullMethod, st.Code().String()).Inc()
		GRPCRequestDuration.WithLabelValues(info.FullMethod).Observe(duration)

		return resp, err
	}
}
