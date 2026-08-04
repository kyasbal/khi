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
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
)

// Register registers all privatecomposer inspection tasks to the registry.
func Register(registry coreinspection.InspectionTaskRegistry) error {
	err := registry.AddInspectionType(privatecomposer_contract.ComposerV3InspectionType)
	if err != nil {
		return err
	}

	scopedCloudLogging := coreinspection.NewScopedRegistry(registry, inspectioncore_contract.InspectionTypeLabelSelector(
		map[string]string{
			inspectioncore_contract.InspectionTypeLabelKeyLogSource:      "cloud_logging",
			inspectioncore_contract.InspectionTypeLabelKeyEnvironment:    "googlecloud",
			inspectioncore_contract.InspectionTypeLabelKeyBasePlatform:   "kubernetes",
			googlecloudcommon_contract.InspectionTypeLabelKeyClusterType: "gke",
			googlecloudcommon_contract.InspectionTypeLabelKeyProduct:     "composer",
		}))

	scopedComposerV3 := coreinspection.NewScopedRegistry(registry, inspectioncore_contract.InspectionTypeLabelSelector(
		map[string]string{
			inspectioncore_contract.InspectionTypeLabelKeyLogSource:      "cloud_logging",
			inspectioncore_contract.InspectionTypeLabelKeyEnvironment:    "googlecloud",
			inspectioncore_contract.InspectionTypeLabelKeyBasePlatform:   "kubernetes",
			googlecloudcommon_contract.InspectionTypeLabelKeyClusterType: "gke",
			googlecloudcommon_contract.InspectionTypeLabelKeyProduct:     "composer",
			privatecommon_contract.InspectionTypeLabelKeyPrivate:         "true",
		}))

	if err := coretask.RegisterTasks(scopedComposerV3,
		ClusterIdentityTask,
		ComposerClusterIdentityTask,
		AutocompleteComposerClusterNamesTask,
		ComposerV3ClusterNamePrefixTask,
	); err != nil {
		return err
	}

	return coretask.RegisterTasks(scopedCloudLogging,
		InputComposerTenantProjectIdTask,
		CloudSQLLogsQueryTask,
		CloudSQLLogsFieldSetReadTask,
		CloudSQLLogsIngesterTask,
		CloudSQLLogsGrouperTask,
		CloudSQLLogsTimelineMapperTask,
		CloudSQLAuditLogsQueryTask,
		CloudSQLAuditLogsFieldSetReadTask,
		CloudSQLAuditLogsIngesterTask,
		CloudSQLAuditLogsGrouperTask,
		CloudSQLAuditLogsTimelineMapperTask,
	)
}
