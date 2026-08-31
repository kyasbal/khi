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
	"strings"
	"time"

	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
)

// LogSorterTask sorts logs by time.
var LogSorterTask = inspectiontaskbase.NewLogSorterByTimeTask(
	privatecsmcp_contract.LogSorterTaskID,
	privatecsmcp_contract.LogQueryTaskID.Ref(),
)

type csmcpLogIngester struct{}

// RawLogTask returns the task ID for the raw log input.
func (i *csmcpLogIngester) RawLogTask() taskid.TaskReference[[]*log.Log] {
	return privatecsmcp_contract.LogSorterTaskID.Ref()
}

// Dependencies returns additional dependencies required for ingestion.
func (i *csmcpLogIngester) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{}
}

// ProcessLog ingests a single log entry into a LogChangeSet.
func (i *csmcpLogIngester) ProcessLog(ctx context.Context, l *log.Log) (*khifilev6.LogChangeSet, error) {
	cs, err := khifilev6.NewLogChangeSet(l)
	if err != nil {
		return nil, err
	}

	cs.SetLogType(privatecsmcp_contract.LogTypeCSMCP)
	cs.SetTimestamp(l.Timestamp)

	if severity, err := googlecloudcommon_contract.ExtractGCPSeverity(l.NodeReader); err == nil && severity != nil {
		cs.SetSeverity(severity)
	}

	if csmcpFS, err := privatecsmcp_contract.ExtractCSMCP(l.NodeReader); err == nil {
		cs.SetSummary(csmcpFS.Message)
		if csmcpFS.Timestamp != nil {
			cs.SetTimestamp(*csmcpFS.Timestamp)
		}
	}

	return cs, nil
}

var _ inspectiontaskbase.LogIngester = (*csmcpLogIngester)(nil)

// LogIngesterTask ingests CSM CP logs.
var LogIngesterTask = inspectiontaskbase.NewLogIngesterTask(
	privatecsmcp_contract.LogIngesterTaskID,
	&csmcpLogIngester{},
)

// LogGrouperTask groups logs, currently all in one group since service is filtered.
var LogGrouperTask = inspectiontaskbase.NewLogGrouperTask(
	privatecsmcp_contract.LogGrouperTaskID,
	privatecsmcp_contract.LogSorterTaskID.Ref(),
	func(ctx context.Context, l *log.Log) string {
		return "csmcp"
	},
)

type csmcpTimelineState struct {
	// ConnectedConns tracks connection IDs that have a "connected" log.
	ConnectedConns map[string]bool
}

type csmcpTimelineMapper struct{}

// PassCount returns 1 to perform one pre-processing pass.
func (m *csmcpTimelineMapper) PassCount() int {
	return 1
}

func getConnKey(pod privatecsmcp_contract.PodIdentifier) string {
	return pod.Namespace + "/" + pod.Name + "/" + pod.ConnectionID
}

// isXDSLog checks if the message starts with an xDS prefix (e.g., ADS:, CDS:).
func isXDSLog(msg string) bool {
	idx := strings.IndexByte(msg, ':')
	if idx <= 0 || idx > 10 {
		return false
	}
	prefix := msg[:idx]
	for i := 0; i < len(prefix); i++ {
		c := prefix[i]
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

// PreProcessLogByGroup checks for "connected" logs.
func (m *csmcpTimelineMapper) PreProcessLogByGroup(ctx context.Context, passIndex int, l *log.Log, prevGroupData *csmcpTimelineState) (*csmcpTimelineState, error) {
	if prevGroupData == nil {
		prevGroupData = &csmcpTimelineState{
			ConnectedConns: make(map[string]bool),
		}
	}

	csmcpFS, err := privatecsmcp_contract.ExtractCSMCP(l.NodeReader)
	if err != nil {
		return prevGroupData, nil // skip if not csmcp log
	}

	msg := csmcpFS.Message
	if isXDSLog(msg) {
		if strings.Contains(msg, "new connection for") {
			for _, pod := range csmcpFS.Pods {
				if pod.ConnectionID != "" {
					prevGroupData.ConnectedConns[getConnKey(pod)] = true
				}
			}
		}
	}
	return prevGroupData, nil
}

// LogIngesterTask returns the task ID for log ingestion.
func (m *csmcpTimelineMapper) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privatecsmcp_contract.LogIngesterTaskID.Ref()
}

// Dependencies returns dependencies needed for mapping.
func (m *csmcpTimelineMapper) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(),
		privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref(),
		googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(),
	}
}

// GroupedLogTask returns the task ID for grouped logs.
func (m *csmcpTimelineMapper) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privatecsmcp_contract.LogGrouperTaskID.Ref()
}

// ProcessLogByGroup maps a log entry to its corresponding timeline paths.
func (m *csmcpTimelineMapper) ProcessLogByGroup(ctx context.Context, l *log.Log, prevGroupData *csmcpTimelineState) (*khifilev6.TimelineChangeSet, *csmcpTimelineState, error) {
	tenantProjectID := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref())
	serviceName := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref())
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref())

	csmcpFS, err := privatecsmcp_contract.ExtractCSMCP(l.NodeReader)
	if err != nil {
		return nil, prevGroupData, err
	}

	cs := khifilev6.NewTimelineChangeSet(l)

	tenantProjectPath := googlecloudcommon_contract.MustGCPProjectTimeline(ctx, tenantProjectID)
	servicePath := privatecsmcp_contract.MustCloudRunServiceTimeline(ctx, tenantProjectPath, serviceName, csmcpFS.InstanceID)

	cs.AddEvent(servicePath)

	// Map to Pods' timelines if pod names are extracted and it's an xDS log
	if isXDSLog(csmcpFS.Message) && clusterIdentity.ClusterName != "" && len(csmcpFS.Pods) > 0 {
		clusterPath := commonlogk8saudit_contract.MustK8sClusterTimeline(ctx, clusterIdentity.ClusterName)
		apiPath := commonlogk8saudit_contract.MustK8sAPIVersionTimeline(ctx, clusterPath, "core/v1")
		kindPath := commonlogk8saudit_contract.MustK8sKindTimeline(ctx, apiPath, "pod")

		for _, pod := range csmcpFS.Pods {
			podName := pod.Name
			namespace := pod.Namespace

			nsPath := commonlogk8saudit_contract.MustK8sNamespaceTimeline(ctx, kindPath, namespace)
			podPath := commonlogk8saudit_contract.MustK8sNamespacedResourceTimeline(ctx, nsPath, podName)
			csmcpPodLogPath := privatecsmcp_contract.MustCSMCPPodLogTimeline(ctx, podPath)

			cs.AddEvent(csmcpPodLogPath)

			if pod.ConnectionID != "" {
				connPath := privatecsmcp_contract.MustCSMCPConnectionTimeline(ctx, podPath, pod.ConnectionID)

				msg := csmcpFS.Message
				if isXDSLog(msg) {
					var changedTime time.Time
					if csmcpFS.Timestamp != nil {
						changedTime = *csmcpFS.Timestamp
					}
					if strings.Contains(msg, "new connection for") {
						cs.AddRevision(connPath, &khifilev6.StagingRevision{
							VerbType:     inspectioncore_contract.VerbUnknown,
							StateType:    privatecsmcp_contract.RevisionStateCSMCPConnectionConnected,
							ResourceBody: nil,
							Principal:    "csm-cp",
							ChangedTime:  changedTime,
						})
					} else if strings.Contains(msg, "terminated") {
						if prevGroupData != nil && !prevGroupData.ConnectedConns[getConnKey(pod)] {
							cs.AddRevision(connPath, &khifilev6.StagingRevision{
								ChangedTime:  time.Unix(0, 0),
								VerbType:     inspectioncore_contract.VerbUnknown,
								StateType:    privatecsmcp_contract.RevisionStateCSMCPConnectionConnectedLogNotFound,
								ResourceBody: nil,
								Principal:    "csm-cp",
							})
						}
						cs.AddRevision(connPath, &khifilev6.StagingRevision{
							VerbType:     inspectioncore_contract.VerbUnknown,
							StateType:    privatecsmcp_contract.RevisionStateCSMCPConnectionTerminated,
							ResourceBody: nil,
							Principal:    "csm-cp",
							ChangedTime:  changedTime,
						})
					}
				}
			}
		}
	}

	return cs, prevGroupData, nil
}

// LogToTimelineMapperTask is the task that maps CSM CP logs to timelines.
var LogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask(
	privatecsmcp_contract.LogToTimelineMapperTaskID,
	&csmcpTimelineMapper{},
	inspectioncore_contract.FeatureTaskLabel(
		"CSM Control Plane Logs (PRIVATE)",
		`Parser and timeline mapping for CSM Control Plane logs in tenant project.
It gathers xDS logs from the Cloud Run in the tenant project and associate them with Pods related to them.
This log needs to query meshconfig.googleapis.com tenant project. 
`,
		21000,
		false,
	),
)

var _ inspectiontaskbase.LogToTimelineMapper[*csmcpTimelineState] = (*csmcpTimelineMapper)(nil)
