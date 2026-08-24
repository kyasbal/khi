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
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud/logestimator"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/gcpqueryutil"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
)

type listLogEntriesTaskSetting struct{}

// DefaultResourceNames implements [googlecloudcommon_contract.StructuredListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) DefaultResourceNames(ctx context.Context) ([]string, error) {
	logSource := coretask.GetTaskResult(ctx, privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref())
	if logSource == nil {
		return []string{}, nil
	}
	return []string{logSource.LogViewResourceName}, nil
}

// Dependencies implements [googlecloudcommon_contract.StructuredListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref(),
		privategkemaster_contract.InputPrivateGKEMasterComponentNameFilterTaskID.Ref(),
	}
}

// QueryName implements [googlecloudcommon_contract.StructuredListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) QueryName() string {
	return "Private GKE master logs"
}

// Queries implements [googlecloudcommon_contract.StructuredListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) Queries(ctx context.Context) ([]*logestimator.StructuredLogQuery, error) {
	logSource := coretask.GetTaskResult(ctx, privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref())
	componentFilter := coretask.GetTaskResult(ctx, privategkemaster_contract.InputPrivateGKEMasterComponentNameFilterTaskID.Ref())
	projectID := ""
	if logSource != nil {
		projectID = logSource.TenantProjectID
	}
	return []*logestimator.StructuredLogQuery{
		GeneratePrivateGKEMasterStructuredQuery(projectID, componentFilter),
	}, nil
}

// TaskID implements [googlecloudcommon_contract.StructuredListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) TaskID() taskid.TaskImplementationID[[]*log.Log] {
	return privategkemaster_contract.ListLogEntriesTaskID
}

// TimePartitionCount implements [googlecloudcommon_contract.StructuredListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) TimePartitionCount(ctx context.Context) (int, error) {
	return 10, nil
}

var _ googlecloudcommon_contract.StructuredListLogEntriesTaskSetting = (*listLogEntriesTaskSetting)(nil)

var listLogEntriesTask = googlecloudcommon_contract.NewStructuredListLogEntriesTask(
	&listLogEntriesTaskSetting{},
)

// GeneratePrivateGKEMasterStructuredQuery generates a structured query for Private GKE master logs.
func GeneratePrivateGKEMasterStructuredQuery(projectID string, componentFilter *gcpqueryutil.SetFilterParseResult) *logestimator.StructuredLogQuery {
	filters := []logestimator.LoggingMonitoringMatcher{
		logestimator.LogID(logestimator.NoneOf(
			"cloudaudit.googleapis.com/activity",
			"cloudaudit.googleapis.com/data_access",
			"compute.googleapis.com/shielded_vm_integrity",
			"serialconsole.googleapis.com/serial_port_1_output",
			"serialconsole.googleapis.com/serial_port_2_output",
			"serialconsole.googleapis.com/serial_port_3_output",
			"serialconsole.googleapis.com/serial_port_debug_output",
		)),
	}

	if componentFilter != nil {
		switch {
		case componentFilter.ValidationError != "":
			filters = append(filters, logestimator.Comment(fmt.Sprintf(`Failed to generate component name filter due to the validation error "%s"`, componentFilter.ValidationError)))
		case componentFilter.SubtractMode:
			if len(componentFilter.Subtractives) == 0 {
				filters = append(filters, logestimator.Comment("no master component filter"))
			} else {
				filters = append(filters, logestimator.LogID(logestimator.NoneOf(componentFilter.Subtractives...)))
			}
		default:
			if len(componentFilter.Additives) > 0 {
				filters = append(filters, logestimator.LogID(logestimator.OneOf(componentFilter.Additives...)))
			}
		}
	}

	filters = append(filters, logestimator.ResourceLabel("project_id", logestimator.Exact(projectID)))

	return &logestimator.StructuredLogQuery{
		Incomplete:    projectID == "",
		ResourceTypes: []string{"container", "gce_instance"},
		// Cloud Monitoring metrics queries are disabled because users do not have permissions
		// to query metrics in the GKE master tenant project.
		IgnoreMetricsResourceType: []string{"container", "gce_instance"},
		Filters:                   filters,
	}
}

// GeneratePrivateGKEMasterQuery formats the Cloud Logging filter query for Private GKE master logs.
func GeneratePrivateGKEMasterQuery(projectID string, componentFilter *gcpqueryutil.SetFilterParseResult) string {
	return GeneratePrivateGKEMasterStructuredQuery(projectID, componentFilter).GenerateCloudLoggingQuery()
}
