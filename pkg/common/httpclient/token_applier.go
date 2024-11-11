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

package httpclient

import (
	"context"
	"net/http"
	"time"
)

// TokenApplyResult contains the operation result of ApplyCurrentToken.
// This contains TokenApplyTime to purge related token only when the token was expired.
// Because the time to purge token and obtaining token is different and these are used in parallel, it can purge newly refreshed token without checking the time to issue the token.
type TokenApplyResult struct {
	TokenObtainedAt time.Time
}

type TokenApplier interface {
	ApplyCurrentToken(ctx context.Context, req *http.Request) (*TokenApplyResult, error)
}

type NopTokenApplier struct{}

var _ TokenApplier = (*NopTokenApplier)(nil)

func (n *NopTokenApplier) ApplyCurrentToken(ctx context.Context, req *http.Request) (*TokenApplyResult, error) {
	return &TokenApplyResult{
		TokenObtainedAt: time.Now(),
	}, nil
}
