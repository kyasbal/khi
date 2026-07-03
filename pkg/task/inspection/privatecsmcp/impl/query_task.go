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

package privatecsmcp_impl

import (
	"context"
	"fmt"

	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
)

// LogQueryTask executes Cloud Logging filter to fetch CSM CP logs.
var LogQueryTask = googlecloudcommon_contract.NewListLogEntriesTask(&csmcpLogQueryTaskSetting{})

type csmcpLogQueryTaskSetting struct{}

// TaskID returns the ID of this task.
func (s *csmcpLogQueryTaskSetting) TaskID() taskid.TaskImplementationID[[]*log.Log] {
	return privatecsmcp_contract.LogQueryTaskID
}

// Dependencies returns the dependencies required by this task.
func (s *csmcpLogQueryTaskSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(),
		privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref(),
	}
}

// Description returns the description of this task.
func (s *csmcpLogQueryTaskSetting) Description() *googlecloudcommon_contract.ListLogEntriesTaskDescription {
	return &googlecloudcommon_contract.ListLogEntriesTaskDescription{
		QueryName: "CSM CP logs",
		ExampleQuery: `resource.type="cloud_run_revision"
AND LOG_ID("run.googleapis.com/stdout")`,
	}
}

// LogFilters returns the log queries to execute.
func (s *csmcpLogQueryTaskSetting) LogFilters(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) ([]string, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref())
	serviceName := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref())

	if tenantProjectID == "" || serviceName == "" {
		return []string{}, nil
	}

	query := fmt.Sprintf(`resource.type="cloud_run_revision"
resource.labels.project_id="%s"
resource.labels.service_name="%s"
LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr")`, tenantProjectID, serviceName)
	return []string{query}, nil
}

// DefaultResourceNames returns the resource names to query logs from.
func (s *csmcpLogQueryTaskSetting) DefaultResourceNames(ctx context.Context) ([]string, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref())
	if tenantProjectID == "" {
		return []string{}, nil
	}
	return []string{fmt.Sprintf("projects/%s", tenantProjectID)}, nil
}

// TimePartitionCount returns the number of time partitions to use.
func (s *csmcpLogQueryTaskSetting) TimePartitionCount(ctx context.Context) (int, error) {
	return 10, nil
}

var _ googlecloudcommon_contract.ListLogEntriesTaskSetting = (*csmcpLogQueryTaskSetting)(nil)
