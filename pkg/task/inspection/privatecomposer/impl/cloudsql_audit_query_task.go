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

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud/logestimator"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
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

// QueryName returns the human-readable name of the query task.
func (s *cloudSQLAuditListLogEntriesTaskSetting) QueryName() string {
	return "Cloud Composer Tenant Project Cloud SQL Audit Logs"
}

// Queries returns the list of structured log queries for Cloud SQL audit logs.
func (s *cloudSQLAuditListLogEntriesTaskSetting) Queries(ctx context.Context) ([]*logestimator.StructuredLogQuery, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref())
	return []*logestimator.StructuredLogQuery{
		GenerateCloudSQLAuditStructuredQuery(tenantProjectID),
	}, nil
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

var _ googlecloudcommon_contract.StructuredListLogEntriesTaskSetting = (*cloudSQLAuditListLogEntriesTaskSetting)(nil)

// CloudSQLAuditLogsQueryTask executes Cloud Logging filter to fetch Cloud SQL activity audit logs in the tenant project.
var CloudSQLAuditLogsQueryTask = googlecloudcommon_contract.NewStructuredListLogEntriesTask(&cloudSQLAuditListLogEntriesTaskSetting{})

// GenerateCloudSQLAuditStructuredQuery generates a structured query for Cloud SQL activity audit logs in the tenant project.
func GenerateCloudSQLAuditStructuredQuery(tenantProjectID string) *logestimator.StructuredLogQuery {
	filters := []logestimator.LoggingMonitoringMatcher{
		logestimator.ResourceLabel("project_id", logestimator.Exact(tenantProjectID)),
		logestimator.LogID(logestimator.Exact("cloudaudit.googleapis.com/activity")),
	}

	return &logestimator.StructuredLogQuery{
		Incomplete:    tenantProjectID == "",
		ResourceTypes: []string{"cloudsql_database"},
		Filters:       filters,
	}
}

// GenerateCloudSQLAuditExampleQuery formats the Cloud Logging filter query for a given tenant project ID.
func GenerateCloudSQLAuditExampleQuery(tenantProjectID string) string {
	return GenerateCloudSQLAuditStructuredQuery(tenantProjectID).GenerateCloudLoggingQuery()
}
