package token

import (
	"context"
	"fmt"
	"log/slog"
)

type TokenRefresher interface {
	Refresh(ctx context.Context) error
}

type NopTokenRefresher struct {
}

// Refresh implements TokenRefresher.
func (n *NopTokenRefresher) Refresh(ctx context.Context) error {
	return nil
}

var _ TokenRefresher = (*NopTokenRefresher)(nil)

// MultiTokenStoreRefresher implements the TokenRefresher to refresh tokens in multiple stores.
// This rotate the next token store to be refreshed but it will ignore if the token is assured to be alive from the expiry time.
type MultiTokenStoreRefresher struct {
	tokenStores             []TokenStore
	nextStoreIndexToRefresh int
}

// Refresh implements TokenRefresher.
func (m *MultiTokenStoreRefresher) Refresh(ctx context.Context) error {
	for i := 0; i < len(m.tokenStores); i++ {
		store := m.tokenStores[m.nextStoreIndexToRefresh]
		m.nextStoreIndexToRefresh = (m.nextStoreIndexToRefresh + 1) % len(m.tokenStores)
		if store.IsTokenValidityAssured(ctx) {
			slog.DebugContext(ctx, fmt.Sprintf("Token for %s shouldn't be expired yet. ignoring token refresh.", store.GetType()))
			continue
		}
		err := store.RefreshToken(ctx)
		if err == nil {
			return nil
		}
		slog.DebugContext(ctx, fmt.Sprintf("token store %s couldn't refresh the new token. Skipping refresh this token store", m.tokenStores[i].GetType()))
	}
	return fmt.Errorf("there were no token store to refresh its token")
}

func NewMultiTokenStoreRefresher(tokenStores ...TokenStore) *MultiTokenStoreRefresher {
	return &MultiTokenStoreRefresher{
		tokenStores: tokenStores,
	}
}

var _ TokenRefresher = (*MultiTokenStoreRefresher)(nil)
