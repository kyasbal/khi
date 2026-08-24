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
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateCloudSQLAuditStructuredQuery(t *testing.T) {
	testCases := []struct {
		name                   string
		tenantProjectID        string
		wantQuery              string
		wantMetricFilters      []string
		wantSupportMetricsFlag bool
		wantIncomplete         bool
	}{
		{
			name:            "valid tenant project id",
			tenantProjectID: "my-tenant-project-tp",
			wantQuery: `resource.type="cloudsql_database"
resource.labels.project_id="my-tenant-project-tp"
LOG_ID("cloudaudit.googleapis.com/activity")`,
			wantMetricFilters: []string{
				`metric.type = "logging.googleapis.com/log_entry_count" AND resource.type = "cloudsql_database" AND resource.labels.project_id = "my-tenant-project-tp" AND metric.labels.log = "cloudaudit.googleapis.com/activity"`,
			},
			wantSupportMetricsFlag: true,
			wantIncomplete:         false,
		},
		{
			name:            "empty tenant project id",
			tenantProjectID: "",
			wantQuery: `resource.type="cloudsql_database"
resource.labels.project_id=""
LOG_ID("cloudaudit.googleapis.com/activity")`,
			wantMetricFilters: []string{
				`metric.type = "logging.googleapis.com/log_entry_count" AND resource.type = "cloudsql_database" AND resource.labels.project_id = "" AND metric.labels.log = "cloudaudit.googleapis.com/activity"`,
			},
			wantSupportMetricsFlag: true,
			wantIncomplete:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sq := GenerateCloudSQLAuditStructuredQuery(tc.tenantProjectID)
			gotQuery := sq.GenerateCloudLoggingQuery()
			if diff := cmp.Diff(tc.wantQuery, gotQuery); diff != "" {
				t.Errorf("GenerateCloudLoggingQuery() mismatch (-want +got):\n%s", diff)
			}

			legacyQuery := GenerateCloudSQLAuditExampleQuery(tc.tenantProjectID)
			if diff := cmp.Diff(gotQuery, legacyQuery); diff != "" {
				t.Errorf("GenerateCloudSQLAuditExampleQuery() mismatch (-want +got):\n%s", diff)
			}

			gotMetrics := sq.GenerateMonitoringMetricFilters()
			if diff := cmp.Diff(tc.wantMetricFilters, gotMetrics); diff != "" {
				t.Errorf("GenerateMonitoringMetricFilters() mismatch (-want +got):\n%s", diff)
			}

			if sq.AllFiltersSupportMetrics() != tc.wantSupportMetricsFlag {
				t.Errorf("AllFiltersSupportMetrics() = %v, want %v", sq.AllFiltersSupportMetrics(), tc.wantSupportMetricsFlag)
			}

			if sq.Incomplete != tc.wantIncomplete {
				t.Errorf("Incomplete = %v, want %v", sq.Incomplete, tc.wantIncomplete)
			}
		})
	}
}

func TestCloudSQLAuditListLogEntriesTaskSetting(t *testing.T) {
	testCases := []struct {
		name              string
		tenantProjectID   string
		wantResourceNames []string
		wantIncomplete    bool
	}{
		{
			name:              "valid tenant project id",
			tenantProjectID:   "my-tenant-project-tp",
			wantResourceNames: []string{"projects/my-tenant-project-tp"},
			wantIncomplete:    false,
		},
		{
			name:              "empty tenant project id returns empty slice",
			tenantProjectID:   "",
			wantResourceNames: []string{},
			wantIncomplete:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = tasktest.WithTaskResult(ctx, privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref(), tc.tenantProjectID)

			setting := &cloudSQLAuditListLogEntriesTaskSetting{}

			gotResources, err := setting.DefaultResourceNames(ctx)
			if err != nil {
				t.Fatalf("DefaultResourceNames() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantResourceNames, gotResources); diff != "" {
				t.Errorf("DefaultResourceNames() mismatch (-want +got):\n%s", diff)
			}

			queries, err := setting.Queries(ctx)
			if err != nil {
				t.Fatalf("Queries() unexpected error: %v", err)
			}
			if len(queries) != 1 {
				t.Fatalf("Queries() returned %d queries, want 1", len(queries))
			}
			if queries[0].Incomplete != tc.wantIncomplete {
				t.Errorf("Queries()[0].Incomplete = %v, want %v", queries[0].Incomplete, tc.wantIncomplete)
			}

			if setting.QueryName() != "Cloud Composer Tenant Project Cloud SQL Audit Logs" {
				t.Errorf("QueryName() = %q, want %q", setting.QueryName(), "Cloud Composer Tenant Project Cloud SQL Audit Logs")
			}
		})
	}
}

func TestCloudSQLAuditLogsQueryTask_DryRun(t *testing.T) {
	t.Parallel()
	startTime := time.Date(2025, time.January, 1, 1, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, time.January, 1, 1, 1, 0, 0, time.UTC)

	clientFactory, err := googlecloud.NewClientFactory()
	if err != nil {
		t.Fatalf("failed to create clientFactory: %v", err)
	}

	testCases := []struct {
		name            string
		tenantProjectID string
		wantIncomplete  bool
		wantQuery       string
	}{
		{
			name:            "valid tenant project id",
			tenantProjectID: "my-tenant-project-tp",
			wantIncomplete:  false,
			wantQuery: `resource.type="cloudsql_database"
resource.labels.project_id="my-tenant-project-tp"
LOG_ID("cloudaudit.googleapis.com/activity")
timestamp >= "2025-01-01T01:00:00+0000"
timestamp <= "2025-01-01T01:01:00+0000"`,
		},
		{
			name:            "empty tenant project id is incomplete",
			tenantProjectID: "",
			wantIncomplete:  true,
			wantQuery: `resource.type="cloudsql_database"
resource.labels.project_id=""
LOG_ID("cloudaudit.googleapis.com/activity")
timestamp >= "2025-01-01T01:00:00+0000"
timestamp <= "2025-01-01T01:01:00+0000"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceNamesInput := googlecloudcommon_contract.NewResourceNamesInput()
			if tc.tenantProjectID != "" {
				resourceNamesInput.UpdateDefaultResourceNamesForQuery(privatecomposer_contract.CloudSQLAuditLogsQueryTaskID.ReferenceIDString(), []string{"projects/" + tc.tenantProjectID})
			}

			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())
			gotLogs, _, err := inspectiontest.RunInspectionTask(ctx, CloudSQLAuditLogsQueryTask, inspectioncore_contract.TaskModeDryRun, map[string]any{},
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputStartTimeTaskID.Ref(), startTime),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputEndTimeTaskID.Ref(), endTime),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.APIClientFactoryTaskID.Ref(), clientFactory),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputLoggingFilterResourceNameTaskID.Ref(), resourceNamesInput),
				tasktest.NewTaskDependencyValuePair(privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref(), tc.tenantProjectID),
			)
			if err != nil {
				t.Fatalf("DryRun returned unexpected error: %v", err)
			}
			if len(gotLogs) != 0 {
				t.Errorf("DryRun should return 0 logs, got %d", len(gotLogs))
			}

			metadata := khictx.MustGetValue(ctx, inspectioncore_contract.InspectionRunMetadata)
			queryMetadata, found := typedmap.Get(metadata, inspectionmetadata.QueryMetadataKey)
			if !found {
				t.Fatalf("QueryMetadata not found in metadata")
			}

			serialized := queryMetadata.ToSerializable().([]*inspectionmetadata.QueryItem)
			if len(serialized) != 1 {
				t.Fatalf("expected 1 QueryItem, got %d", len(serialized))
			}

			if serialized[0].Incomplete != tc.wantIncomplete {
				t.Errorf("Incomplete mismatch: got %v, want %v", serialized[0].Incomplete, tc.wantIncomplete)
			}

			if diff := cmp.Diff(tc.wantQuery, serialized[0].Query); diff != "" {
				t.Errorf("Query mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
