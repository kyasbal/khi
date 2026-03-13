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

package privatecomposerv3_impl

import (
	"context"
	"testing"

	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposerv3_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposerv3/contract"
)

func TestClusterIdentityTask(t *testing.T) {
	// Mock required dependencies
	mockTenantIDTask := tasktest.StubTaskFromReferenceID(privatecomposerv3_contract.InputComposerV3TenantProjectIdTaskID.Ref(), "my-tenant-project-tp", nil)
	mockClusterNameTask := tasktest.StubTaskFromReferenceID(googlecloudk8scommon_contract.InputClusterNameTaskID.Ref(), "my-cluster", nil)
	mockLocationsTask := tasktest.StubTaskFromReferenceID(googlecloudcommon_contract.InputLocationsTaskID.Ref(), "us-central1-c", nil)
	mockPrefixTask := tasktest.StubTaskFromReferenceID(googlecloudk8scommon_contract.ClusterNamePrefixTaskRef, "my-cluster-prefix-", nil)

	result, _, err := inspectiontest.RunInspectionTaskWithDependency(
		inspectiontest.WithDefaultTestInspectionTaskContext(context.Background()),
		ClusterIdentityTask,
		[]coretask.UntypedTask{mockTenantIDTask, mockClusterNameTask, mockLocationsTask, mockPrefixTask},
		inspectioncore_contract.TaskModeRun,
		map[string]any{},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ProjectID != "my-tenant-project-tp" {
		t.Errorf("expected actual project ID %q, but got %q", "my-tenant-project-tp", result.ProjectID)
	}
	if result.ClusterName != "my-cluster" {
		t.Errorf("expected actual cluster name %q, but got %q", "my-cluster", result.ClusterName)
	}
}

func TestComposerClusterIdentityTask(t *testing.T) {
	// Mock required dependencies
	mockProjectIDTask := tasktest.StubTaskFromReferenceID(googlecloudcommon_contract.InputProjectIdTaskID.Ref(), "my-customer-project", nil)
	mockLocationsTask := tasktest.StubTaskFromReferenceID(googlecloudcommon_contract.InputLocationsTaskID.Ref(), "us-central1-c", nil)

	result, _, err := inspectiontest.RunInspectionTaskWithDependency(
		inspectiontest.WithDefaultTestInspectionTaskContext(context.Background()),
		ComposerClusterIdentityTask,
		[]coretask.UntypedTask{mockProjectIDTask, mockLocationsTask},
		inspectioncore_contract.TaskModeRun,
		map[string]any{},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ProjectID != "my-customer-project" {
		t.Errorf("expected actual project ID %q, but got %q", "my-customer-project", result.ProjectID)
	}
	if result.ClusterName != "" {
		t.Errorf("expected actual cluster name %q, but got %q", "", result.ClusterName)
	}
	if result.ClusterTypePrefix != "" {
		t.Errorf("expected actual cluster type prefix %q, but got %q", "", result.ClusterTypePrefix)
	}
}
