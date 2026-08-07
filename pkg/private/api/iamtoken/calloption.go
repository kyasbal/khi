// Copyright 2026 Google LLC
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
	"context"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typeddict"
	"google.golang.org/grpc/metadata"
)

const iamTokenKey = "x-goog-iam-authorization-token"

// IAMTokenCallOptionInjectorOption is a CallOptionInjectorOption that injects an IAM token into API calls.
type IAMTokenCallOptionInjectorOption struct {
	containerSpecificTokens *typeddict.TypedDict[string]
}

func NewInjector() *IAMTokenCallOptionInjectorOption {
	return &IAMTokenCallOptionInjectorOption{
		containerSpecificTokens: typeddict.NewTypedDict[string](),
	}
}

// ApplyToCallContext implements googlecloud.CallOptionInjectorOption.
// It appends the IAM token to the outgoing context metadata for gRPC calls.
func (o *IAMTokenCallOptionInjectorOption) ApplyToCallContext(ctx context.Context, container googlecloud.ResourceContainer) context.Context {
	token := o.getTokenFor(container)
	if token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, iamTokenKey, token)
}

// ApplyToRawHTTPHeader implements googlecloud.CallOptionInjectorOption.
// It sets the IAM token in the HTTP header for REST calls.
func (o *IAMTokenCallOptionInjectorOption) ApplyToRawHTTPHeader(header http.Header, container googlecloud.ResourceContainer) {
	token := o.getTokenFor(container)
	if token == "" {
		return
	}
	header.Set(iamTokenKey, token)
}

// SetTokenFor sets the IAM token for the given container.
func (o *IAMTokenCallOptionInjectorOption) SetTokenFor(container googlecloud.ResourceContainer, token string) {
	typeddict.Set(o.containerSpecificTokens, container.Identifier(), token)
}

// HasTokenFor checks if the given container has an IAM token.
func (o *IAMTokenCallOptionInjectorOption) HasTokenFor(container googlecloud.ResourceContainer) bool {
	if container == nil {
		return false
	}
	ci := container.Identifier()
	_, found := typeddict.Get(o.containerSpecificTokens, ci)
	return found
}

// RegisteredProjectIDs returns a list of Google Cloud project IDs that have IAM tokens registered.
func (o *IAMTokenCallOptionInjectorOption) RegisteredProjectIDs() []string {
	if o == nil || o.containerSpecificTokens == nil {
		return nil
	}
	keys := o.containerSpecificTokens.Keys()
	var projectIDs []string
	for _, key := range keys {
		if strings.HasPrefix(key, "projects/") {
			projectIDs = append(projectIDs, strings.TrimPrefix(key, "projects/"))
		}
	}
	return projectIDs
}

func (o *IAMTokenCallOptionInjectorOption) getTokenFor(container googlecloud.ResourceContainer) string {
	if container == nil {
		return ""
	}
	ci := container.Identifier()
	token, found := typeddict.Get(o.containerSpecificTokens, ci)
	if !found {
		return ""
	}
	return token
}

var _ googlecloud.CallOptionInjectorOption = (*IAMTokenCallOptionInjectorOption)(nil)
