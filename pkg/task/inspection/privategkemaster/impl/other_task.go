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
	"github.com/GoogleCloudPlatform/khi/pkg/model/enum"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8sauditv2_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8sauditv2/contract"
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

type otherLogToTimelineMapperTaskSetting struct {
	uidPrefixTokenCandidates []rune
}

func (p *otherLogToTimelineMapperTaskSetting) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privategkemaster_contract.OtherGrouperTaskID.Ref()
}

func (p *otherLogToTimelineMapperTaskSetting) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.LogIngesterTaskID.Ref()
}

func (p *otherLogToTimelineMapperTaskSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		commonlogk8sauditv2_contract.ResourceUIDPatternFinderTaskID.Ref(),
		googlecloudk8scommon_contract.ClusterIndentityTaskID.Ref(),
	}
}

func (p *otherLogToTimelineMapperTaskSetting) ProcessLogByGroup(ctx context.Context, l *log.Log, cs *history.ChangeSet, builder *history.Builder, prevGroupData struct{}) (struct{}, error) {
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIndentityTaskID.Ref())
	masterLogField := log.MustGetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
	commonLogField := log.MustGetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{})

	for _, path := range masterLogField.ResourcePaths(clusterIdentity.ClusterName) {
		cs.AddEvent(path)
	}

	summary := commonLogField.Message
	if masterLogField.StructuredBody != nil {
		msg, err := masterLogField.StructuredBody.MainMessage()
		if err == nil {
			summary = msg
		}
	}

	cs.SetLogSummary(summary)

	finder := coretask.GetTaskResult(ctx, commonlogk8sauditv2_contract.ResourceUIDPatternFinderTaskID.Ref())
	resources := patternfinder.FindAllWithStarterRunes(commonLogField.Message, finder, false, p.uidPrefixTokenCandidates...)
	writtenResourcePaths := map[string]struct{}{}
	for _, resource := range resources {
		path := resource.Value.ResourcePathString()
		if _, ok := writtenResourcePaths[path]; ok {
			continue
		}
		cs.AddEvent(resourcepath.ResourcePath{
			Path:               path,
			ParentRelationship: enum.RelationshipChild,
		})
		writtenResourcePaths[path] = struct{}{}
	}
	return struct{}{}, nil
}

var otherLogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask(
	privategkemaster_contract.OtherLogToTimelineMapperTaskID,
	&otherLogToTimelineMapperTaskSetting{
		uidPrefixTokenCandidates: []rune{
			'"', ' ', '\'', '=',
		},
	},
)
