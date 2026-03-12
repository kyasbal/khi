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
	"testing"

	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
	"github.com/google/go-cmp/cmp"
)

func TestInputGKEMasterLogSourceTask_ValidatorAndConverter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
		wantValue *privategkemaster_contract.LogSource
	}{
		{
			name:      "valid input with storageScope and project",
			input:     "https://pantheon.corp.google.com/logs/query;query=resource.labels.project_id=%22my-tenant-project%22;storageScope=storage%2Cprojects%2Fgke-prod-bq-logs;Foo=Bar?project=gke-prod-bq-logs",
			wantValid: true,
			wantValue: &privategkemaster_contract.LogSource{
				TenantProjectID:     "my-tenant-project",
				LogViewResourceName: "projects/gke-prod-bq-logs",
			},
		},
		{
			name:      "complex real input",
			input:     `https://pantheon.corp.google.com/logs/query;query=resource.type=%22container%22%0Aresource.labels.container_name=%22cloud-controller-manager%22%0Aresource.labels.namespace_id=%22kube-system%22%0Aresource.labels.pod_id:%22cloud-controller-manager%22%0Aresource.labels.cluster_name=%22same-name%22%0Aresource.labels.project_id=%22foobarqux-tp%22;storageScope=storage%2Cprojects%2Fgke-prod-bq-logs%2Flocations%2Fasia-northeast1%2Fbuckets%2Fprod-kcp-logging-asia-northeast1-shard_7_of_16%2Fviews%2F_AllLogs?project=gke-prod-bq-logs`,
			wantValid: true,
			wantValue: &privategkemaster_contract.LogSource{
				TenantProjectID:     "foobarqux-tp",
				LogViewResourceName: "projects/gke-prod-bq-logs/locations/asia-northeast1/buckets/prod-kcp-logging-asia-northeast1-shard_7_of_16/views/_AllLogs",
			},
		},
		{
			name:      "invalid input without project",
			input:     "https://pantheon.corp.google.com/logs/query;query=...;storageScope=storage%2Cprojects%2Fgke-prod-bq-logs",
			wantValid: false,
			wantValue: nil,
		},
		{
			name:      "invalid input without storageScope",
			input:     "https://pantheon.corp.google.com/logs/query;query=...?project=my-project",
			wantValid: false,
			wantValue: nil,
		},
		{
			name:      "new url format with logName",
			input:     "https://pantheon.corp.google.com/logs/query;query=%28resource.type%3D%22gce_instance%22%20OR%20resource.type%3D%22container%22%29%20AND%20logName%3D%22projects%2Fabcdefg-tp%2Flogs%2Fkube-master-installation%22;storageScope=storage,projects%2Fgke-prod-bq-logs%2Flocations%2Fus%2Fbuckets%2Fprod-kcp-logging-us-central1-a-shard_4_of_16%2Fviews%2Fcluster-view-ikakeru-b290;duration=P2D?project=gke-prod-bq-logs",
			wantValid: true,
			wantValue: &privategkemaster_contract.LogSource{
				TenantProjectID:     "abcdefg-tp",
				LogViewResourceName: "projects/gke-prod-bq-logs/locations/us/buckets/prod-kcp-logging-us-central1-a-shard_4_of_16/views/cluster-view-ikakeru-b290",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(context.Background())
			ctx = tasktest.WithTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(), googlecloudk8scommon_contract.GoogleCloudClusterIdentity{ProjectID: "test-project"})
			ctx = tasktest.WithTaskResult(ctx, privatecommon_contract.JustificationFormTaskID.Ref(), "b/123456")

			// Test Validator
			msg, err := validateMasterLogLink(ctx, tt.input)
			if err != nil {
				t.Fatalf("unexpected error from validator: %v", err)
			}
			isValid := msg == ""
			if isValid != tt.wantValid {
				t.Errorf("Validator(%q) valid = %v, want %v (msg: %q)", tt.input, isValid, tt.wantValid, msg)
			}

			if !tt.wantValid {
				return
			}

			// Test Converter
			got, err := convertMasterLogLink(ctx, tt.input)
			if err != nil {
				t.Fatalf("unexpected error from converter: %v", err)
			}
			if diff := cmp.Diff(tt.wantValue, got); diff != "" {
				t.Errorf("Converter(%q) mismatch (-want +got):\n%s", tt.input, diff)
			}
		})
	}
}
