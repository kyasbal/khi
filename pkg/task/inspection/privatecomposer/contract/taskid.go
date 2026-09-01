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
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
)

const PrivateComposerTaskIDPrefix = "privatecomposer/"

// InputComposerTenantProjectIdTaskID is the task ID for the tenant project ID of Managed Airflow.
var InputComposerTenantProjectIdTaskID = taskid.NewDefaultImplementationID[string](PrivateComposerTaskIDPrefix + "input-tenant-project-id")

// ComposerV3ClusterNamePrefixTaskID is the task id for the task that returns the GKE cluster name prefix used by Managed Airflow 3.
var ComposerV3ClusterNamePrefixTaskID = taskid.NewImplementationID(googlecloudk8scommon_contract.ClusterNamePrefixTaskRef, "gcp-composer-v3")

// CloudSQLLogsQueryTaskID is the task ID for querying Cloud SQL logs in the tenant project.
var CloudSQLLogsQueryTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateComposerTaskIDPrefix + "cloudsql-query")

// CloudSQLLogsIngesterTaskID is the task ID for ingesting Cloud SQL logs.
var CloudSQLLogsIngesterTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateComposerTaskIDPrefix + "cloudsql-log-ingester")

// CloudSQLLogsGrouperTaskID is the task ID for grouping Cloud SQL logs by database instance.
var CloudSQLLogsGrouperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateComposerTaskIDPrefix + "cloudsql-log-grouper")

// CloudSQLLogsTimelineMapperTaskID is the task ID for mapping Cloud SQL logs to timelines.
var CloudSQLLogsTimelineMapperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.TimelineMapperResult](PrivateComposerTaskIDPrefix + "cloudsql-timeline-mapper")

// CloudSQLAuditLogsQueryTaskID is the task ID for querying Cloud SQL activity audit logs in the tenant project.
var CloudSQLAuditLogsQueryTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateComposerTaskIDPrefix + "cloudsql-audit-query")

// CloudSQLAuditLogsIngesterTaskID is the task ID for ingesting Cloud SQL activity audit logs.
var CloudSQLAuditLogsIngesterTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateComposerTaskIDPrefix + "cloudsql-audit-log-ingester")

// CloudSQLAuditLogsGrouperTaskID is the task ID for grouping Cloud SQL activity audit logs by database instance.
var CloudSQLAuditLogsGrouperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateComposerTaskIDPrefix + "cloudsql-audit-log-grouper")

// CloudSQLAuditLogsTimelineMapperTaskID is the task ID for mapping Cloud SQL activity audit logs to timelines.
var CloudSQLAuditLogsTimelineMapperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.TimelineMapperResult](PrivateComposerTaskIDPrefix + "cloudsql-audit-timeline-mapper")
