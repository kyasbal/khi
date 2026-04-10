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
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
	privatecomposerv3_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposerv3/contract"
)

// Register registers all privatecomposerv3 inspection tasks to the registry.
func Register(registry coreinspection.InspectionTaskRegistry) error {
	err := registry.AddInspectionType(privatecomposerv3_contract.ComposerV3InspectionType)
	if err != nil {
		return err
	}
	scoped := coreinspection.NewScopedRegistry(registry, inspectioncore_contract.InspectionTypeLabelSelector(
		map[string]string{
			inspectioncore_contract.InspectionTypeLabelKeyLogSource:      "cloud_logging",
			inspectioncore_contract.InspectionTypeLabelKeyEnvironment:    "googlecloud",
			inspectioncore_contract.InspectionTypeLabelKeyBasePlatform:   "kubernetes",
			googlecloudcommon_contract.InspectionTypeLabelKeyClusterType: "gke",
			googlecloudcommon_contract.InspectionTypeLabelKeyProduct:     "composer",
			privatecommon_contract.InspectionTypeLabelKeyPrivate:         "true",
		}))
	return coretask.RegisterTasks(scoped,
		InputComposerV3TenantProjectIdTask,
		ClusterIdentityTask,
		ComposerClusterIdentityTask,
		AutocompleteComposerClusterNamesTask,
		ComposerV3ClusterNamePrefixTask,
	)
}
