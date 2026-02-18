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

	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
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
		&googlecloudlogk8scontrolplane_contract.K8sSchedulerComponentFieldSetReader{},
	},
)

var schedulerGrouperTask = inspectiontaskbase.NewLogGrouperTask(
	privategkemaster_contract.SchedulerLogGrouperTaskID,
	privategkemaster_contract.SchedulerLogFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, l *log.Log) string {
		return "" // No grouping needed
	},
)

var schedulerLogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask[struct{}](
	privategkemaster_contract.SchedulerLogToTimelineMapperTaskID,
	&schedulerLogToTimelineMapperTaskSetting{},
)

type schedulerLogToTimelineMapperTaskSetting struct {
}

// Dependencies implements inspectiontaskbase.LogToTimelineMapper.
func (o *schedulerLogToTimelineMapperTaskSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		googlecloudk8scommon_contract.ClusterIndentityTaskID.Ref(),
	}
}

// GroupedLogTask implements inspectiontaskbase.LogToTimelineMapper.
func (o *schedulerLogToTimelineMapperTaskSetting) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privategkemaster_contract.SchedulerLogGrouperTaskID.Ref()
}

// LogIngesterTask implements inspectiontaskbase.LogToTimelineMapper.
func (o *schedulerLogToTimelineMapperTaskSetting) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.LogIngesterTaskID.Ref()
}

// ProcessLogByGroup implements inspectiontaskbase.LogToTimelineMapper.
func (o *schedulerLogToTimelineMapperTaskSetting) ProcessLogByGroup(ctx context.Context, l *log.Log, cs *history.ChangeSet, builder *history.Builder, prevGroupData struct{}) (struct{}, error) {
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIndentityTaskID.Ref())
	masterFieldSet, err := log.GetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
	if err != nil {
		return struct{}{}, err
	}
	commonMainMessage, err := log.GetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{})
	if err != nil {
		return struct{}{}, err
	}
	schedulerMessageFieldSet, err := log.GetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sSchedulerComponentFieldSet{})
	if err != nil {
		return struct{}{}, err
	}

	cs.SetLogSummary(commonMainMessage.Message)
	for _, path := range masterFieldSet.ResourcePaths(clusterIdentity.ClusterName) {
		cs.AddEvent(path)
	}
	if schedulerMessageFieldSet.HasPodField() {
		cs.AddEvent(schedulerMessageFieldSet.ResourcePath())
	}
	return struct{}{}, nil
}

var _ inspectiontaskbase.LogToTimelineMapper[struct{}] = (*schedulerLogToTimelineMapperTaskSetting)(nil)
