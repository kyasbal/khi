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

package privatecomposer_impl

import (
	"context"

	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
)

// ComposerV3ClusterNamePrefixTask is the task that returns the Managed Airflow 3 cluster name prefix.
var ComposerV3ClusterNamePrefixTask = coretask.NewTask(
	privatecomposer_contract.ComposerV3ClusterNamePrefixTaskID,
	[]taskid.UntypedTaskReference{},
	func(ctx context.Context) (googlecloudk8scommon_contract.ClusterPrefixPolicy, error) {
		return googlecloudk8scommon_contract.ClusterPrefixPolicy{}, nil
	},
	inspectioncore_contract.InspectionTypeLabel(privatecomposer_contract.InspectionTypeId),
)
