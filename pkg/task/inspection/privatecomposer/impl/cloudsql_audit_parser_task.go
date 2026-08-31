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
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	pb "github.com/GoogleCloudPlatform/khi/pkg/generated/khifile/v6"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
)

type cloudSQLAuditLogsIngester struct{}

// RawLogTask returns the reference to the task providing raw audit logs.
func (i *cloudSQLAuditLogsIngester) RawLogTask() taskid.TaskReference[[]*log.Log] {
	return privatecomposer_contract.CloudSQLAuditLogsQueryTaskID.Ref()
}

// Dependencies returns additional task dependencies required by the ingester.
func (i *cloudSQLAuditLogsIngester) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{}
}

// ProcessLog parses log entries and populates change set metadata for Cloud SQL audit logs.
func (i *cloudSQLAuditLogsIngester) ProcessLog(ctx context.Context, l *log.Log) (*khifilev6.LogChangeSet, error) {
	cs, err := khifilev6.NewLogChangeSet(l)
	if err != nil {
		return nil, err
	}

	cs.SetLogType(privatecomposer_contract.LogTypeCloudSQLAudit)
	cs.SetTimestamp(l.Timestamp)

	if sev, err := googlecloudcommon_contract.ExtractGCPSeverity(l.NodeReader); err == nil && sev != nil {
		cs.SetSeverity(sev)
	}

	summarySet := false
	if mainMsg, err := googlecloudcommon_contract.ExtractGCPMainMessage(l.NodeReader); err == nil && mainMsg != "" {
		cs.SetSummary(mainMsg)
		summarySet = true
	}
	if !summarySet {
		if cloudSQLFS, err := privatecomposer_contract.ExtractCloudSQL(l.NodeReader); err == nil && cloudSQLFS.Summary != "" {
			cs.SetSummary(cloudSQLFS.Summary)
		}
	}

	return cs, nil
}

var _ inspectiontaskbase.LogIngester = (*cloudSQLAuditLogsIngester)(nil)

// CloudSQLAuditLogsIngesterTask is the log ingester task for Cloud SQL activity audit logs.
var CloudSQLAuditLogsIngesterTask = inspectiontaskbase.NewLogIngesterTask(
	privatecomposer_contract.CloudSQLAuditLogsIngesterTaskID,
	&cloudSQLAuditLogsIngester{},
)

// CloudSQLAuditLogsGrouperTask groups Cloud SQL activity audit logs into a single group to share state across operation logs.
var CloudSQLAuditLogsGrouperTask = inspectiontaskbase.NewLogGrouperTask(
	privatecomposer_contract.CloudSQLAuditLogsGrouperTaskID,
	privatecomposer_contract.CloudSQLAuditLogsQueryTaskID.Ref(),
	func(ctx context.Context, l *log.Log) string {
		return "cloudsql-audit"
	},
)

type cloudSQLAuditTimelineState struct {
	Tracker              *googlecloudcommon_contract.GCPOperationTracker
	OperationDatabaseIDs map[string]string
}

func newCloudSQLAuditTimelineState() *cloudSQLAuditTimelineState {
	return &cloudSQLAuditTimelineState{
		Tracker:              googlecloudcommon_contract.NewGCPOperationTracker(),
		OperationDatabaseIDs: make(map[string]string),
	}
}

type cloudSQLAuditLogsTimelineMapper struct {
	inspectiontaskbase.SinglePassMapperBase[*cloudSQLAuditTimelineState]
}

// LogIngesterTask returns the reference to the Cloud SQL audit log ingester task.
func (m *cloudSQLAuditLogsTimelineMapper) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privatecomposer_contract.CloudSQLAuditLogsIngesterTaskID.Ref()
}

// Dependencies returns additional dependencies needed for mapping audit log timelines.
func (m *cloudSQLAuditLogsTimelineMapper) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{}
}

// GroupedLogTask returns the reference to the Cloud SQL audit log grouper task.
func (m *cloudSQLAuditLogsTimelineMapper) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privatecomposer_contract.CloudSQLAuditLogsGrouperTaskID.Ref()
}

// ProcessLogByGroup maps individual Cloud SQL activity audit logs to timeline events, operations, and revisions.
func (m *cloudSQLAuditLogsTimelineMapper) ProcessLogByGroup(ctx context.Context, l *log.Log, state *cloudSQLAuditTimelineState) (*khifilev6.TimelineChangeSet, *cloudSQLAuditTimelineState, error) {
	if state == nil {
		state = newCloudSQLAuditTimelineState()
	}

	audit, err := googlecloudcommon_contract.ExtractGCPAuditLog(l.NodeReader)
	if err != nil {
		return nil, state, err
	}

	tenantProjectID := audit.ProjectID
	if tenantProjectID == "" {
		tenantProjectID = "unknown"
	}

	var databaseID string
	if cloudSQLFS, err := privatecomposer_contract.ExtractCloudSQL(l.NodeReader); err == nil && cloudSQLFS.DatabaseID != "" {
		databaseID = cloudSQLFS.DatabaseID
		if audit.OperationID != "" && !audit.ImmediateOperation() {
			state.OperationDatabaseIDs[audit.OperationID] = databaseID
		}
	} else if audit.OperationID != "" {
		databaseID = state.OperationDatabaseIDs[audit.OperationID]
	}

	if audit.OperationID != "" && audit.Ending() {
		delete(state.OperationDatabaseIDs, audit.OperationID)
	}

	if databaseID == "" {
		databaseID = "unknown"
	}

	instancePath := privatecomposer_contract.MustCloudSQLInstanceTimeline(ctx, tenantProjectID, databaseID)
	cs := khifilev6.NewTimelineChangeSet(l)

	if !audit.ImmediateOperation() {
		shortMethodName := "unknown"
		if audit.MethodName != "" {
			methodNameSplitted := strings.Split(audit.MethodName, ".")
			if last := methodNameSplitted[len(methodNameSplitted)-1]; last != "" {
				shortMethodName = last
			}
		}
		operationPath := googlecloudcommon_contract.MustGCPOperationTimeline(ctx, instancePath, shortMethodName, audit.OperationID)
		state.Tracker.ProcessOperationLog(ctx, cs, operationPath, &audit, l.Timestamp)
	}

	isCreate := audit.MethodName == "cloudsql.instances.create"
	isDelete := audit.MethodName == "cloudsql.instances.delete"

	if (isCreate || isDelete) && audit.Status <= 0 {
		principal := audit.PrincipalEmail
		if principal == "" {
			principal = "unknown"
		}

		if isCreate {
			var resBody structured.Node
			if audit.Request != nil {
				resBody = audit.Request.Node
			}

			// When a creation ending log observed but no starting operation found, the creation is previously started.
			if audit.Ending() && !state.Tracker.HasStarted(audit.OperationID) {
				cs.AddRevision(instancePath, &khifilev6.StagingRevision{
					VerbType:     commonlogk8saudit_contract.VerbCreate,
					StateType:    privatecomposer_contract.RevisionStateCloudSQLInstanceProvisioningLogNotFound,
					Principal:    principal,
					ChangedTime:  time.Unix(0, 0),
					ResourceBody: nil,
				})
				state.Tracker.MarkResourceRevision(instancePath)
			}

			var stateType *pb.RevisionState
			if audit.Starting() {
				stateType = privatecomposer_contract.RevisionStateCloudSQLInstanceProvisioning
			} else {
				stateType = privatecomposer_contract.RevisionStateCloudSQLInstanceExisting
			}

			cs.AddRevision(instancePath, &khifilev6.StagingRevision{
				VerbType:     commonlogk8saudit_contract.VerbCreate,
				StateType:    stateType,
				Principal:    principal,
				ChangedTime:  l.Timestamp,
				ResourceBody: resBody,
			})
			state.Tracker.MarkResourceRevision(instancePath)
		} else if isDelete {
			// When a deletion start log observed but no other resource revision before, prev status must be existing.
			if audit.Starting() && !state.Tracker.HasResourceRevision(instancePath) {
				cs.AddRevision(instancePath, &khifilev6.StagingRevision{
					VerbType:     commonlogk8saudit_contract.VerbCreate,
					StateType:    privatecomposer_contract.RevisionStateCloudSQLInstanceExistingLogNotFound,
					Principal:    principal,
					ChangedTime:  time.Unix(0, 0),
					ResourceBody: nil,
				})
				state.Tracker.MarkResourceRevision(instancePath)
			}
			// When a deletion ending log observed but no starting operation found, the deletion is previously started.
			if audit.Ending() && !state.Tracker.HasStarted(audit.OperationID) {
				cs.AddRevision(instancePath, &khifilev6.StagingRevision{
					VerbType:     commonlogk8saudit_contract.VerbDelete,
					StateType:    privatecomposer_contract.RevisionStateCloudSQLInstanceDeletingLogNotFound,
					Principal:    principal,
					ChangedTime:  time.Unix(0, 0),
					ResourceBody: nil,
				})
				state.Tracker.MarkResourceRevision(instancePath)
			}

			var stateType *pb.RevisionState
			if audit.Starting() {
				stateType = privatecomposer_contract.RevisionStateCloudSQLInstanceDeleting
			} else {
				stateType = privatecomposer_contract.RevisionStateCloudSQLInstanceDeleted
			}

			cs.AddRevision(instancePath, &khifilev6.StagingRevision{
				VerbType:    commonlogk8saudit_contract.VerbDelete,
				StateType:   stateType,
				Principal:   principal,
				ChangedTime: l.Timestamp,
			})
			state.Tracker.MarkResourceRevision(instancePath)
		}
	} else {
		cs.AddEvent(instancePath)
	}

	return cs, state, nil
}

var _ inspectiontaskbase.LogToTimelineMapper[*cloudSQLAuditTimelineState] = (*cloudSQLAuditLogsTimelineMapper)(nil)

// CloudSQLAuditLogsTimelineMapperTask maps Cloud SQL activity audit logs to timelines.
var CloudSQLAuditLogsTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask(
	privatecomposer_contract.CloudSQLAuditLogsTimelineMapperTaskID,
	&cloudSQLAuditLogsTimelineMapper{},
	inspectioncore_contract.FeatureTaskLabel(
		"[PRIVATE] Cloud SQL Audit Logs",
		"Gather Cloud SQL activity audit logs (cloudaudit.googleapis.com/activity) from the Managed Airflow tenant project to visualize database create/delete operations and events on timelines.",
		2610,
		false,
	),
)
