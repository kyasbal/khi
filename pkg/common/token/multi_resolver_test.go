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
					t.Error("Expected an error but no error returned")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token.RawToken != tt.expectedToken.RawToken {
					t.Errorf("Unexpected token.RawToken: got %s, want %s", token.RawToken, tt.expectedToken.RawToken)
				}
				if !token.ValidAtLeastUntil.Equal(tt.expectedToken.ValidAtLeastUntil) {
					t.Errorf("Unexpected token.ValidAtLeastUntil: got %v, want %v", token.ValidAtLeastUntil, tt.expectedToken.ValidAtLeastUntil)
				}
			}
		})
	}
}
