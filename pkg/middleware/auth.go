package middleware

import (
	"net/http"
	"strings"

	"github.com/thekrauss/beto-shared/pkg/authz"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

// retourne un middleware qui valide les tokens via Keystone
func AuthMiddleware(validator *authz.KeystoneValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			//  claims dans le context
			ctx := authz.WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
