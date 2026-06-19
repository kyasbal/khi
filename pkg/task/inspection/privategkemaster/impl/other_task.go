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

	"github.com/GoogleCloudPlatform/khi/pkg/common/patternfinder"
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

var otherLogFilterTask = inspectiontaskbase.NewLogFilterTask(
	privategkemaster_contract.OtherLogFilterTaskID,
	privategkemaster_contract.CommonFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, l *log.Log) bool {
		componentFieldSet, err := log.GetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
		if err != nil {
			return false
		}
		return componentFieldSet.PrivateGKEMasterParserType() == privategkemaster_contract.PrivateGKEMasterParserTypeOther
	},
)

var otherLogFieldSetReaderTask = inspectiontaskbase.NewFieldSetReadTask(privategkemaster_contract.OtherLogFieldSetReaderTaskID,
	privategkemaster_contract.OtherLogFilterTaskID.Ref(),
	[]log.FieldSetReader{}, // No additional fields to read for "Other"
)

var otherGrouperTask = inspectiontaskbase.NewLogGrouperTask(
	privategkemaster_contract.OtherGrouperTaskID,
	privategkemaster_contract.OtherLogFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, l *log.Log) string {
		masterLogField := log.MustGetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
		return masterLogField.PodID
	},
)

// OtherTimelineMapper maps other control plane logs to timeline paths.
type OtherTimelineMapper struct {
	inspectiontaskbase.StatelessMapperBase
	uidPrefixTokenCandidates []rune
}

// Dependencies implements inspectiontaskbase.LogToTimelineMapperV2.
func (p *OtherTimelineMapper) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		commonlogk8saudit_contract.ResourceUIDPatternFinderTaskID.Ref(),
		googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(),
	}
}

// GroupedLogTask implements inspectiontaskbase.LogToTimelineMapperV2.
func (p *OtherTimelineMapper) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privategkemaster_contract.OtherGrouperTaskID.Ref()
}

// LogIngesterTask implements inspectiontaskbase.LogToTimelineMapperV2.
func (p *OtherTimelineMapper) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.LogIngesterTaskID.Ref()
}

// ProcessLogByGroup implements inspectiontaskbase.LogToTimelineMapperV2.
func (p *OtherTimelineMapper) ProcessLogByGroup(ctx context.Context, l *log.Log, _ struct{}) (*khifilev6.TimelineChangeSet, struct{}, error) {
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref())
	masterLogField := log.MustGetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
	commonLogField := log.MustGetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{})

	cs := khifilev6.NewTimelineChangeSet(l)
	writtenResourcePaths := map[uint32]struct{}{}
	for _, tPath := range masterLogField.ResourceTimelines(ctx, clusterIdentity.ClusterName) {
		cs.AddEvent(tPath)
		writtenResourcePaths[tPath.ID] = struct{}{}
	}

	finder := coretask.GetTaskResult(ctx, commonlogk8saudit_contract.ResourceUIDPatternFinderTaskID.Ref())
	resources := patternfinder.FindAllWithStarterRunes(commonLogField.Message, finder, false, p.uidPrefixTokenCandidates...)
	for _, resource := range resources {
		tPath := commonlogk8saudit_contract.MustResourceTimeline(ctx, clusterIdentity.ClusterName, resource.Value)
		if _, ok := writtenResourcePaths[tPath.ID]; ok {
			continue
		}
		cs.AddEvent(tPath)
		writtenResourcePaths[tPath.ID] = struct{}{}
	}
	return cs, struct{}{}, nil
}

var _ inspectiontaskbase.LogToTimelineMapperV2[struct{}] = (*OtherTimelineMapper)(nil)

var otherLogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTaskV2(
	privategkemaster_contract.OtherLogToTimelineMapperTaskID,
	&OtherTimelineMapper{
		uidPrefixTokenCandidates: []rune{
			'"', ' ', '\'', '=',
		},
	},
)
