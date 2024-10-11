package token

import (
	"context"
)

// TokenResolver resolve the new token
type TokenResolver interface {
	// Resolve returns the token newly resolved.
	Resolve(ctx context.Context) (*Token, error)
}
