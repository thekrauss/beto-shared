package middleware

import (
	"net/http"

	"github.com/thekrauss/beto-shared/pkg/config"
	"github.com/thekrauss/beto-shared/pkg/errors"
	"github.com/thekrauss/beto-shared/pkg/redis"
)

// limits HTTP requests per key
func RateLimitMiddleware(keyFunc func(r *http.Request) string, cfg config.RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			allowed, err := redis.AllowRequest(r.Context(), key, cfg.Limit, cfg.Window)
			if err != nil {
				status, body := errors.ToHTTPError(errors.New(errors.CodeInternal, "rate limit check failed"))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				w.Write(body)
				return
			}

			if !allowed {
				status, body := errors.ToHTTPError(errors.New(errors.CodeTimeout, "rate limit exceeded"))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				w.Write(body)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
