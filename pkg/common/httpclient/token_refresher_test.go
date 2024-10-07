package httpclient

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

type tokenResolverSpy struct {
	CallCount      int
	ParentResolver token.TokenResolver
}

// Resolve implements token.TokenResolver.
func (t *tokenResolverSpy) Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error) {
	t.CallCount++
	return t.ParentResolver.Resolve(ctx, expiredTokens)
}

func newTokenResolverSpy(resolver token.TokenResolver) *tokenResolverSpy {
	return &tokenResolverSpy{
		CallCount:      0,
		ParentResolver: resolver,
	}
}

var _ token.TokenResolver = (*tokenResolverSpy)(nil)

type testTokenResolver struct {
	Tokens []string
}

// Resolve implements token.TokenResolver.
func (t *testTokenResolver) Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error) {
	if len(t.Tokens) == 0 {
		return "", token.ErrNoNewTokenResolved
	}
	token := t.Tokens[0]
	if len(token) > 0 {
		t.Tokens = t.Tokens[1:]
	}
	return token, nil
}

var _ token.TokenResolver = (*testTokenResolver)(nil)

func newTestTokenResolver(tokens []string) *testTokenResolver {
	return &testTokenResolver{
		Tokens: tokens,
	}
}

func TestMultiTokenStoreRefresher_Refresh(t *testing.T) {
	spy1 := newTokenResolverSpy(newTestTokenResolver([]string{"token1-0", "token1-1", "token1-2"}))
	spy2 := newTokenResolverSpy(newTestTokenResolver([]string{"token2-0", "token2-1", "token2-2"}))
	store1 := token.NewBasicTokenStore("store1", spy1)
	store2 := token.NewBasicTokenStore("store1", spy2)
	refresher := NewMultiTokenStoreRefresher(store1, store2)
	store1.GetToken(context.Background())
	store2.GetToken(context.Background())
	if spy1.CallCount != 1 {
		t.Errorf("token1 resolver must only called once but called %d", spy1.CallCount)
	}
	if spy2.CallCount != 1 {
		t.Errorf("token2 resolver must only called once but called %d", spy2.CallCount)
	}
	refresher.Refresh(context.Background())
	store1.GetToken(context.Background())
	store2.GetToken(context.Background())
	if spy1.CallCount != 2 {
		t.Errorf("token1 resolver must only called once but called %d", spy1.CallCount)
	}
	if spy2.CallCount != 1 {
		t.Errorf("token2 resolver must not called but called %d", spy2.CallCount)
	}
	refresher.Refresh(context.Background())
	store1.GetToken(context.Background())
	store2.GetToken(context.Background())
	if spy1.CallCount != 3 {
		t.Errorf("token1 resolver must only called once but called %d", spy1.CallCount)
	}
	if spy2.CallCount != 1 {
		t.Errorf("token2 resolver must only called once but called %d", spy2.CallCount)
	}
	refresher.Refresh(context.Background())
	store1.GetToken(context.Background())
	store2.GetToken(context.Background())
	if spy1.CallCount != 4 {
		t.Errorf("token1 resolver must only called twice but called %d", spy1.CallCount)
	}
	if spy2.CallCount != 2 {
		t.Errorf("token2 resolver must only called once but called %d", spy2.CallCount)
	}
}
