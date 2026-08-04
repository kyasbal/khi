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

package privatecomposer_impl

import (
	"context"
	"testing"

	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
	"github.com/google/go-cmp/cmp"
)

func TestClusterIdentityTask(t *testing.T) {
	testCases := []struct {
		name         string
		tenantID     string
		clusterName  string
		location     string
		prefixPolicy googlecloudk8scommon_contract.ClusterPrefixPolicy
		want         googlecloudk8scommon_contract.GoogleCloudClusterIdentity
	}{
		{
			name:         "resolves cluster identity with tenant project id",
			tenantID:     "my-tenant-project-tp",
			clusterName:  "my-cluster",
			location:     "us-central1-c",
			prefixPolicy: googlecloudk8scommon_contract.ClusterPrefixPolicy{},
			want: googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
				ProjectID:    "my-tenant-project-tp",
				ClusterName:  "my-cluster",
				Location:     "us-central1-c",
				PrefixPolicy: googlecloudk8scommon_contract.ClusterPrefixPolicy{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTenantIDTask := tasktest.StubTaskFromReferenceID(privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref(), tc.tenantID, nil)
			mockClusterNameTask := tasktest.StubTaskFromReferenceID(googlecloudk8scommon_contract.InputClusterNameTaskID.Ref(), tc.clusterName, nil)
			mockLocationsTask := tasktest.StubTaskFromReferenceID(googlecloudcommon_contract.InputLocationsTaskID.Ref(), tc.location, nil)
			mockPrefixTask := tasktest.StubTaskFromReferenceID(googlecloudk8scommon_contract.ClusterNamePrefixTaskRef, tc.prefixPolicy, nil)

			result, _, err := inspectiontest.RunInspectionTaskWithDependency(
				inspectiontest.WithDefaultTestInspectionTaskContext(context.Background()),
				ClusterIdentityTask,
				[]coretask.UntypedTask{mockTenantIDTask, mockClusterNameTask, mockLocationsTask, mockPrefixTask},
				inspectioncore_contract.TaskModeRun,
				map[string]any{},
			)
			if err != nil {
				t.Fatalf("unexpected error running task: %v", err)
			}

			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("ClusterIdentityTask mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestComposerClusterIdentityTask(t *testing.T) {
	testCases := []struct {
		name      string
		projectID string
		location  string
		want      googlecloudk8scommon_contract.GoogleCloudClusterIdentity
	}{
		{
			name:      "resolves composer cluster identity with customer project id",
			projectID: "my-customer-project",
			location:  "us-central1-c",
			want: googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
				ProjectID: "my-customer-project",
				Location:  "us-central1-c",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProjectIDTask := tasktest.StubTaskFromReferenceID(googlecloudcommon_contract.InputProjectIdTaskID.Ref(), tc.projectID, nil)
			mockLocationsTask := tasktest.StubTaskFromReferenceID(googlecloudcommon_contract.InputLocationsTaskID.Ref(), tc.location, nil)

			result, _, err := inspectiontest.RunInspectionTaskWithDependency(
				inspectiontest.WithDefaultTestInspectionTaskContext(context.Background()),
				ComposerClusterIdentityTask,
				[]coretask.UntypedTask{mockProjectIDTask, mockLocationsTask},
				inspectioncore_contract.TaskModeRun,
				map[string]any{},
			)
			if err != nil {
				t.Fatalf("unexpected error running task: %v", err)
			}

			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("ComposerClusterIdentityTask mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
