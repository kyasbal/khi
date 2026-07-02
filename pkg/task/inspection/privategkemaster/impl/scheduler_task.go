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

package privategkemaster_impl

import (
	"context"

	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/logutil"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	googlecloudlogk8scontrolplane_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8scontrolplane/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
)

var schedulerLogFilterTask = inspectiontaskbase.NewLogFilterTask(
	privategkemaster_contract.SchedulerLogFilterTaskID,
	privategkemaster_contract.CommonFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, l *log.Log) bool {
		componentFieldSet, err := log.GetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
		if err != nil {
			return false
		}
		return componentFieldSet.PrivateGKEMasterParserType() == privategkemaster_contract.PrivateGKEMasterParserTypeScheduler
	},
)

var schedulerLogFieldSetReaderTask = inspectiontaskbase.NewFieldSetReadTask(privategkemaster_contract.SchedulerLogFieldSetReaderTaskID,
	privategkemaster_contract.SchedulerLogFilterTaskID.Ref(),
	[]log.FieldSetReader{
		&googlecloudlogk8scontrolplane_contract.K8sSchedulerComponentFieldSetReader{
			KLogParser: logutil.NewKLogTextParser(false),
		},
	},
)

var schedulerGrouperTask = inspectiontaskbase.NewLogGrouperTask(
	privategkemaster_contract.SchedulerLogGrouperTaskID,
	privategkemaster_contract.SchedulerLogFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, l *log.Log) string {
		return "" // No grouping needed
	},
)

// SchedulerTimelineMapper maps scheduler logs to timeline paths.
type SchedulerTimelineMapper struct {
	inspectiontaskbase.StatelessMapperBase
}

// Dependencies implements inspectiontaskbase.LogToTimelineMapper.
func (m *SchedulerTimelineMapper) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(),
	}
}

// GroupedLogTask implements inspectiontaskbase.LogToTimelineMapper.
func (m *SchedulerTimelineMapper) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privategkemaster_contract.SchedulerLogGrouperTaskID.Ref()
}

// LogIngesterTask implements inspectiontaskbase.LogToTimelineMapper.
func (m *SchedulerTimelineMapper) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.LogIngesterTaskID.Ref()
}

// ProcessLogByGroup implements inspectiontaskbase.LogToTimelineMapper.
func (m *SchedulerTimelineMapper) ProcessLogByGroup(ctx context.Context, l *log.Log, _ struct{}) (*khifilev6.TimelineChangeSet, struct{}, error) {
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref())
	masterFieldSet, err := log.GetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
	if err != nil {
		return nil, struct{}{}, err
	}
	schedulerMessageFieldSet, err := log.GetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sSchedulerComponentFieldSet{})
	if err != nil {
		return nil, struct{}{}, err
	}

	cs := khifilev6.NewTimelineChangeSet(l)
	for _, tPath := range masterFieldSet.ResourceTimelines(ctx, clusterIdentity.ClusterName) {
		cs.AddEvent(tPath)
	}
	if schedulerMessageFieldSet.HasPodField() {
		clusterTimeline := commonlogk8saudit_contract.MustK8sClusterTimeline(ctx, clusterIdentity.ClusterName)
		apiVersionTimeline := commonlogk8saudit_contract.MustK8sAPIVersionTimeline(ctx, clusterTimeline, "core/v1")
		kindTimeline := commonlogk8saudit_contract.MustK8sKindTimeline(ctx, apiVersionTimeline, "pod")
		namespaceTimeline := commonlogk8saudit_contract.MustK8sNamespaceTimeline(ctx, kindTimeline, schedulerMessageFieldSet.PodNamespace)
		podTimeline := commonlogk8saudit_contract.MustK8sNamespacedResourceTimeline(ctx, namespaceTimeline, schedulerMessageFieldSet.PodName)
		cs.AddEvent(podTimeline)
	}
	return cs, struct{}{}, nil
}

var _ inspectiontaskbase.LogToTimelineMapper[struct{}] = (*SchedulerTimelineMapper)(nil)

var schedulerLogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask(
	privategkemaster_contract.SchedulerLogToTimelineMapperTaskID,
	&SchedulerTimelineMapper{},
)
