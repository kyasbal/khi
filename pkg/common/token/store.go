package token

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var ErrNoNewTokenResolved = errors.New("no new token resolved")

type TokenStore interface {
	task.CachableDependency
	GetType() string
	// GetToken returns the current token. It can come from cache or newly resolved from TokenResolver.
	GetToken(ctx context.Context) (string, error)
	// RefreshToken try to get new token
	RefreshToken(ctx context.Context) (string, error)
}

// BasicTokenStore provides feature to refresh token and return cached token.
// BasicTokenStore memory the expired tokens and it calls resolvers in order to get new token after MarkTokenExpired called.
type BasicTokenStore struct {
	tokenType     string
	resolver      TokenResolver
	tokenLock     sync.Mutex
	expiredTokens map[string]interface{}
	lastToken     string
}

// Digest implements TokenStore.
func (b *BasicTokenStore) Digest() string {
	return b.tokenType
}

func NewBasicTokenStore(tokenType string, resolver TokenResolver) *BasicTokenStore {
	return &BasicTokenStore{
		tokenType:     tokenType,
		resolver:      resolver,
		expiredTokens: map[string]interface{}{},
		lastToken:     "",
	}
}

func (b *BasicTokenStore) GetType() string {
	return b.tokenType
}

// GetToken implements TokenStore.
func (b *BasicTokenStore) GetToken(ctx context.Context) (string, error) {
	b.tokenLock.Lock()
	defer b.tokenLock.Unlock()
	if b.lastToken == "" {
		return b.refreshTokenWithoutLock(ctx)
	}
	return b.lastToken, nil
}

// RefreshToken implements TokenStore.
func (b *BasicTokenStore) RefreshToken(ctx context.Context) (string, error) {
	b.tokenLock.Lock()
	defer b.tokenLock.Unlock()
	return b.refreshTokenWithoutLock(ctx)
}

func (b *BasicTokenStore) refreshTokenWithoutLock(ctx context.Context) (string, error) {
	slog.DebugContext(ctx, fmt.Sprintf("Current token for %s is expired. Refreshing a new token", b.tokenType))
	if b.lastToken != "" {
		b.expiredTokens[b.lastToken] = struct{}{}
	}
	token, err := b.resolver.Resolve(ctx, b.expiredTokens)
	if err != nil {
		return "", err
	}
	if b.lastToken == token {
		return "", ErrNoNewTokenResolved
	}
	b.lastToken = token
	return token, err
}

var _ TokenStore = &BasicTokenStore{}
