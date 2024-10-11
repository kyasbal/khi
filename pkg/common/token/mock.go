package token

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type SpyTokenResolver struct {
	tokenResponse      []*Token
	callCount          int
	delayInMillisecond int
}

func NewSpyTokenResolverWithDelay(delayInMillisecond int, tokenResponse ...*Token) *SpyTokenResolver {
	return &SpyTokenResolver{
		tokenResponse:      tokenResponse,
		delayInMillisecond: delayInMillisecond,
	}
}

func NewSpyTokenResolver(tokenResponse ...*Token) *SpyTokenResolver {
	return &SpyTokenResolver{
		tokenResponse: tokenResponse,
	}
}

func (m *SpyTokenResolver) Resolve(ctx context.Context) (*Token, error) {
	if m.callCount < len(m.tokenResponse) {
		m.callCount++
		time.Sleep(time.Duration(m.delayInMillisecond * int(time.Millisecond)))
		return m.tokenResponse[m.callCount-1], nil
	} else {
		return nil, fmt.Errorf("no expected token response supplied")
	}
}

var _ TokenResolver = &SpyTokenResolver{}

type MockErrorTokenResolver struct{}

func NewMockErrorTokenResolver() *MockErrorTokenResolver {
	return &MockErrorTokenResolver{}
}

// Resolve implements TokenResolver.
func (m *MockErrorTokenResolver) Resolve(ctx context.Context) (*Token, error) {
	return nil, errors.New("test error")
}

var _ TokenResolver = &MockErrorTokenResolver{}
