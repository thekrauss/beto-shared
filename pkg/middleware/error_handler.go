package middleware

import (
	"net/http"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

// attrape les erreurs et renvoie un JSON unifié
func ErrorHandlerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				status, body := errors.ToHTTPError(errors.New(errors.CodeInternal, "internal server error"))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				w.Write(body)
			}
		}()
		next(w, r)
	}
}
