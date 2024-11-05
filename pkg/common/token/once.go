package token

import (
	"context"
)

// OnceTokenResolver resolves the token from somewhere only onetime. This is mainly used for resolving token from CLI arguments or environment variables.
type OnceTokenResolver struct {
	// resolver can be called twice if the resolver can't return token at the previous call.
	resolver func() string
	// tokenResolved is set to true once this resolver returns a token.
	tokenResolved bool
}

func NewOnceTokenResolver(resolver func() string) *OnceTokenResolver {
	return &OnceTokenResolver{
		resolver:      resolver,
		tokenResolved: false,
	}
}

// Resolve implements TokenResolver.
func (e *OnceTokenResolver) Resolve(ctx context.Context) (*Token, error) {
	if !e.tokenResolved {
		token := e.resolver()
		if token != "" {
			e.tokenResolved = true
			return New(token), nil
		}
	}
	if e.tokenResolved {
		return nil, ErrNoNewTokenResolved
	}
	return nil, ErrNoValidTokenResolved
}

var _ TokenResolver = (*OnceTokenResolver)(nil)
