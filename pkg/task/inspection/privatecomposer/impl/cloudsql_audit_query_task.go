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

package privatecomposer_impl

import (
	"context"
	"fmt"

	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
)

type cloudSQLAuditListLogEntriesTaskSetting struct{}

// TaskID returns the task ID for the Cloud SQL activity audit list log entries task.
func (s *cloudSQLAuditListLogEntriesTaskSetting) TaskID() taskid.TaskImplementationID[[]*log.Log] {
	return privatecomposer_contract.CloudSQLAuditLogsQueryTaskID
}

// Dependencies returns the dependent tasks required before querying Cloud SQL audit logs.
func (s *cloudSQLAuditListLogEntriesTaskSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref(),
	}
}

// Description returns the metadata description of this audit log query task.
func (s *cloudSQLAuditListLogEntriesTaskSetting) Description() *googlecloudcommon_contract.ListLogEntriesTaskDescription {
	return &googlecloudcommon_contract.ListLogEntriesTaskDescription{
		QueryName:    "Cloud Composer Tenant Project Cloud SQL Audit Logs",
		ExampleQuery: generateCloudSQLAuditExampleQuery("sample-tenant-project-tp"),
	}
}

// LogFilters returns the Cloud Logging filter queries for Cloud SQL activity audit logs in the tenant project.
func (s *cloudSQLAuditListLogEntriesTaskSetting) LogFilters(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) ([]string, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref())
	if tenantProjectID == "" {
		return []string{}, nil
	}
	return []string{generateCloudSQLAuditExampleQuery(tenantProjectID)}, nil
}

// DefaultResourceNames returns the parent resource names to query logs from.
func (s *cloudSQLAuditListLogEntriesTaskSetting) DefaultResourceNames(ctx context.Context) ([]string, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref())
	if tenantProjectID == "" {
		return []string{}, nil
	}
	return []string{fmt.Sprintf("projects/%s", tenantProjectID)}, nil
}

// TimePartitionCount returns the number of time partitions used when querying logs.
func (s *cloudSQLAuditListLogEntriesTaskSetting) TimePartitionCount(ctx context.Context) (int, error) {
	return 5, nil
}

var _ googlecloudcommon_contract.ListLogEntriesTaskSetting = (*cloudSQLAuditListLogEntriesTaskSetting)(nil)

// CloudSQLAuditLogsQueryTask executes Cloud Logging filter to fetch Cloud SQL activity audit logs in the tenant project.
var CloudSQLAuditLogsQueryTask = googlecloudcommon_contract.NewListLogEntriesTask(&cloudSQLAuditListLogEntriesTaskSetting{})

// generateCloudSQLAuditExampleQuery formats the Cloud Logging filter query for a given tenant project ID.
func generateCloudSQLAuditExampleQuery(tenantProjectID string) string {
	return fmt.Sprintf(`resource.type="cloudsql_database"
resource.labels.project_id="%s"

LOG_ID("cloudaudit.googleapis.com/activity")`, tenantProjectID)
}
