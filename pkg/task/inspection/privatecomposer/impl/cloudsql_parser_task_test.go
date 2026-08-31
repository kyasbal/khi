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

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	core_contract "github.com/GoogleCloudPlatform/khi/pkg/task/core/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testchangeset"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testlog"
)

func TestCloudSQLLogsIngester_ProcessLog(t *testing.T) {
	testCases := []struct {
		name   string
		input  *log.Log
		assert func(t *testing.T, cs *khifilev6.LogChangeSet)
	}{
		{
			name: "sample log with LOG summary extraction",
			input: testlog.NewMockLog(
				time.Date(2026, 8, 4, 14, 52, 17, 0, time.UTC),
				inspectioncore_contract.DefaultSeverityFieldSet{
					Severity: inspectioncore_contract.SeverityInfo,
				},
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
					Summary:    "connection received: host=127.0.0.1 port=60716",
				},
				googlecloudcommon_contract.GCPMainMessageFieldSet{
					MainMessage: "2026-08-04 14:52:17.232 UTC [115855]: [1-1] db=[unknown],user=[unknown] LOG:  connection received: host=127.0.0.1 port=60716",
				},
			),
			assert: func(t *testing.T, cs *khifilev6.LogChangeSet) {
				testchangeset.AssertLog(t, cs).
					HasSummary("connection received: host=127.0.0.1 port=60716").
					HasSeverity(inspectioncore_contract.SeverityInfo)
			},
		},
		{
			name: "fallback to main message if summary is empty",
			input: testlog.NewMockLog(
				time.Date(2026, 8, 4, 14, 52, 17, 0, time.UTC),
				inspectioncore_contract.DefaultSeverityFieldSet{
					Severity: inspectioncore_contract.SeverityError,
				},
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
					Summary:    "",
				},
				googlecloudcommon_contract.GCPMainMessageFieldSet{
					MainMessage: "unhandled error occurred",
				},
			),
			assert: func(t *testing.T, cs *khifilev6.LogChangeSet) {
				testchangeset.AssertLog(t, cs).
					HasSummary("unhandled error occurred").
					HasSeverity(inspectioncore_contract.SeverityError)
			},
		},
	}

	ingester := &cloudSQLLogsIngester{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cs, err := ingester.ProcessLog(ctx, tc.input)
			if err != nil {
				t.Fatalf("ProcessLog() returned unexpected error: %v", err)
			}
			tc.assert(t, cs)
		})
	}
}

func TestCloudSQLLogsTimelineMapper_ProcessLogByGroup(t *testing.T) {
	builder := khifilev6.NewBuilder()
	expectedPath := privatecomposer_contract.MustCloudSQLLogTimeline(
		khictx.WithValue(context.Background(), inspectioncore_contract.Builder, builder),
		"my-tenant-project-tp",
		"us-central1-composer-2-170061d7-sql",
		"postgres.log",
	)

	testCases := []struct {
		name      string
		inputLog  *log.Log
		projectID string
		assert    func(t *testing.T, cs *khifilev6.TimelineChangeSet)
	}{
		{
			name: "adds event to cloudsql log timeline",
			inputLog: testlog.NewMockLog(
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID:  "us-central1-composer-2-170061d7-sql",
					LogFileName: "postgres.log",
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(expectedPath)
			},
		},
	}

	mapper := &cloudSQLLogsTimelineMapper{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := khictx.WithValue(context.Background(), inspectioncore_contract.Builder, builder)
			taskDependentValues := typedmap.NewTypedMap()
			typedmap.Set(taskDependentValues, typedmap.NewTypedKey[string](privatecomposer_contract.InputComposerTenantProjectIdTaskID.ReferenceIDString()), tc.projectID)
			ctx = khictx.WithValue(ctx, core_contract.TaskResultMapContextKey, taskDependentValues)

			cs, _, err := mapper.ProcessLogByGroup(ctx, tc.inputLog, struct{}{})
			if err != nil {
				t.Fatalf("ProcessLogByGroup() returned unexpected error: %v", err)
			}
			tc.assert(t, cs)
		})
	}
}
