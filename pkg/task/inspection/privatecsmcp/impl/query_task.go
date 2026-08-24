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

package privatecsmcp_impl

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud/logestimator"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
)

// LogQueryTask executes Cloud Logging filter to fetch CSM CP logs.
var LogQueryTask = googlecloudcommon_contract.NewStructuredListLogEntriesTask(&csmcpLogQueryTaskSetting{})

type csmcpLogQueryTaskSetting struct{}

// TaskID returns the ID of this task.
func (s *csmcpLogQueryTaskSetting) TaskID() taskid.TaskImplementationID[[]*log.Log] {
	return privatecsmcp_contract.LogQueryTaskID
}

// Dependencies returns the dependencies required by this task.
func (s *csmcpLogQueryTaskSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(),
		privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref(),
	}
}

// QueryName returns the human-readable name of the query task.
func (s *csmcpLogQueryTaskSetting) QueryName() string {
	return "CSM CP logs"
}

// Queries returns the list of structured log queries for CSM CP logs.
func (s *csmcpLogQueryTaskSetting) Queries(ctx context.Context) ([]*logestimator.StructuredLogQuery, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref())
	serviceName := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref())

	return []*logestimator.StructuredLogQuery{
		GenerateCSMCPStructuredQuery(tenantProjectID, serviceName),
	}, nil
}

// DefaultResourceNames returns the resource names to query logs from.
func (s *csmcpLogQueryTaskSetting) DefaultResourceNames(ctx context.Context) ([]string, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref())
	if tenantProjectID == "" {
		return []string{}, nil
	}
	return []string{fmt.Sprintf("projects/%s", tenantProjectID)}, nil
}

// TimePartitionCount returns the number of time partitions to use.
func (s *csmcpLogQueryTaskSetting) TimePartitionCount(ctx context.Context) (int, error) {
	return 10, nil
}

var _ googlecloudcommon_contract.StructuredListLogEntriesTaskSetting = (*csmcpLogQueryTaskSetting)(nil)

// GenerateCSMCPStructuredQuery generates a structured query for CSM CP logs.
func GenerateCSMCPStructuredQuery(tenantProjectID, serviceName string) *logestimator.StructuredLogQuery {
	filters := []logestimator.LoggingMonitoringMatcher{
		logestimator.ResourceLabel("project_id", logestimator.Exact(tenantProjectID)),
		logestimator.ResourceLabel("service_name", logestimator.Exact(serviceName)),
		logestimator.LogID(logestimator.OneOf(
			"run.googleapis.com/stdout",
			"run.googleapis.com/stderr",
		)),
	}

	return &logestimator.StructuredLogQuery{
		Incomplete:    tenantProjectID == "" || serviceName == "",
		ResourceTypes: []string{"cloud_run_revision"},
		Filters:       filters,
	}
}

// GenerateCSMCPQuery formats the Cloud Logging filter query for CSM CP logs.
func GenerateCSMCPQuery(tenantProjectID, serviceName string) string {
	return GenerateCSMCPStructuredQuery(tenantProjectID, serviceName).GenerateCloudLoggingQuery()
}
