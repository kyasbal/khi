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

package privategkemaster_impl

import (
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
)

// Register registers the private GKE master tasks.
func Register(registry coreinspection.InspectionTaskRegistry) error {
	scopedWithLogs := coreinspection.NewScopedRegistry(registry, inspectioncore_contract.InspectionTypeLabelSelector(map[string]string{
		inspectioncore_contract.InspectionTypeLabelKeyLogSource:      "cloud_logging",
		inspectioncore_contract.InspectionTypeLabelKeyEnvironment:    "googlecloud",
		inspectioncore_contract.InspectionTypeLabelKeyBasePlatform:   "kubernetes",
		googlecloudcommon_contract.InspectionTypeLabelKeyClusterType: "gke",
	}))
	if err := coretask.RegisterTasks(scopedWithLogs, listLogEntriesTask); err != nil {
		return err
	}
	scoped := coreinspection.NewScopedRegistry(registry, inspectioncore_contract.InspectionTypeLabelSelector(map[string]string{
		inspectioncore_contract.InspectionTypeLabelKeyEnvironment:    "googlecloud",
		inspectioncore_contract.InspectionTypeLabelKeyBasePlatform:   "kubernetes",
		googlecloudcommon_contract.InspectionTypeLabelKeyClusterType: "gke",
	}))
	return coretask.RegisterTasks(scoped,
		InputGKEMasterLogSourceTask,
		InputPrivateGKEMasterComponentNameFilterTask,
		logIngesterTask,
		CommonFieldSetReaderTask,
		schedulerLogFilterTask,
		schedulerLogFieldSetReaderTask,
		schedulerGrouperTask,
		schedulerLogToTimelineMapperTask,
		controllerManagerLogFilterTask,
		controllerManagerLogFieldSetReaderTask,
		controllerManagerGrouperTask,
		controllerManagerLogToTimelineMapperTask,
		otherLogFilterTask,
		otherLogFieldSetReaderTask,
		otherGrouperTask,
		otherLogToTimelineMapperTask,
		KubeletLogFilterTask,
		KubeletLogGroupTask,
		KubeletLogLogToTimelineMapperTask,
		ContainerdLogFilterTask,
		ContainerdLogGroupTask,
		ContainerIDDiscoveryTask,
		PodSandboxIDDiscoveryTask,
		ContainerdLogLogToTimelineMapperTask,
		TailTask,
	)
}
