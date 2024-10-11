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
func (m *MultiTokenResolver) Resolve(ctx context.Context) (*Token, error) {
	resultErrors := []error{
		ErrNoValidTokenResolved,
	}
	for _, resolver := range m.resolvers {
		token, err := resolver.Resolve(ctx)
		if err != nil {
			resultErrors = append(resultErrors, err)
			continue
		}
		return token, nil
	}
	return nil, errors.Join(resultErrors...)
}

var _ TokenResolver = &MultiTokenResolver{}
