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

package privatecomposer_contract

import (
	"testing"

	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	"github.com/google/go-cmp/cmp"
)

func TestMustCloudSQLInstanceTimeline(t *testing.T) {
	testCases := []struct {
		name           string
		projectID      string
		instanceID     string
		wantName       string
		wantParentName string
	}{
		{
			name:           "valid project and instance id",
			projectID:      "my-tenant-project-tp",
			instanceID:     "us-central1-instance-sql",
			wantName:       "us-central1-instance-sql",
			wantParentName: "my-tenant-project-tp",
		},
		{
			name:           "empty project and instance id default to unknown",
			projectID:      "",
			instanceID:     "",
			wantName:       "unknown",
			wantParentName: "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())

			path := MustCloudSQLInstanceTimeline(ctx, tc.projectID, tc.instanceID)

			if diff := cmp.Diff(tc.wantName, path.Name.Resolve()); diff != "" {
				t.Errorf("MustCloudSQLInstanceTimeline() name mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantParentName, path.Parent.Name.Resolve()); diff != "" {
				t.Errorf("MustCloudSQLInstanceTimeline() parent name mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(TimelineTypeCloudSQLInstance.GetId(), path.Type.GetId()); diff != "" {
				t.Errorf("MustCloudSQLInstanceTimeline() type mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMustCloudSQLLogTimeline(t *testing.T) {
	testCases := []struct {
		name           string
		projectID      string
		instanceID     string
		logFileName    string
		wantName       string
		wantParentName string
	}{
		{
			name:           "valid project, instance id, and log file name",
			projectID:      "my-tenant-project-tp",
			instanceID:     "us-central1-instance-sql",
			logFileName:    "postgres.log",
			wantName:       "postgres.log",
			wantParentName: "us-central1-instance-sql",
		},
		{
			name:           "empty log file name defaults to unknown",
			projectID:      "my-tenant-project-tp",
			instanceID:     "us-central1-instance-sql",
			logFileName:    "",
			wantName:       "unknown",
			wantParentName: "us-central1-instance-sql",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())

			path := MustCloudSQLLogTimeline(ctx, tc.projectID, tc.instanceID, tc.logFileName)

			if diff := cmp.Diff(tc.wantName, path.Name.Resolve()); diff != "" {
				t.Errorf("MustCloudSQLLogTimeline() name mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantParentName, path.Parent.Name.Resolve()); diff != "" {
				t.Errorf("MustCloudSQLLogTimeline() parent name mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(TimelineTypeCloudSQLLog.GetId(), path.Type.GetId()); diff != "" {
				t.Errorf("MustCloudSQLLogTimeline() type mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
