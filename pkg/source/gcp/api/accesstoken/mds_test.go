package accesstoken

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil"
	"github.com/google/go-cmp/cmp"
)

type mockMDSResponseHttpClient struct {
	response string
	err      error
}

func newMockMDSResponseHttpClient(response string, err error) *mockMDSResponseHttpClient {
	return &mockMDSResponseHttpClient{
		response: response,
		err:      err,
	}
}

// DoWithContext implements httpclient.HttpClient.
func (m *mockMDSResponseHttpClient) DoWithContext(ctx context.Context, request *http.Request) (*http.Response, error) {
	if m.err != nil {
		return testutil.ResponseFromString(http.StatusOK, ""), m.err
	}
	return testutil.ResponseFromString(http.StatusOK, m.response), nil
}

var _ httpclient.HttpClient[*http.Response] = (*mockMDSResponseHttpClient)(nil)

func TestMDSTokenResolver(t *testing.T) {
	tests := []struct {
		name             string
		client           *httpclient.JSONReponseHttpClient[MDSResponse]
		wantExpireAround time.Time
		want             string
		wantErr          bool
		expiredToken     map[string]interface{}
	}{
		{
			name:             "MDSTokenResolver should return the token from the metadata server",
			client:           httpclient.NewJsonResponseHttpClient[MDSResponse](newMockMDSResponseHttpClient("{ \"access_token\": \"test-fake-token\",\"expires_in\": 180}", nil)),
			want:             "test-fake-token",
			wantErr:          false,
			wantExpireAround: time.Now().Add(time.Minute * 3),
		},
		{
			name:    "MDSTokenResolver should return error if client returns error",
			client:  httpclient.NewJsonResponseHttpClient[MDSResponse](newMockMDSResponseHttpClient("", errors.New("test-error"))),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetadataServerAccessTokenResolver(tt.client)
			got, err := m.Resolve(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("MDSTokenResolver.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, got.RawToken); diff != "" {
					t.Errorf("MDSTokenResolver.Resolve() mismatch (-want +got):\n%s", diff)
				}
				if !got.ValidAtLeastUntil.Before(tt.wantExpireAround.Add(time.Second*10)) || !got.ValidAtLeastUntil.After(tt.wantExpireAround.Add(-time.Second*10)) {
					t.Errorf("MDSTokenResolver.Resolve() returns expire time far different from the expected time. want:%v, got:%v", tt.wantExpireAround, got.ValidAtLeastUntil)
				}
			}
		})
	}
}
