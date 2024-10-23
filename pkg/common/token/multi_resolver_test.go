package token

import (
	"context"
	"testing"
)

func TestMultiTokenResolver_Resolve(t *testing.T) {
	testCases := []struct {
		name          string
		resolvers     []TokenResolver
		wantErr       bool
		expectedToken *Token
	}{
		{
			name:          "without any resolvers",
			resolvers:     make([]TokenResolver, 0),
			wantErr:       true,
			expectedToken: nil,
		},
		{
			name: "with the first successful resolver",
			resolvers: []TokenResolver{
				NewSpyTokenResolver(New("foo")),
			},
			wantErr:       false,
			expectedToken: New("foo"),
		},
		{
			name: "with the errornous resolver and successful resolver",
			resolvers: []TokenResolver{
				NewMockErrorTokenResolver(),
				NewSpyTokenResolver(New("foo")),
			},
			wantErr:       false,
			expectedToken: New("foo"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewMultiTokenResolver(tt.resolvers...)

			token, err := resolver.Resolve(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("got nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("got %v, want nil", err)
				}
				if token.RawToken != tt.expectedToken.RawToken {
					t.Errorf("got raw token %s, want %s", token.RawToken, tt.expectedToken.RawToken)
				}
				if !token.ValidAtLeastUntil.Equal(tt.expectedToken.ValidAtLeastUntil) {
					t.Errorf("got expiry %v, want %v", token.ValidAtLeastUntil, tt.expectedToken.ValidAtLeastUntil)
				}
			}
		})
	}
}
