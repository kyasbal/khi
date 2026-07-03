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

	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
	"github.com/google/go-cmp/cmp"
)

func TestCSMCPLogQueryTaskSetting_LogFilters(t *testing.T) {
	testCases := []struct {
		name        string
		projectID   string
		serviceName string
		want        []string
	}{
		{
			name:        "all inputs provided",
			projectID:   "test-project",
			serviceName: "test-service",
			want: []string{
				`resource.type="cloud_run_revision"
resource.labels.project_id="test-project"
resource.labels.service_name="test-service"
LOG_ID("run.googleapis.com/stdout") OR LOG_ID("run.googleapis.com/stderr")`,
			},
		},
		{
			name:        "missing project id",
			projectID:   "",
			serviceName: "test-service",
			want:        []string{},
		},
		{
			name:        "missing service name",
			projectID:   "test-project",
			serviceName: "",
			want:        []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			ctx = tasktest.WithTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(), tc.projectID)
			ctx = tasktest.WithTaskResult(ctx, privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref(), tc.serviceName)

			setting := &csmcpLogQueryTaskSetting{}
			got, err := setting.LogFilters(ctx, inspectioncore_contract.TaskModeRun)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("LogFilters() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCSMCPLogQueryTaskSetting_DefaultResourceNames(t *testing.T) {
	testCases := []struct {
		name      string
		projectID string
		want      []string
	}{
		{
			name:      "project ID provided",
			projectID: "test-project",
			want:      []string{"projects/test-project"},
		},
		{
			name:      "missing project ID",
			projectID: "",
			want:      []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			ctx = tasktest.WithTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(), tc.projectID)

			setting := &csmcpLogQueryTaskSetting{}
			got, err := setting.DefaultResourceNames(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("DefaultResourceNames() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
