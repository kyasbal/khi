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

package privatecsmcp_contract

import (
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
)

// TaskIDPrefix is the prefix used for all CSM CP task IDs.
const TaskIDPrefix = "csmcp.khi.google.com/"

// InputCSMTenantProjectIDTaskID is the task ID for resolving the CSM Tenant Project ID input.
var InputCSMTenantProjectIDTaskID = taskid.NewDefaultImplementationID[string](TaskIDPrefix + "input-tenant-project-id")

// InputCSMCPCloudRunServiceNameTaskID is the task ID for resolving the CSM CP Cloud Run Service Name input.
var InputCSMCPCloudRunServiceNameTaskID = taskid.NewDefaultImplementationID[string](TaskIDPrefix + "input-cloudrun-service-name")

// AutocompleteCSMCPCloudRunServiceNameTaskID is the task ID for the Cloud Run Service Name autocomplete.
var AutocompleteCSMCPCloudRunServiceNameTaskID = taskid.NewDefaultImplementationID[*inspectioncore_contract.AutocompleteResult[string]](TaskIDPrefix + "autocomplete-cloudrun-service-name")

// LogQueryTaskID is the task ID for executing the log query.
var LogQueryTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "query")

// LogSorterTaskID is the task ID for sorting logs by time.
var LogSorterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "log-sorter")

// LogIngesterTaskID is the task ID for ingesting parsed logs into LogChangeSets.
var LogIngesterTaskID = taskid.NewDefaultImplementationID[struct{}](TaskIDPrefix + "log-ingester")

// LogGrouperTaskID is the task ID for grouping parsed logs.
var LogGrouperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](TaskIDPrefix + "log-grouper")

// LogToTimelineMapperTaskID is the task ID for mapping log groups to timelines.
var LogToTimelineMapperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.TimelineMapperResult](TaskIDPrefix + "timeline-mapper")
