package iamtoken

import (
	"errors"
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

// GCPIAMTokenProvider is an implementation of HTTPHeaderProvider for IAM token.
type GCPIAMTokenProvider struct {
	IAMToken token.TokenStore
}

func NewHeaderProvider(tokenStore token.TokenStore) *GCPIAMTokenProvider {
	return &GCPIAMTokenProvider{
		IAMToken: tokenStore,
	}
}

// AddHeader implements httpclient.HTTPHeaderProvider.
func (a *GCPIAMTokenProvider) AddHeader(req *http.Request) error {
	token, err := a.IAMToken.GetToken(req.Context())
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("iam token is empty")
	}
	req.Header.Set("x-goog-iam-authorization-token", token.RawToken)
	return nil
}

var _ httpclient.HTTPHeaderProvider = (*GCPIAMTokenProvider)(nil)
