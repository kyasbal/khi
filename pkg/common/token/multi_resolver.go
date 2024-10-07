package token

import (
	"context"
	"errors"
)

var ErrNoValidTokenResolved = errors.New("no valid token returned")

type MultiTokenResolver struct {
	resolvers []TokenResolver
}

func NewMultiTokenResolver(resolvers ...TokenResolver) *MultiTokenResolver {
	return &MultiTokenResolver{
		resolvers: resolvers,
	}
}

// Resolve implements TokenResolver.
func (m *MultiTokenResolver) Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error) {
	for _, resolver := range m.resolvers {
		token, err := resolver.Resolve(ctx, expiredTokens)
		if err != nil {
			continue
		}
		return token, nil
	}
	return "", ErrNoValidTokenResolved
}

var _ TokenResolver = &MultiTokenResolver{}
