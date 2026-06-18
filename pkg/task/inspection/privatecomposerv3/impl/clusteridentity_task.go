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

package privatecomposerv3_impl

import (
	"context"

	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudclustercomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudclustercomposer/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposerv3_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposerv3/contract"
)

var ClusterIdentityTask = inspectiontaskbase.NewInspectionTask(taskid.NewImplementationID(googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(), privatecomposerv3_contract.InspectionTypeId), []taskid.UntypedTaskReference{
	privatecomposerv3_contract.InputComposerV3TenantProjectIdTaskID.Ref(),
	googlecloudk8scommon_contract.InputClusterNameTaskID.Ref(),
	googlecloudcommon_contract.InputLocationsTaskID.Ref(),
	googlecloudk8scommon_contract.ClusterNamePrefixTaskRef,
}, func(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) (googlecloudk8scommon_contract.GoogleCloudClusterIdentity, error) {
	projectId := coretask.GetTaskResult(ctx, privatecomposerv3_contract.InputComposerV3TenantProjectIdTaskID.Ref())
	clusterName := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.InputClusterNameTaskID.Ref())
	location := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputLocationsTaskID.Ref())
	prefixPolicy := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterNamePrefixTaskRef)

	return googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
		ProjectID:    projectId,
		Location:     location,
		ClusterName:  clusterName,
		PrefixPolicy: prefixPolicy,
	}, nil
},
	inspectioncore_contract.InspectionTypeLabel(privatecomposerv3_contract.InspectionTypeId),
	coretask.WithSelectionPriority(100),
)

// ComposerClusterIdentityTask is an override for googlecloudclustercomposer_contract.ClusterIdentityTaskID
// that ensures Composer queries read from the original customer project ID rather than the tenant project ID.
var ComposerClusterIdentityTask = inspectiontaskbase.NewInspectionTask(taskid.NewImplementationID(googlecloudclustercomposer_contract.ClusterIdentityTaskID.Ref(), privatecomposerv3_contract.InspectionTypeId), []taskid.UntypedTaskReference{
	googlecloudcommon_contract.InputProjectIdTaskID.Ref(),
	googlecloudcommon_contract.InputLocationsTaskID.Ref(),
}, func(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) (googlecloudk8scommon_contract.GoogleCloudClusterIdentity, error) {
	projectId := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputProjectIdTaskID.Ref())
	location := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputLocationsTaskID.Ref())

	return googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
		ProjectID: projectId,
		Location:  location,
	}, nil
},
	inspectioncore_contract.InspectionTypeLabel(privatecomposerv3_contract.InspectionTypeId),
	coretask.WithSelectionPriority(100),
)
