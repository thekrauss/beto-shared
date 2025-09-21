package authz

import (
	"net/http"
	"strings"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

// KeystoneAuthMiddleware protège une route HTTP
func KeystoneAuthMiddleware(validator *KeystoneValidator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			status, body := errors.ToHTTPError(errors.New(errors.CodeUnauthorized, "missing bearer token"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			w.Write(body)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validator.ValidateToken(r.Context(), token)
		if err != nil {
			status, body := errors.ToHTTPError(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			w.Write(body)
			return
		}

		ctx := WithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
