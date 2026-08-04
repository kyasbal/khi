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

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	core_contract "github.com/GoogleCloudPlatform/khi/pkg/task/core/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
	"github.com/google/go-cmp/cmp"
)

func TestCloudSQLListLogEntriesTaskSetting(t *testing.T) {
	testCases := []struct {
		name              string
		tenantProjectID   string
		wantLogFilters    []string
		wantResourceNames []string
	}{
		{
			name:            "valid tenant project id",
			tenantProjectID: "my-tenant-project-tp",
			wantLogFilters: []string{
				`resource.type="cloudsql_database"
resource.labels.project_id="my-tenant-project-tp"

LOG_ID("cloudsql.googleapis.com/postgres.log") OR LOG_ID("cloudsql.googleapis.com/postgres-audit.log") OR LOG_ID("cloudsql.googleapis.com/postgres-upgrade.log")`,
			},
			wantResourceNames: []string{"projects/my-tenant-project-tp"},
		},
		{
			name:              "empty tenant project id returns empty slice",
			tenantProjectID:   "",
			wantLogFilters:    []string{},
			wantResourceNames: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			taskDependentValues := typedmap.NewTypedMap()
			typedmap.Set(taskDependentValues, typedmap.NewTypedKey[string](privatecomposer_contract.InputComposerTenantProjectIdTaskID.ReferenceIDString()), tc.tenantProjectID)
			ctx = khictx.WithValue(ctx, core_contract.TaskResultMapContextKey, taskDependentValues)

			setting := &cloudSQLListLogEntriesTaskSetting{}
			taskMode := inspectioncore_contract.TaskModeDryRun

			gotFilters, err := setting.LogFilters(ctx, taskMode)
			if err != nil {
				t.Fatalf("LogFilters() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantLogFilters, gotFilters); diff != "" {
				t.Errorf("LogFilters() mismatch (-want +got):\n%s", diff)
			}

			gotResources, err := setting.DefaultResourceNames(ctx)
			if err != nil {
				t.Fatalf("DefaultResourceNames() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantResourceNames, gotResources); diff != "" {
				t.Errorf("DefaultResourceNames() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
