package token

import (
	"context"
)

// TokenResolver resolve the new token
type TokenResolver interface {
	// Resolve returns the token newly resolved.
	// expiredTokens are just for information to each resolvers. Each resolvers can return the token even if it was expired in the map when the resolver believe the token is active yet.
	Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error)
}
