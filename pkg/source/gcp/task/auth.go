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

package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/api/iamtoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api/accesstoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const GCPApiClientTaskId = GCPPrefix + "api-client"

var GCPApiClientTask = task.NewProcessorTask(GCPApiClientTaskId,
	[]string{},
	func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		return api.NewGCPClient(accesstoken.DefaultAccessTokenStore, iamtoken.DefaultIAMTokenStore, *parameters.Auth.QuotaProjectID)
	})

func GetGCPApiClientFromTaskVariable(v *task.VariableSet) (api.GCPClient, error) {
	return task.GetTypedVariableFromTaskVariable[api.GCPClient](v, GCPApiClientTaskId, nil)
}
