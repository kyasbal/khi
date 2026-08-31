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
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testchangeset"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testlog"
)

func TestCloudSQLAuditLogsIngester_ProcessLog(t *testing.T) {
	testCases := []struct {
		name   string
		input  *log.Log
		assert func(t *testing.T, cs *khifilev6.LogChangeSet)
	}{
		{
			name: "audit log sets summary from main message",
			input: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
				inspectioncore_contract.DefaultSeverityFieldSet{
					Severity: inspectioncore_contract.SeverityInfo,
				},
				googlecloudcommon_contract.GCPMainMessageFieldSet{
					MainMessage: "google.cloud.sql.v1beta4.CloudSqlInstancesService.Create",
				},
			),
			assert: func(t *testing.T, cs *khifilev6.LogChangeSet) {
				testchangeset.AssertLog(t, cs).
					HasSummary("google.cloud.sql.v1beta4.CloudSqlInstancesService.Create").
					HasSeverity(inspectioncore_contract.SeverityInfo)
			},
		},
	}

	ingester := &cloudSQLAuditLogsIngester{}
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

func TestCloudSQLAuditLogsTimelineMapper_ProcessLogByGroup(t *testing.T) {
	builder := khifilev6.NewBuilder()
	testCtx := khictx.WithValue(context.Background(), inspectioncore_contract.Builder, builder)
	expectedInstancePath := privatecomposer_contract.MustCloudSQLInstanceTimeline(
		testCtx,
		"my-tenant-project-tp",
		"us-central1-composer-2-170061d7-sql",
	)
	expectedCreateOpPath := googlecloudcommon_contract.MustGCPOperationTimeline(
		testCtx,
		expectedInstancePath,
		"create",
		"op-12345",
	)
	expectedDeleteOpPath := googlecloudcommon_contract.MustGCPOperationTimeline(
		testCtx,
		expectedInstancePath,
		"delete",
		"op-67890",
	)
	expectedPatchOpPath := googlecloudcommon_contract.MustGCPOperationTimeline(
		testCtx,
		expectedInstancePath,
		"patch",
		"op-patch-1",
	)

	testCases := []struct {
		name      string
		inputLog  *log.Log
		projectID string
		state     *cloudSQLAuditTimelineState
		assert    func(t *testing.T, cs *khifilev6.TimelineChangeSet)
	}{
		{
			name: "create operation start generates provisioning revision and operation start revision",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:      "my-tenant-project-tp",
					MethodName:     "cloudsql.instances.create",
					OperationID:    "op-12345",
					OperationFirst: true,
					ResourceName:   "projects/my-tenant-project-tp/instances/sql-inst",
					Status:         0,
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbCreate,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceProvisioning,
					}).
					HasRevision(expectedCreateOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationStart,
						StateType:   googlecloudcommon_contract.RevisionStateOperationStarted,
					})
			},
		},
		{
			name: "create operation finish with prior start generates existing revision and operation succeed revision",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:     "my-tenant-project-tp",
					MethodName:    "cloudsql.instances.create",
					OperationID:   "op-12345",
					OperationLast: true,
					ResourceName:  "projects/my-tenant-project-tp/instances/sql-inst",
					Status:        0,
				},
			),
			projectID: "my-tenant-project-tp",
			state: func() *cloudSQLAuditTimelineState {
				s := newCloudSQLAuditTimelineState()
				s.OperationDatabaseIDs["op-12345"] = "us-central1-composer-2-170061d7-sql"
				s.Tracker.ProcessOperationLog(context.Background(), khifilev6.NewTimelineChangeSet(testlog.NewMockLog()), expectedCreateOpPath, &googlecloudcommon_contract.GCPAuditLogFieldSet{
					OperationID:    "op-12345",
					OperationFirst: true,
				}, time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC))
				return s
			}(),
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbCreate,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceExisting,
					}).
					HasRevision(expectedCreateOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationFinish,
						StateType:   googlecloudcommon_contract.RevisionStateOperationSucceed,
					})
			},
		},
		{
			name: "create operation finish without start generates log not found revisions",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:     "my-tenant-project-tp",
					MethodName:    "cloudsql.instances.create",
					OperationID:   "op-12345",
					OperationLast: true,
					ResourceName:  "projects/my-tenant-project-tp/instances/sql-inst",
					Status:        0,
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Unix(0, 0),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbCreate,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceProvisioningLogNotFound,
					}).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbCreate,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceExisting,
					}).
					HasRevision(expectedCreateOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Unix(0, 0),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationStart,
						StateType:   googlecloudcommon_contract.RevisionStateOperationStartedLogNotFound,
					}).
					HasRevision(expectedCreateOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationFinish,
						StateType:   googlecloudcommon_contract.RevisionStateOperationSucceed,
					})
			},
		},
		{
			name: "create immediate operation generates existing revision without operation timeline",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:      "my-tenant-project-tp",
					MethodName:     "cloudsql.instances.create",
					OperationFirst: true,
					OperationLast:  true,
					ResourceName:   "projects/my-tenant-project-tp/instances/sql-inst",
					Status:         0,
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbCreate,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceExisting,
					})
			},
		},
		{
			name: "create operation failed generates event on instance timeline and failed operation revision",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:     "my-tenant-project-tp",
					MethodName:    "cloudsql.instances.create",
					OperationID:   "op-12345",
					OperationLast: true,
					ResourceName:  "projects/my-tenant-project-tp/instances/sql-inst",
					Status:        3, // INVALID_ARGUMENT
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(expectedInstancePath).
					HasNoRevision(expectedInstancePath).
					HasRevision(expectedCreateOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 5, 0, 0, time.UTC),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationFinish,
						StateType:   googlecloudcommon_contract.RevisionStateOperationFailed,
					})
			},
		},
		{
			name: "delete operation start generates deleting revision and operation start revision",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 10, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:      "my-tenant-project-tp",
					MethodName:     "cloudsql.instances.delete",
					OperationID:    "op-67890",
					OperationFirst: true,
					ResourceName:   "projects/my-tenant-project-tp/instances/sql-inst",
					Status:         0,
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Unix(0, 0),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbCreate,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceExistingLogNotFound,
					}).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 10, 0, 0, time.UTC),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbDelete,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceDeleting,
					}).
					HasRevision(expectedDeleteOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 10, 0, 0, time.UTC),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationStart,
						StateType:   googlecloudcommon_contract.RevisionStateOperationStarted,
					})
			},
		},
		{
			name: "delete operation finish with prior start generates deleted revision and operation succeed revision",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 15, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:     "my-tenant-project-tp",
					MethodName:    "cloudsql.instances.delete",
					OperationID:   "op-67890",
					OperationLast: true,
					ResourceName:  "projects/my-tenant-project-tp/instances/sql-inst",
					Status:        0,
				},
			),
			projectID: "my-tenant-project-tp",
			state: func() *cloudSQLAuditTimelineState {
				s := newCloudSQLAuditTimelineState()
				s.OperationDatabaseIDs["op-67890"] = "us-central1-composer-2-170061d7-sql"
				s.Tracker.MarkResourceRevision(expectedInstancePath)
				s.Tracker.ProcessOperationLog(context.Background(), khifilev6.NewTimelineChangeSet(testlog.NewMockLog()), expectedDeleteOpPath, &googlecloudcommon_contract.GCPAuditLogFieldSet{
					OperationID:    "op-67890",
					OperationFirst: true,
				}, time.Date(2026, 8, 4, 15, 10, 0, 0, time.UTC))
				return s
			}(),
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 15, 0, 0, time.UTC),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbDelete,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceDeleted,
					}).
					HasRevision(expectedDeleteOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 15, 0, 0, time.UTC),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationFinish,
						StateType:   googlecloudcommon_contract.RevisionStateOperationSucceed,
					})
			},
		},
		{
			name: "delete immediate operation generates deleted revision without operation timeline",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 10, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:      "my-tenant-project-tp",
					MethodName:     "cloudsql.instances.delete",
					OperationFirst: true,
					OperationLast:  true,
					ResourceName:   "projects/my-tenant-project-tp/instances/sql-inst",
					Status:         0,
				},
			),
			projectID: "my-tenant-project-tp",
			state: func() *cloudSQLAuditTimelineState {
				s := newCloudSQLAuditTimelineState()
				s.Tracker.MarkResourceRevision(expectedInstancePath)
				return s
			}(),
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasRevision(expectedInstancePath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 10, 0, 0, time.UTC),
						Principal:   "unknown",
						VerbType:    commonlogk8saudit_contract.VerbDelete,
						StateType:   privatecomposer_contract.RevisionStateCloudSQLInstanceDeleted,
					})
			},
		},
		{
			name: "patch long running operation generates event on instance timeline and operation timeline",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 20, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:      "my-tenant-project-tp",
					MethodName:     "cloudsql.instances.patch",
					OperationID:    "op-patch-1",
					OperationFirst: true,
					ResourceName:   "projects/my-tenant-project-tp/instances/sql-inst",
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(expectedInstancePath).
					HasRevision(expectedPatchOpPath, &khifilev6.StagingRevision{
						ChangedTime: time.Date(2026, 8, 4, 15, 20, 0, 0, time.UTC),
						Principal:   "",
						VerbType:    googlecloudcommon_contract.VerbOperationStart,
						StateType:   googlecloudcommon_contract.RevisionStateOperationStarted,
					})
			},
		},
		{
			name: "other method immediate generates event on instance timeline",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 20, 0, 0, time.UTC),
				privatecomposer_contract.CloudSQLFieldSet{
					DatabaseID: "us-central1-composer-2-170061d7-sql",
				},
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:    "my-tenant-project-tp",
					MethodName:   "cloudsql.instances.patch",
					ResourceName: "projects/my-tenant-project-tp/instances/sql-inst",
				},
			),
			projectID: "my-tenant-project-tp",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(expectedInstancePath)
			},
		},
		{
			name: "fallback to databaseID stored in state for same OperationID when log is missing databaseID",
			inputLog: testlog.NewMockLog(
				time.Date(2026, 8, 4, 15, 30, 0, 0, time.UTC),
				googlecloudcommon_contract.GCPAuditLogFieldSet{
					ProjectID:    "my-tenant-project-tp",
					MethodName:   "cloudsql.instances.patch",
					OperationID:  "op-patch-1",
					ResourceName: "projects/my-tenant-project-tp/instances/sql-inst",
				},
			),
			projectID: "",
			state: func() *cloudSQLAuditTimelineState {
				s := newCloudSQLAuditTimelineState()
				s.OperationDatabaseIDs["op-patch-1"] = "us-central1-composer-2-170061d7-sql"
				return s
			}(),
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(expectedInstancePath)
			},
		},
	}

	mapper := &cloudSQLAuditLogsTimelineMapper{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := khictx.WithValue(context.Background(), inspectioncore_contract.Builder, builder)
			taskDependentValues := typedmap.NewTypedMap()
			typedmap.Set(taskDependentValues, typedmap.NewTypedKey[string](privatecomposer_contract.InputComposerTenantProjectIdTaskID.ReferenceIDString()), tc.projectID)
			ctx = khictx.WithValue(ctx, core_contract.TaskResultMapContextKey, taskDependentValues)

			cs, _, err := mapper.ProcessLogByGroup(ctx, tc.inputLog, tc.state)
			if err != nil {
				t.Fatalf("ProcessLogByGroup() returned unexpected error: %v", err)
			}
			tc.assert(t, cs)
		})
	}
}
