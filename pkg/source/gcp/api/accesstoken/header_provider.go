package accesstoken

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

// GCPAccessTokenHeaderProvider is an implementation of HTTPHeaderProvider for access token.
type GCPAccessTokenHeaderProvider struct {
	AccessToken token.TokenStore
}

func NewHeaderProvider(tokenStore token.TokenStore) *GCPAccessTokenHeaderProvider {
	return &GCPAccessTokenHeaderProvider{
		AccessToken: tokenStore,
	}
}

// AddHeader implements httpclient.HTTPHeaderProvider.
func (a *GCPAccessTokenHeaderProvider) AddHeader(req *http.Request) error {
	token, err := a.AccessToken.GetToken(req.Context())
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("access token is empty")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.RawToken))
	return nil
}

var _ httpclient.HTTPHeaderProvider = (*GCPAccessTokenHeaderProvider)(nil)
