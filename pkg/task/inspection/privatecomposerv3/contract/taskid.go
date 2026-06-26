// Copyright 2025 Google LLC
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

package privatecomposerv3_contract

import (
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
)

const PrivateComposerV3TaskIDPrefix = "privatecomposerv3/"

// InputComposerV3TenantProjectIdTaskID is the task ID for the tenant project ID of Managed Airflow 3.
var InputComposerV3TenantProjectIdTaskID = taskid.NewDefaultImplementationID[string](PrivateComposerV3TaskIDPrefix + "input-tenant-project-id")

// ComposerV3ClusterNamePrefixTaskID is the task id for the task that returns the GKE cluster name prefix used by Managed Airflow 3.
var ComposerV3ClusterNamePrefixTaskID = taskid.NewImplementationID(googlecloudk8scommon_contract.ClusterNamePrefixTaskRef, "gcp-composer-v3")
