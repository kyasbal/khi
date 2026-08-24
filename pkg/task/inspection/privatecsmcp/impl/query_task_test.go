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

package privatecsmcp_impl

import (
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
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateCSMCPStructuredQuery(t *testing.T) {
	testCases := []struct {
		name                   string
		projectID              string
		serviceName            string
		wantQuery              string
		wantMetricFilters      []string
		wantSupportMetricsFlag bool
		wantIncomplete         bool
	}{
		{
			name:        "all inputs provided",
			projectID:   "test-project",
			serviceName: "test-service",
			wantQuery: `resource.type="cloud_run_revision"
resource.labels.project_id="test-project"
resource.labels.service_name="test-service"
(LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr"))`,
			wantMetricFilters: []string{
				`metric.type = "logging.googleapis.com/log_entry_count" AND resource.type = "cloud_run_revision" AND resource.labels.project_id = "test-project" AND resource.labels.service_name = "test-service" AND metric.labels.log = one_of("run.googleapis.com/stdout", "run.googleapis.com/stderr")`,
			},
			wantSupportMetricsFlag: true,
			wantIncomplete:         false,
		},
		{
			name:        "missing project id",
			projectID:   "",
			serviceName: "test-service",
			wantQuery: `resource.type="cloud_run_revision"
resource.labels.project_id=""
resource.labels.service_name="test-service"
(LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr"))`,
			wantMetricFilters: []string{
				`metric.type = "logging.googleapis.com/log_entry_count" AND resource.type = "cloud_run_revision" AND resource.labels.project_id = "" AND resource.labels.service_name = "test-service" AND metric.labels.log = one_of("run.googleapis.com/stdout", "run.googleapis.com/stderr")`,
			},
			wantSupportMetricsFlag: true,
			wantIncomplete:         true,
		},
		{
			name:        "missing service name",
			projectID:   "test-project",
			serviceName: "",
			wantQuery: `resource.type="cloud_run_revision"
resource.labels.project_id="test-project"
resource.labels.service_name=""
(LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr"))`,
			wantMetricFilters: []string{
				`metric.type = "logging.googleapis.com/log_entry_count" AND resource.type = "cloud_run_revision" AND resource.labels.project_id = "test-project" AND resource.labels.service_name = "" AND metric.labels.log = one_of("run.googleapis.com/stdout", "run.googleapis.com/stderr")`,
			},
			wantSupportMetricsFlag: true,
			wantIncomplete:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sq := GenerateCSMCPStructuredQuery(tc.projectID, tc.serviceName)
			gotQuery := sq.GenerateCloudLoggingQuery()
			if diff := cmp.Diff(tc.wantQuery, gotQuery); diff != "" {
				t.Errorf("GenerateCloudLoggingQuery() mismatch (-want +got):\n%s", diff)
			}

			legacyQuery := GenerateCSMCPQuery(tc.projectID, tc.serviceName)
			if diff := cmp.Diff(gotQuery, legacyQuery); diff != "" {
				t.Errorf("GenerateCSMCPQuery() mismatch (-want +got):\n%s", diff)
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

func TestCSMCPLogQueryTaskSetting(t *testing.T) {
	testCases := []struct {
		name              string
		projectID         string
		serviceName       string
		wantResourceNames []string
		wantIncomplete    bool
	}{
		{
			name:              "all inputs provided",
			projectID:         "test-project",
			serviceName:       "test-service",
			wantResourceNames: []string{"projects/test-project"},
			wantIncomplete:    false,
		},
		{
			name:              "missing project ID",
			projectID:         "",
			serviceName:       "test-service",
			wantResourceNames: []string{},
			wantIncomplete:    true,
		},
		{
			name:              "missing service name",
			projectID:         "test-project",
			serviceName:       "",
			wantResourceNames: []string{"projects/test-project"},
			wantIncomplete:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			ctx = tasktest.WithTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(), tc.projectID)
			ctx = tasktest.WithTaskResult(ctx, privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref(), tc.serviceName)

			setting := &csmcpLogQueryTaskSetting{}

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

			if setting.QueryName() != "CSM CP logs" {
				t.Errorf("QueryName() = %q, want %q", setting.QueryName(), "CSM CP logs")
			}
		})
	}
}

func TestLogQueryTask_DryRun(t *testing.T) {
	t.Parallel()
	startTime := time.Date(2025, time.January, 1, 1, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, time.January, 1, 1, 1, 0, 0, time.UTC)

	clientFactory, err := googlecloud.NewClientFactory()
	if err != nil {
		t.Fatalf("failed to create clientFactory: %v", err)
	}

	testCases := []struct {
		name           string
		projectID      string
		serviceName    string
		wantIncomplete bool
		wantQuery      string
	}{
		{
			name:           "all inputs provided",
			projectID:      "test-project",
			serviceName:    "test-service",
			wantIncomplete: false,
			wantQuery: `resource.type="cloud_run_revision"
resource.labels.project_id="test-project"
resource.labels.service_name="test-service"
(LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr"))
timestamp >= "2025-01-01T01:00:00+0000"
timestamp <= "2025-01-01T01:01:00+0000"`,
		},
		{
			name:           "missing service name is incomplete",
			projectID:      "test-project",
			serviceName:    "",
			wantIncomplete: true,
			wantQuery: `resource.type="cloud_run_revision"
resource.labels.project_id="test-project"
resource.labels.service_name=""
(LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr"))
timestamp >= "2025-01-01T01:00:00+0000"
timestamp <= "2025-01-01T01:01:00+0000"`,
		},
		{
			name:           "missing project id is incomplete",
			projectID:      "",
			serviceName:    "test-service",
			wantIncomplete: true,
			wantQuery: `resource.type="cloud_run_revision"
resource.labels.project_id=""
resource.labels.service_name="test-service"
(LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr"))
timestamp >= "2025-01-01T01:00:00+0000"
timestamp <= "2025-01-01T01:01:00+0000"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceNamesInput := googlecloudcommon_contract.NewResourceNamesInput()
			if tc.projectID != "" {
				resourceNamesInput.UpdateDefaultResourceNamesForQuery(privatecsmcp_contract.LogQueryTaskID.ReferenceIDString(), []string{"projects/" + tc.projectID})
			}

			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())
			gotLogs, _, err := inspectiontest.RunInspectionTask(ctx, LogQueryTask, inspectioncore_contract.TaskModeDryRun, map[string]any{},
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputStartTimeTaskID.Ref(), startTime),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputEndTimeTaskID.Ref(), endTime),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.APIClientFactoryTaskID.Ref(), clientFactory),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputLoggingFilterResourceNameTaskID.Ref(), resourceNamesInput),
				tasktest.NewTaskDependencyValuePair(privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(), tc.projectID),
				tasktest.NewTaskDependencyValuePair(privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref(), tc.serviceName),
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
