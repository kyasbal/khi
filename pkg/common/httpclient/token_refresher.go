package httpclient

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

type TokenRefresher interface {
	Refresh(ctx context.Context)
}

type NopTokenRefresher struct {
}

// Refresh implements TokenRefresher.
func (n *NopTokenRefresher) Refresh(ctx context.Context) {
}

var _ TokenRefresher = (*NopTokenRefresher)(nil)

// MultiTokenStoreRefresher implements the TokenRefresher for the clients uses multiple token store.
// This refresher will refresh token only on a single token store. It will refresh the next store on the next call.
// This will prevent requesting refresh the token that is actually not needed to update.
type MultiTokenStoreRefresher struct {
	nextRefreshIndex int
	tokenStores      []token.TokenStore
}

// Refresh implements TokenRefresher.
func (m *MultiTokenStoreRefresher) Refresh(ctx context.Context) {
	for i := 0; i < len(m.tokenStores); i++ {
		_, err := m.tokenStores[i].RefreshToken(ctx)
		if err != nil {
			slog.DebugContext(ctx, fmt.Sprintf("token store %s couldn't refresh the new token. Skipping refresh this token store", m.tokenStores[i].GetType()))
			continue
		} else {
			return
		}
	}
}

func NewMultiTokenStoreRefresher(tokenStores ...token.TokenStore) *MultiTokenStoreRefresher {
	return &MultiTokenStoreRefresher{
		nextRefreshIndex: 0,
		tokenStores:      tokenStores,
	}
}

var _ TokenRefresher = (*MultiTokenStoreRefresher)(nil)
