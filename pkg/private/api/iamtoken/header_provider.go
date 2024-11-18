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
