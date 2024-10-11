package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

type GCPTokenApplier struct {
	accessTokenStore token.TokenStore
	iamTokenStore    token.TokenStore
}

func NewGCPTokenApplier(accessTokenStore token.TokenStore, iamTokenStore token.TokenStore) *GCPTokenApplier {
	return &GCPTokenApplier{
		accessTokenStore: accessTokenStore,
		iamTokenStore:    iamTokenStore,
	}
}

var _ httpclient.TokenApplier = (*GCPTokenApplier)(nil)

func (g *GCPTokenApplier) ApplyCurrentToken(ctx context.Context, req *http.Request) (*httpclient.TokenApplyResult, error) {
	tokenObtainTime := time.Now()
	token, err := g.accessTokenStore.GetToken(req.Context())
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("access token is empty")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.RawToken))
	if use, found := os.LookupEnv("USE_IAM_TOKEN"); found && use != "false" {
		iamToken, err := g.iamTokenStore.GetToken(req.Context())
		if err != nil {
			return nil, errors.New("USE_IAM_TOKEN environment variable is given but no IAM token resolved")
		}
		req.Header.Set("x-goog-iam-authorization-token", iamToken.RawToken)
	}
	return &httpclient.TokenApplyResult{
		TokenObtainedAt: tokenObtainTime,
	}, nil
}
