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
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/gcpqueryutil"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
	"github.com/google/go-cmp/cmp"
)

func TestGeneratePrivateGKEMasterStructuredQuery(t *testing.T) {
	testCases := []struct {
		name                   string
		projectID              string
		componentFilter        *gcpqueryutil.SetFilterParseResult
		wantQuery              string
		wantMetricFilters      []string
		wantSupportMetricsFlag bool
		wantIncomplete         bool
	}{
		{
			name:      "subtractive filter excluding kube-apiserver",
			projectID: "tp-12345",
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				SubtractMode: true,
				Subtractives: []string{"kube-apiserver"},
			},
			wantQuery: `resource.type=("container" OR "gce_instance")
-LOG_ID("cloudaudit.googleapis.com/activity")
-LOG_ID("cloudaudit.googleapis.com/data_access")
-LOG_ID("compute.googleapis.com/shielded_vm_integrity")
-LOG_ID("serialconsole.googleapis.com/serial_port_1_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_2_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_3_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_debug_output")
-LOG_ID("kube-apiserver")
resource.labels.project_id="tp-12345"`,
			wantMetricFilters:      nil,
			wantSupportMetricsFlag: false,
			wantIncomplete:         false,
		},
		{
			name:      "subtractive filter with no exclusions",
			projectID: "tp-12345",
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				SubtractMode: true,
				Subtractives: []string{},
			},
			wantQuery: `resource.type=("container" OR "gce_instance")
-LOG_ID("cloudaudit.googleapis.com/activity")
-LOG_ID("cloudaudit.googleapis.com/data_access")
-LOG_ID("compute.googleapis.com/shielded_vm_integrity")
-LOG_ID("serialconsole.googleapis.com/serial_port_1_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_2_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_3_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_debug_output")
-- no master component filter
resource.labels.project_id="tp-12345"`,
			wantMetricFilters:      nil,
			wantSupportMetricsFlag: false,
			wantIncomplete:         false,
		},
		{
			name:      "additive filter with single component",
			projectID: "tp-12345",
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				SubtractMode: false,
				Additives:    []string{"kube-scheduler"},
			},
			wantQuery: `resource.type=("container" OR "gce_instance")
-LOG_ID("cloudaudit.googleapis.com/activity")
-LOG_ID("cloudaudit.googleapis.com/data_access")
-LOG_ID("compute.googleapis.com/shielded_vm_integrity")
-LOG_ID("serialconsole.googleapis.com/serial_port_1_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_2_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_3_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_debug_output")
LOG_ID("kube-scheduler")
resource.labels.project_id="tp-12345"`,
			wantMetricFilters:      nil,
			wantSupportMetricsFlag: false,
			wantIncomplete:         false,
		},
		{
			name:      "filter with validation error",
			projectID: "tp-12345",
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				ValidationError: "invalid characters",
			},
			wantQuery: `resource.type=("container" OR "gce_instance")
-LOG_ID("cloudaudit.googleapis.com/activity")
-LOG_ID("cloudaudit.googleapis.com/data_access")
-LOG_ID("compute.googleapis.com/shielded_vm_integrity")
-LOG_ID("serialconsole.googleapis.com/serial_port_1_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_2_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_3_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_debug_output")
-- Failed to generate component name filter due to the validation error "invalid characters"
resource.labels.project_id="tp-12345"`,
			wantMetricFilters:      nil,
			wantSupportMetricsFlag: false,
			wantIncomplete:         false,
		},
		{
			name:            "empty project ID is incomplete",
			projectID:       "",
			componentFilter: nil,
			wantQuery: `resource.type=("container" OR "gce_instance")
-LOG_ID("cloudaudit.googleapis.com/activity")
-LOG_ID("cloudaudit.googleapis.com/data_access")
-LOG_ID("compute.googleapis.com/shielded_vm_integrity")
-LOG_ID("serialconsole.googleapis.com/serial_port_1_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_2_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_3_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_debug_output")
resource.labels.project_id=""`,
			wantMetricFilters:      nil,
			wantSupportMetricsFlag: false,
			wantIncomplete:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sq := GeneratePrivateGKEMasterStructuredQuery(tc.projectID, tc.componentFilter)
			gotQuery := sq.GenerateCloudLoggingQuery()
			if diff := cmp.Diff(tc.wantQuery, gotQuery); diff != "" {
				t.Errorf("GenerateCloudLoggingQuery() mismatch (-want +got):\n%s", diff)
			}

			legacyQuery := GeneratePrivateGKEMasterQuery(tc.projectID, tc.componentFilter)
			if diff := cmp.Diff(gotQuery, legacyQuery); diff != "" {
				t.Errorf("GeneratePrivateGKEMasterQuery() mismatch (-want +got):\n%s", diff)
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

func TestListLogEntriesTaskSetting(t *testing.T) {
	testCases := []struct {
		name              string
		logSource         *privategkemaster_contract.LogSource
		componentFilter   *gcpqueryutil.SetFilterParseResult
		wantResourceNames []string
		wantIncomplete    bool
	}{
		{
			name: "logSource provided",
			logSource: &privategkemaster_contract.LogSource{
				TenantProjectID:     "tp-gke-12345",
				LogViewResourceName: "projects/tp-gke-12345/locations/global/buckets/_Default/views/_AllLogs",
			},
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				SubtractMode: true,
				Subtractives: []string{"kube-apiserver"},
			},
			wantResourceNames: []string{
				"projects/tp-gke-12345/locations/global/buckets/_Default/views/_AllLogs",
			},
			wantIncomplete: false,
		},
		{
			name:      "nil logSource returns empty slice and incomplete query",
			logSource: nil,
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				SubtractMode: true,
			},
			wantResourceNames: []string{},
			wantIncomplete:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			ctx = tasktest.WithTaskResult(ctx, privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref(), tc.logSource)
			ctx = tasktest.WithTaskResult(ctx, privategkemaster_contract.InputPrivateGKEMasterComponentNameFilterTaskID.Ref(), tc.componentFilter)

			setting := &listLogEntriesTaskSetting{}

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

			if setting.QueryName() != "Private GKE master logs" {
				t.Errorf("QueryName() = %q, want %q", setting.QueryName(), "Private GKE master logs")
			}

			timePartitionCount, err := setting.TimePartitionCount(ctx)
			if err != nil {
				t.Fatalf("TimePartitionCount() unexpected error: %v", err)
			}
			if timePartitionCount != 10 {
				t.Errorf("TimePartitionCount() = %d, want 10", timePartitionCount)
			}
		})
	}
}

func TestListLogEntriesTask_DryRun(t *testing.T) {
	t.Parallel()
	startTime := time.Date(2025, time.January, 1, 1, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, time.January, 1, 1, 1, 0, 0, time.UTC)

	clientFactory, err := googlecloud.NewClientFactory()
	if err != nil {
		t.Fatalf("failed to create clientFactory: %v", err)
	}

	testCases := []struct {
		name            string
		logSource       *privategkemaster_contract.LogSource
		componentFilter *gcpqueryutil.SetFilterParseResult
		wantIncomplete  bool
		wantQuery       string
	}{
		{
			name: "logSource provided",
			logSource: &privategkemaster_contract.LogSource{
				TenantProjectID:     "tp-12345",
				LogViewResourceName: "projects/tp-12345/locations/global/buckets/_Default/views/_AllLogs",
			},
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				SubtractMode: true,
				Subtractives: []string{"kube-apiserver"},
			},
			wantIncomplete: false,
			wantQuery: `resource.type=("container" OR "gce_instance")
-LOG_ID("cloudaudit.googleapis.com/activity")
-LOG_ID("cloudaudit.googleapis.com/data_access")
-LOG_ID("compute.googleapis.com/shielded_vm_integrity")
-LOG_ID("serialconsole.googleapis.com/serial_port_1_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_2_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_3_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_debug_output")
-LOG_ID("kube-apiserver")
resource.labels.project_id="tp-12345"
timestamp >= "2025-01-01T01:00:00+0000"
timestamp <= "2025-01-01T01:01:00+0000"`,
		},
		{
			name:      "nil logSource is incomplete",
			logSource: nil,
			componentFilter: &gcpqueryutil.SetFilterParseResult{
				SubtractMode: true,
			},
			wantIncomplete: true,
			wantQuery: `resource.type=("container" OR "gce_instance")
-LOG_ID("cloudaudit.googleapis.com/activity")
-LOG_ID("cloudaudit.googleapis.com/data_access")
-LOG_ID("compute.googleapis.com/shielded_vm_integrity")
-LOG_ID("serialconsole.googleapis.com/serial_port_1_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_2_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_3_output")
-LOG_ID("serialconsole.googleapis.com/serial_port_debug_output")
-- no master component filter
resource.labels.project_id=""
timestamp >= "2025-01-01T01:00:00+0000"
timestamp <= "2025-01-01T01:01:00+0000"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceNamesInput := googlecloudcommon_contract.NewResourceNamesInput()
			if tc.logSource != nil {
				resourceNamesInput.UpdateDefaultResourceNamesForQuery(privategkemaster_contract.ListLogEntriesTaskID.ReferenceIDString(), []string{tc.logSource.LogViewResourceName})
			}

			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())
			gotLogs, _, err := inspectiontest.RunInspectionTask(ctx, listLogEntriesTask, inspectioncore_contract.TaskModeDryRun, map[string]any{},
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputStartTimeTaskID.Ref(), startTime),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputEndTimeTaskID.Ref(), endTime),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.APIClientFactoryTaskID.Ref(), clientFactory),
				tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputLoggingFilterResourceNameTaskID.Ref(), resourceNamesInput),
				tasktest.NewTaskDependencyValuePair(privategkemaster_contract.InputGKEMasterLogSourceTaskID.Ref(), tc.logSource),
				tasktest.NewTaskDependencyValuePair(privategkemaster_contract.InputPrivateGKEMasterComponentNameFilterTaskID.Ref(), tc.componentFilter),
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
