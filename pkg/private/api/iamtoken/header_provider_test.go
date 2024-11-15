package iamtoken

import (
	"net/http"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

func TestGCPAccessTokenHeaderProvider(t *testing.T) {
	testCases := []struct {
		name          string
		tokenStore    token.TokenStore
		wantIAMHeader string
		wantError     bool
	}{
		{
			name:          "success",
			tokenStore:    token.NewBasicTokenStore("test", token.NewSpyTokenResolver(token.New("test-token"))),
			wantIAMHeader: "test-token",
			wantError:     false,
		},
		{
			name:          "empty token",
			tokenStore:    token.NewBasicTokenStore("test", token.NewMockErrorTokenResolver()),
			wantIAMHeader: "",
			wantError:     true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p := NewHeaderProvider(tt.tokenStore)
			req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			err = p.AddHeader(req)
			if (err != nil) != tt.wantError {
				t.Errorf("GCPIAMTokenProvider.AddHeader() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError {
				gotAuthorization := req.Header.Get("x-goog-iam-authorization-token")
				if gotAuthorization != tt.wantIAMHeader {
					t.Errorf("GCPIAMTokenProvider.AddHeader() = %v, want %v", gotAuthorization, tt.wantIAMHeader)
				}
			}
		})
	}

}
