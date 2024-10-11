package token

import (
	"context"
	"testing"
	"time"
)

func TestNewMultiTokenStoreRefresher(t *testing.T) {
	refresher := NewMultiTokenStoreRefresher(NewBasicTokenStore("foo", NewSpyTokenResolver()), NewBasicTokenStore("bar", NewSpyTokenResolver()))

	if refresher.nextStoreIndexToRefresh != 0 {
		t.Errorf("Expected refresher.nextStoreIndexToRefresh to be 0, but got %d", refresher.nextStoreIndexToRefresh)
	}
}

func TestMultiTokenStoreRefresher_Refresh(t *testing.T) {
	expireOnFuture := time.Date(2300, time.January, 1, 0, 0, 0, 0, time.UTC)
	testCase := []struct {
		name                            string
		stores                          []TokenStore
		nextStoreIndexToRefreshFirst    int
		nextStoreIndexToRefreshExpected int
		expectedRawTokens               []string
		wantErr                         bool
	}{
		{
			name: "refreshing the first token",
			stores: []TokenStore{
				NewBasicTokenStore("store1", NewSpyTokenResolver(New("foo"), New("bar"))),
				NewBasicTokenStore("store2", NewSpyTokenResolver(New("foo"), New("bar"))),
			},
			nextStoreIndexToRefreshFirst:    0,
			nextStoreIndexToRefreshExpected: 1,
			wantErr:                         false,
			expectedRawTokens: []string{
				"bar", "foo",
			},
		},
		{
			name: "refreshing the last token",
			stores: []TokenStore{
				NewBasicTokenStore("store1", NewSpyTokenResolver(New("foo"), New("bar"))),
				NewBasicTokenStore("store2", NewSpyTokenResolver(New("foo"), New("bar"))),
			},
			nextStoreIndexToRefreshFirst:    1,
			nextStoreIndexToRefreshExpected: 0,
			wantErr:                         false,
			expectedRawTokens: []string{
				"foo", "bar",
			},
		},
		{
			name: "refreshing the last token because the first token is not yet expired",
			stores: []TokenStore{
				NewBasicTokenStore("store1", NewSpyTokenResolver(NewWithExpiry("foo", expireOnFuture), New("bar"))),
				NewBasicTokenStore("store2", NewSpyTokenResolver(New("foo"), New("bar"))),
			},
			nextStoreIndexToRefreshFirst:    0,
			nextStoreIndexToRefreshExpected: 0,
			wantErr:                         false,
			expectedRawTokens: []string{
				"foo", "bar",
			},
		},
		{
			name: "return error when all store has valid tokens",
			stores: []TokenStore{
				NewBasicTokenStore("store1", NewSpyTokenResolver(NewWithExpiry("foo", expireOnFuture))),
				NewBasicTokenStore("store1", NewSpyTokenResolver(NewWithExpiry("foo", expireOnFuture))),
			},
			wantErr:                         true,
			nextStoreIndexToRefreshFirst:    0,
			nextStoreIndexToRefreshExpected: 0,
			expectedRawTokens:               []string{},
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMultiTokenStoreRefresher(tt.stores...)
			store.nextStoreIndexToRefresh = tt.nextStoreIndexToRefreshFirst
			for _, store := range tt.stores {
				_, _ = store.GetToken(context.Background())
			}

			err := store.Refresh(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Error("Expected an error but no error returned")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.nextStoreIndexToRefreshExpected != store.nextStoreIndexToRefresh {
					t.Errorf("Expected nextStoreIndexToRefresh to be %d, but got %d", tt.nextStoreIndexToRefreshExpected, store.nextStoreIndexToRefresh)
				}
				if len(tt.expectedRawTokens) != len(tt.stores) {
					t.Errorf("Expected len(expectedRawTokens) to be %d, but got %d", len(tt.expectedRawTokens), len(tt.stores))
				}

				for i, store := range tt.stores {
					token, err := store.GetToken(context.Background())

					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					if tt.expectedRawTokens[i] != token.RawToken {
						t.Errorf("Expected token.RawToken to be %s, but got %s", tt.expectedRawTokens[i], token.RawToken)
					}
				}
			}
		})
	}
}
