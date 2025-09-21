package authz

import "context"

type contextKey string

const claimsKey contextKey = "authz_claims"

// injecte les claims dans le contexte
func WithClaims(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

// récupère les claims depuis le contexte
func GetClaims(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*Claims)
	return c, ok
}

// vérifie si l’utilisateur a un rôle donné
func RequireRole(ctx context.Context, role string) bool {
	c, ok := GetClaims(ctx)
	if !ok {
		return false
	}
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}
