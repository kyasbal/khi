package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
	"github.com/google/go-cmp/cmp"
)

func wrapPointer[T any](t T) *T {
	return &t
}

func TestGCPTokenApplier(t *testing.T) {
	testCases := []struct {
		name               string
		accessTokenError   bool
		iamTokenError      bool
		wantAuthorization  string
		wantIamTokenHeader string
		wantErr            bool
		before             func()
		after              func()
	}{
		{
			name:               "with accesstoken and iamtoken",
			accessTokenError:   false,
			iamTokenError:      false,
			wantAuthorization:  "Bearer accesstoken1",
			wantIamTokenHeader: "iamtoken1",
			wantErr:            false,
			before: func() {
				parameters.Private.InspectionMode = wrapPointer(true)
			},
			after: func() {
				parameters.Private.InspectionMode = nil
			},
		},
		{
			name:               "with accesstoken only",
			accessTokenError:   false,
			iamTokenError:      true,
			wantAuthorization:  "Bearer accesstoken1",
			wantIamTokenHeader: "",
			wantErr:            false,
			before: func() {
				parameters.Private.InspectionMode = wrapPointer(false)
			},
			after: func() {
				parameters.Private.InspectionMode = nil
			},
		},
		{
			name:               "with accesstoken and iamtoken but USE_IAM_TOKEN is unset",
			accessTokenError:   false,
			iamTokenError:      false,
			wantAuthorization:  "Bearer accesstoken1",
			wantIamTokenHeader: "",
			wantErr:            false,
			before: func() {
				parameters.Private.InspectionMode = wrapPointer(false)
			},
			after: func() {
				parameters.Private.InspectionMode = nil
			},
		},
		{
			name:               "USE_IAM_TOKEN is set but iamtoken returns an error",
			accessTokenError:   false,
			iamTokenError:      true,
			wantAuthorization:  "Bearer accesstoken1",
			wantIamTokenHeader: "",
			wantErr:            true,
			before: func() {
				parameters.Private.InspectionMode = wrapPointer(true)
			},
			after: func() {
				parameters.Private.InspectionMode = nil
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var accessTokenStore token.TokenStore
			var iamTokenStore token.TokenStore
			if tt.accessTokenError {
				accessTokenStore = token.NewBasicTokenStore("accesstoken", token.NewMockErrorTokenResolver())
			} else {
				accessTokenStore = token.NewBasicTokenStore("accesstoken", token.NewSpyTokenResolver(token.New("accesstoken1")))
			}
			if tt.iamTokenError {
				iamTokenStore = token.NewBasicTokenStore("iamtoken", token.NewMockErrorTokenResolver())
			} else {
				iamTokenStore = token.NewBasicTokenStore("iamtoken", token.NewSpyTokenResolver(token.New("iamtoken1")))
			}
			tt.before()
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
			tt.after()
		})
	}
}
