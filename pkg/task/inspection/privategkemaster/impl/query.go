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
	"strings"

	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/enum"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
)

type listLogEntriesTaskSetting struct{}

// DefaultResourceNames implements [googlecloudcommon_contract.ListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) DefaultResourceNames(ctx context.Context) ([]string, error) {
	logSource := coretask.GetTaskResult(ctx, privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref())
	if logSource == nil {
		return []string{
			"! This field is filled after you provide the master logs link.",
		}, nil
	}
	return []string{logSource.LogViewResourceName}, nil
}

func (l *listLogEntriesTaskSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref(),
		privategkemaster_contract.InputPrivateGKEMasterComponentNameFilterTaskID.Ref(),
	}
}

// Description implements [googlecloudcommon_contract.ListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) Description() *googlecloudcommon_contract.ListLogEntriesTaskDescription {
	return &googlecloudcommon_contract.ListLogEntriesTaskDescription{
		QueryName:      "Private GKE master logs",
		DefaultLogType: enum.LogTypeControlPlaneComponent,
		ExampleQuery: `resource.type=("container" OR "gce_instance")
-log_id("cloudaudit.googleapis.com/activity")
-log_id("cloudaudit.googleapis.com/data_access")
-log_id("compute.googleapis.com/shielded_vm_integrity")
-logName: "serialconsole.googleapis.com"
resource.labels.project_id="tp-????????"`,
	}
}

// LogFilters implements [googlecloudcommon_contract.ListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) LogFilters(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) ([]string, error) {
	logSource := coretask.GetTaskResult(ctx, privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref())
	componentFilter := coretask.GetTaskResult(ctx, privategkemaster_contract.InputPrivateGKEMasterComponentNameFilterTaskID.Ref())
	projectID := "(Project ID is filled after you provide the master logs link)"
	if logSource != nil {
		projectID = logSource.TenantProjectID
	}
	componentFilterStr := ""
	if !componentFilter.SubtractMode {
		components := componentFilter.AdditivesWithQuotes()
		if len(components) > 0 {
			logIdComponents := make([]string, 0, len(components))
			for _, component := range components {
				logIdComponents = append(logIdComponents, fmt.Sprintf("log_id(%s)", component))
			}
			componentFilterStr = fmt.Sprintf("\n%s", strings.Join(logIdComponents, " OR "))
		}
	} else {
		if len(componentFilter.Subtractives) == 0 {
			componentFilterStr = "-- no master component filter"
		} else {
			components := componentFilter.SubtractivesWithQuotes()
			if len(components) > 0 {
				logIdComponents := make([]string, 0, len(components))
				for _, component := range components {
					logIdComponents = append(logIdComponents, fmt.Sprintf("-log_id(%s)", component))
				}
				componentFilterStr = fmt.Sprintf("\n%s", strings.Join(logIdComponents, "\n"))
			}
		}
	}
	return []string{fmt.Sprintf(`resource.type=("container" OR "gce_instance")
-log_id("cloudaudit.googleapis.com/activity")
-log_id("cloudaudit.googleapis.com/data_access")
-log_id("compute.googleapis.com/shielded_vm_integrity")
-logName: "serialconsole.googleapis.com"
%s
resource.labels.project_id="%s"`, componentFilterStr, projectID)}, nil
}

// TaskID implements [googlecloudcommon_contract.ListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) TaskID() taskid.TaskImplementationID[[]*log.Log] {
	return privategkemaster_contract.ListLogEntriesTaskID
}

// TimePartitionCount implements [googlecloudcommon_contract.ListLogEntriesTaskSetting].
func (l *listLogEntriesTaskSetting) TimePartitionCount(ctx context.Context) (int, error) {
	return 10, nil
}

var _ googlecloudcommon_contract.ListLogEntriesTaskSetting = (*listLogEntriesTaskSetting)(nil)

var listLogEntriesTask = googlecloudcommon_contract.NewListLogEntriesTask(
	&listLogEntriesTaskSetting{},
)
