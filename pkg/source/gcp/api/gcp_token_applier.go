// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
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
	if parameters.Private.InspectionMode != nil && *parameters.Private.InspectionMode {
		iamToken, err := g.iamTokenStore.GetToken(req.Context())
		if err != nil {
			return nil, errors.New("failed to get iam token")
		}
		req.Header.Set("x-goog-iam-authorization-token", iamToken.RawToken)
	}
	return &httpclient.TokenApplyResult{
		TokenObtainedAt: tokenObtainTime,
	}, nil
}
