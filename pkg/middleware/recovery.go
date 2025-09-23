package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

// captures panics and returns an error 500
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				status, body := errors.ToHTTPError(
					errors.New(errors.CodeInternal, "internal server error"),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				w.Write(body)
				debug.PrintStack()
			}
		}()
		next.ServeHTTP(w, r)
	})
}
