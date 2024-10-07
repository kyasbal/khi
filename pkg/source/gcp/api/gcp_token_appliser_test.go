package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
	"github.com/google/go-cmp/cmp"
)

type mockTokenStore struct {
	Token string
	Err   error
	Type  string
}

// GetType implements token.TokenStore.
func (m *mockTokenStore) GetType() string {
	return m.Type
}

// Digest implements token.TokenStore.
func (m *mockTokenStore) Digest() string {
	return ""
}

// MarkTokenExpired implements token.TokenStore.
func (m *mockTokenStore) RefreshToken(ctx context.Context) (string, error) {
	return m.Token, nil
}

func (m *mockTokenStore) GetToken(ctx context.Context) (string, error) {
	return m.Token, m.Err
}

var _ token.TokenStore = (*mockTokenStore)(nil)

func TestGCPTokenApplier_ApplyCurrentToken(t *testing.T) {
	tests := []struct {
		name               string
		accessToken        string
		accessTokenErr     error
		iamToken           string
		iamTokenErr        error
		wantAuthorization  string
		wantIamTokenHeader string
		wantErr            bool
		before             func()
		after              func()
	}{
		{
			name:               "ApplyCurrentToken should set access token and iam token to the request header",
			accessToken:        "test-access-token",
			accessTokenErr:     nil,
			iamToken:           "test-iam-token",
			iamTokenErr:        nil,
			wantAuthorization:  "Bearer test-access-token",
			wantIamTokenHeader: "test-iam-token",
			wantErr:            false,
			before: func() {
				os.Setenv("USE_IAM_TOKEN", "true")
			},
			after: func() {
				os.Unsetenv("USE_IAM_TOKEN")
			},
		},
		{
			name:               "ApplyCurrentToken should return error if access token is empty",
			accessToken:        "",
			accessTokenErr:     nil,
			iamToken:           "test-iam-token",
			iamTokenErr:        nil,
			wantAuthorization:  "",
			wantIamTokenHeader: "",
			wantErr:            true,
		},
		{
			name:               "ApplyCurrentToken should return error if access token store returns error",
			accessToken:        "",
			accessTokenErr:     errors.New("test-error"),
			iamToken:           "test-iam-token",
			iamTokenErr:        nil,
			wantAuthorization:  "",
			wantIamTokenHeader: "",
			wantErr:            true,
			before: func() {
				os.Setenv("USE_IAM_TOKEN", "true")
			},
			after: func() {
				os.Unsetenv("USE_IAM_TOKEN")
			},
		},
		{
			name:               "ApplyCurrentToken should not set iam token if iam token is empty",
			accessToken:        "test-access-token",
			accessTokenErr:     nil,
			iamToken:           "",
			iamTokenErr:        nil,
			wantAuthorization:  "Bearer test-access-token",
			wantIamTokenHeader: "",
			wantErr:            false,
			before: func() {
				os.Setenv("USE_IAM_TOKEN", "true")
			},
			after: func() {
				os.Unsetenv("USE_IAM_TOKEN")
			},
		},
		{
			name:               "ApplyCurrentToken should not set iam token if iam token store returns error",
			accessToken:        "test-access-token",
			accessTokenErr:     nil,
			iamToken:           "",
			iamTokenErr:        errors.New("test-error"),
			wantAuthorization:  "Bearer test-access-token",
			wantIamTokenHeader: "",
			wantErr:            false,
			before: func() {
				os.Setenv("USE_IAM_TOKEN", "true")
			},
			after: func() {
				os.Unsetenv("USE_IAM_TOKEN")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			if tt.after != nil {
				defer tt.after()
			}
			accessTokenStore := &mockTokenStore{
				Type:  "accesstoken",
				Token: tt.accessToken,
				Err:   tt.accessTokenErr,
			}
			iamTokenStore := &mockTokenStore{
				Type:  "iamtoken",
				Token: tt.iamToken,
				Err:   tt.iamTokenErr,
			}
			g := &GCPTokenApplier{
				accessTokenStore: accessTokenStore,
				iamTokenStore:    iamTokenStore,
			}
			req, _ := http.NewRequest("GET", "https://example.com", nil)
			to, err := g.ApplyCurrentToken(context.Background(), req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GCPTokenApplier.ApplyCurrentToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				gotAuthorization := req.Header.Get("Authorization")
				if diff := cmp.Diff(tt.wantAuthorization, gotAuthorization); diff != "" {
					t.Errorf("GCPTokenApplier.ApplyCurrentToken() mismatch (-want +got):\n%s", diff)
				}
				gotIamTokenHeader := req.Header.Get("x-goog-iam-authorization-token")
				if diff := cmp.Diff(tt.wantIamTokenHeader, gotIamTokenHeader); diff != "" {
					t.Errorf("GCPTokenApplier.ApplyCurrentToken() mismatch (-want +got):\n%s", diff)
				}
				if time.Since(to.TokenObtainedAt) < 0 {
					t.Errorf("The time to obtain token is pointing the future")
				}
			}
		})
	}
}
