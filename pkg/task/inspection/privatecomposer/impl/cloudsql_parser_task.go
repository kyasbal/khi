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

	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
)

type cloudSQLLogsIngester struct{}

// RawLogTask returns the reference to the task providing raw logs.
func (i *cloudSQLLogsIngester) RawLogTask() taskid.TaskReference[[]*log.Log] {
	return privatecomposer_contract.CloudSQLLogsQueryTaskID.Ref()
}

// Dependencies returns additional task dependencies required by the ingester.
func (i *cloudSQLLogsIngester) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{}
}

// ProcessLog parses log entries and populates change set metadata for Cloud SQL logs.
func (i *cloudSQLLogsIngester) ProcessLog(ctx context.Context, l *log.Log) (*khifilev6.LogChangeSet, error) {
	cs, err := khifilev6.NewLogChangeSet(l)
	if err != nil {
		return nil, err
	}

	cs.SetLogType(privatecomposer_contract.LogTypeCloudSQL)
	cs.SetTimestamp(l.Timestamp)

	if sev, err := googlecloudcommon_contract.ExtractGCPSeverity(l.NodeReader); err == nil && sev != nil {
		cs.SetSeverity(sev)
	}

	summarySet := false
	if cloudSQLFS, err := privatecomposer_contract.ExtractCloudSQL(l.NodeReader); err == nil && cloudSQLFS.Summary != "" {
		cs.SetSummary(cloudSQLFS.Summary)
		summarySet = true
	}
	if !summarySet {
		if mainMsg, err := googlecloudcommon_contract.ExtractGCPMainMessage(l.NodeReader); err == nil && mainMsg != "" {
			cs.SetSummary(mainMsg)
		}
	}

	return cs, nil
}

var _ inspectiontaskbase.LogIngester = (*cloudSQLLogsIngester)(nil)

// CloudSQLLogsIngesterTask is the log ingester task for Cloud SQL database logs.
var CloudSQLLogsIngesterTask = inspectiontaskbase.NewLogIngesterTask(
	privatecomposer_contract.CloudSQLLogsIngesterTaskID,
	&cloudSQLLogsIngester{},
)

// CloudSQLLogsGrouperTask groups Cloud SQL database logs by database instance ID.
var CloudSQLLogsGrouperTask = inspectiontaskbase.NewLogGrouperTask(
	privatecomposer_contract.CloudSQLLogsGrouperTaskID,
	privatecomposer_contract.CloudSQLLogsQueryTaskID.Ref(),
	func(ctx context.Context, l *log.Log) string {
		if cloudSQLFS, err := privatecomposer_contract.ExtractCloudSQL(l.NodeReader); err == nil && cloudSQLFS.DatabaseID != "" {
			return cloudSQLFS.DatabaseID
		}
		return "unknown"
	},
)

type cloudSQLLogsTimelineMapper struct {
	inspectiontaskbase.StatelessMapperBase
}

// LogIngesterTask returns the reference to the Cloud SQL log ingester task.
func (m *cloudSQLLogsTimelineMapper) LogIngesterTask() taskid.TaskReference[struct{}] {
	return privatecomposer_contract.CloudSQLLogsIngesterTaskID.Ref()
}

// Dependencies returns additional dependencies needed for mapping timelines.
func (m *cloudSQLLogsTimelineMapper) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref(),
	}
}

// GroupedLogTask returns the reference to the Cloud SQL log grouper task.
func (m *cloudSQLLogsTimelineMapper) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privatecomposer_contract.CloudSQLLogsGrouperTaskID.Ref()
}

// ProcessLogByGroup maps individual Cloud SQL log entries to timeline events.
func (m *cloudSQLLogsTimelineMapper) ProcessLogByGroup(ctx context.Context, l *log.Log, _ struct{}) (*khifilev6.TimelineChangeSet, struct{}, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref())
	if tenantProjectID == "" {
		tenantProjectID = "unknown"
	}

	databaseID := "unknown"
	logFileName := "unknown"
	if cloudSQLFS, err := privatecomposer_contract.ExtractCloudSQL(l.NodeReader); err == nil {
		if cloudSQLFS.DatabaseID != "" {
			databaseID = cloudSQLFS.DatabaseID
		}
		if cloudSQLFS.LogFileName != "" {
			logFileName = cloudSQLFS.LogFileName
		}
	}

	targetPath := privatecomposer_contract.MustCloudSQLLogTimeline(ctx, tenantProjectID, databaseID, logFileName)
	cs := khifilev6.NewTimelineChangeSet(l)
	cs.AddEvent(targetPath)

	return cs, struct{}{}, nil
}

var _ inspectiontaskbase.LogToTimelineMapper[struct{}] = (*cloudSQLLogsTimelineMapper)(nil)

// CloudSQLLogsTimelineMapperTask maps Cloud SQL database logs to timelines.
var CloudSQLLogsTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask(
	privatecomposer_contract.CloudSQLLogsTimelineMapperTaskID,
	&cloudSQLLogsTimelineMapper{},
	inspectioncore_contract.FeatureTaskLabel(
		"[PRIVATE] Cloud SQL Logs",
		"Gather Cloud SQL database engine logs (postgres.log, postgres-audit.log, postgres-upgrade.log) from the Managed Airflow tenant project to visualize database operations on timelines.",
		2600,
		false,
	),
)
