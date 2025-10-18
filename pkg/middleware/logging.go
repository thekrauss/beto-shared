package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/thekrauss/beto-shared/pkg/logger"
	"google.golang.org/grpc"
)

// responseRecorder permet de capturer le status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware log chaque requête HTTP avec Zap
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rr, r)

		duration := time.Since(start)
		reqID := GetRequestID(r.Context())

		logger.FromContext(r.Context()).Infow("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rr.statusCode,
			"duration", duration.String(),
			"request_id", reqID,
			"user_agent", r.UserAgent(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

// GRPCLoggingInterceptor log chaque appel gRPC
func GRPCLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.FromContext(ctx).Errorw("gRPC request failed",
				"method", info.FullMethod,
				"duration", duration.String(),
				"error", err,
			)
			return nil, err
		}

		logger.FromContext(ctx).Infow("gRPC request",
			"method", info.FullMethod,
			"duration", duration.String(),
		)
		return resp, nil
	}
}
