package token

import (
	"context"
	"errors"
)

type spyTokenResolver struct {
	tokenResponse string
	callCount     int
}

func newSpyTokenResolver(tokenResponse string) *spyTokenResolver {
	return &spyTokenResolver{
		tokenResponse: tokenResponse,
	}
}

func (m *spyTokenResolver) Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error) {
	m.callCount++
	return m.tokenResponse, nil
}

var _ TokenResolver = &spyTokenResolver{}

type mockErrorTokenResolver struct{}

func newMockErrorTokenResolver() *mockErrorTokenResolver {
	return &mockErrorTokenResolver{}
}

// Resolve implements TokenResolver.
func (m *mockErrorTokenResolver) Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error) {
	return "", errors.New("test error")
}

var _ TokenResolver = &mockErrorTokenResolver{}
