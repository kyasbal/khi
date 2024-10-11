package accesstoken

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

type MDSResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type MDSTokenResolver struct {
	client *httpclient.JSONReponseHttpClient[MDSResponse]
}

func NewMetadataServerAccessTokenResolver(client *httpclient.JSONReponseHttpClient[MDSResponse]) *MDSTokenResolver {
	return &MDSTokenResolver{
		client: client,
	}
}

// Resolve implements token.TokenResolver.
func (m *MDSTokenResolver) Resolve(ctx context.Context) (*token.Token, error) {
	slog.InfoContext(ctx, `Environment variable "GCP_ACCESS_TOKEN" was not found. Trying to get access token from metadata server`)
	req, err := http.NewRequest("GET", metadataServerAddress, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Metadata-Flavor", "Google")
	response, _, err := m.client.DoWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	if response.AccessToken != "" {
		return token.NewWithExpiry(response.AccessToken, time.Now().Add(time.Duration(response.ExpiresIn-1)*time.Second)), nil
	}
	return nil, fmt.Errorf("failed to get access token")
}

var metadataServerAddress = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"

var _ token.TokenResolver = (*MDSTokenResolver)(nil)
