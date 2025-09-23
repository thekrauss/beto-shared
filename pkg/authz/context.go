package authz

import "context"

type contextKey string

const claimsKey contextKey = "authz_claims"

// injects claims into context
func WithClaims(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

func GetClaims(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*Claims)
	return c, ok
}

// collects claims from context
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
