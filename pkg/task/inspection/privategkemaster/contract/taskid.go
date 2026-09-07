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

package privategkemaster_contract

import (
	"github.com/GoogleCloudPlatform/khi/pkg/common/patternfinder"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/gcpqueryutil"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudlogk8snode_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8snode/contract"
)

// PrivateGKEMasterCommonTaskIDPrefix is the prefix for private GKE master tasks.
const PrivateGKEMasterCommonTaskIDPrefix = "private.khi.google.com/gke/master/"

// InputGKEMasterLogSourceTaskID is the task ID to input the master log source.
var InputGKEMasterLogSourceTaskID = taskid.NewDefaultImplementationID[*LogSource](PrivateGKEMasterCommonTaskIDPrefix + "log-source")

// InputPrivateGKEMasterComponentNameFilterTaskID is the task ID to input the master component name filter.
var InputPrivateGKEMasterComponentNameFilterTaskID = taskid.NewDefaultImplementationID[*gcpqueryutil.SetFilterParseResult](PrivateGKEMasterCommonTaskIDPrefix + "component-name-filter")

// ListLogEntriesTaskID is the task ID to list log entries.
var ListLogEntriesTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateGKEMasterCommonTaskIDPrefix + "query")

// LogIngesterTaskID is the task ID to finalize the logs to be included in the final output.
var LogIngesterTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "log-ingester")

// LogGrouperTaskID is deprecated.
var LogGrouperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateGKEMasterCommonTaskIDPrefix + "log-grouper")

// LogToTimelineMapperTaskID is deprecated.
var LogToTimelineMapperTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "log-to-timeline-mapper")

// Scheduler

// SchedulerLogFilterTaskID is the task ID to filter scheduler logs.
var SchedulerLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateGKEMasterCommonTaskIDPrefix + "scheduler/filter")

// SchedulerLogGrouperTaskID is the task ID to group scheduler logs.
var SchedulerLogGrouperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateGKEMasterCommonTaskIDPrefix + "scheduler/grouper")

// SchedulerLogToTimelineMapperTaskID is the task ID to map scheduler logs to the timeline.
var SchedulerLogToTimelineMapperTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "scheduler/log-to-timeline-mapper")

// Controller Manager

// ControllerManagerLogFilterTaskID is the task ID to filter controller manager logs.
var ControllerManagerLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateGKEMasterCommonTaskIDPrefix + "controller-manager/filter")

// ControllerManagerGrouperTaskID is the task ID to group controller manager logs.
var ControllerManagerGrouperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateGKEMasterCommonTaskIDPrefix + "controller-manager/grouper")

// ControllerManagerLogToTimelineMapperTaskID is the task ID to map controller manager logs to the timeline.
var ControllerManagerLogToTimelineMapperTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "controller-manager/log-to-timeline-mapper")

// Other

// OtherLogFilterTaskID is the task ID to filter other logs.
var OtherLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateGKEMasterCommonTaskIDPrefix + "other/filter")

// OtherGrouperTaskID is the task ID to group other logs.
var OtherGrouperTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateGKEMasterCommonTaskIDPrefix + "other/grouper")

// OtherLogToTimelineMapperTaskID is the task ID to map other logs to the timeline.
var OtherLogToTimelineMapperTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "other/log-to-timeline-mapper")

// TailTaskID is the task ID to ensure all logs are processed.
var TailTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "tail")

// Kubelet task IDs

// KubeletLogFilterTaskID is the task ID to filter kubelet logs.
var KubeletLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateGKEMasterCommonTaskIDPrefix + "kubelet/filter")

// KubeletLogGroupTaskID is the task ID to group kubelet logs.
var KubeletLogGroupTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateGKEMasterCommonTaskIDPrefix + "kubelet/grouper")

// KubeletLogLogToTimelineMapperTaskID is the task ID to map kubelet logs to the timeline.
var KubeletLogLogToTimelineMapperTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "kubelet/log-to-timeline-mapper")

// Containerd task IDs

// ContainerdLogFilterTaskID is the task ID to filter containerd logs.
var ContainerdLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](PrivateGKEMasterCommonTaskIDPrefix + "containerd/filter")

// ContainerdLogGroupTaskID is the task ID to group containerd logs.
var ContainerdLogGroupTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](PrivateGKEMasterCommonTaskIDPrefix + "containerd/grouper")

// ContainerdLogLogToTimelineMapperTaskID is the task ID to map containerd logs to the timeline.
var ContainerdLogLogToTimelineMapperTaskID = taskid.NewDefaultImplementationID[struct{}](PrivateGKEMasterCommonTaskIDPrefix + "containerd/log-to-timeline-mapper")

// PodSandboxIDDiscoveryTaskID is the task ID to discover pod sandbox IDs.
var PodSandboxIDDiscoveryTaskID = taskid.NewDefaultImplementationID[patternfinder.PatternFinder[*googlecloudlogk8snode_contract.PodSandboxIDInfo]](PrivateGKEMasterCommonTaskIDPrefix + "containerd/pod-sandbox-id-discovery")

// ContainerIDDiscoveryTaskID is the task ID to discover container IDs.
var ContainerIDDiscoveryTaskID = taskid.NewDefaultImplementationID[commonlogk8saudit_contract.ContainerIDToContainerIdentity](PrivateGKEMasterCommonTaskIDPrefix + "containerd/container-id-discovery")
