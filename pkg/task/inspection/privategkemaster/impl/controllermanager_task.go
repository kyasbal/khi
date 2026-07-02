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

func kindToKLogFieldPair(apiVersion string, kind string, klogField string, isNamespaced bool) *googlecloudlogk8scontrolplane_contract.KindToKLogFieldPairData {
	return &googlecloudlogk8scontrolplane_contract.KindToKLogFieldPairData{
		APIVersion:   apiVersion,
		KindName:     kind,
		KLogField:    klogField,
		IsNamespaced: isNamespaced,
	}
}

var controllerManagerLogFilterTask = inspectiontaskbase.NewLogFilterTask(
	privategkemaster_contract.ControllerManagerLogFilterTaskID,
	privategkemaster_contract.CommonFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, l *log.Log) bool {
		componentFieldSet, err := log.GetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
		if err != nil {
			return false
		}
		return componentFieldSet.PrivateGKEMasterParserType() == privategkemaster_contract.PrivateGKEMasterParserTypeControllerManager
	},
)

var controllerManagerLogFieldSetReaderTask = inspectiontaskbase.NewFieldSetReadTask(privategkemaster_contract.ControllerManagerLogFieldSetReaderTaskID,
	privategkemaster_contract.ControllerManagerLogFilterTaskID.Ref(),
	[]log.FieldSetReader{
		&googlecloudlogk8scontrolplane_contract.K8sControllerManagerComponentFieldSetReader{
			WellKnownSourceLocationToControllerMap: map[string]string{
				"namespace_controller.go":      "namespace-controller",
				"resource_quota_controller.go": "resourcequota-controller",
				"requestheader_controller.go":  "requestheader-controller",
				"pv_protection_controller.go":  "persistentvolume-protection-controller",
			},
			WellKnownKindToKLogFieldPairs: []*googlecloudlogk8scontrolplane_contract.KindToKLogFieldPairData{
				kindToKLogFieldPair("apps/v1", "deployment", "deployment", true),
				kindToKLogFieldPair("apps/v1", "replicaset", "replicaSet", true),
				kindToKLogFieldPair("apps/v1", "statefulset", "statefulSet", true),
				kindToKLogFieldPair("apps/v1", "daemonset", "daemonSet", true),
				kindToKLogFieldPair("batch/v1", "cronjob", "cronjob", true),
				kindToKLogFieldPair("batch/v1", "job", "job", true),
				kindToKLogFieldPair("policy/v1", "poddisruptionbudget", "podDisruptionBudget", true),
				kindToKLogFieldPair("certificates.k8s.io/v1", "certificatesigningrequest", "csr", false),
				kindToKLogFieldPair("core/v1", "persistentvolumeclaim", "PVC", true),
				kindToKLogFieldPair("core/v1", "persistentvolume", "volumeName", false),
				kindToKLogFieldPair("core/v1", "service", "service", true),
				kindToKLogFieldPair("core/v1", "node", "node", false),
				kindToKLogFieldPair("core/v1", "pod", "pod", true),
				kindToKLogFieldPair("core/v1", "namespace", "namespace", false),
			},
			KLogParser: logutil.NewKLogTextParser(false),
		},
	},
)

var controllerManagerGrouperTask = inspectiontaskbase.NewLogGrouperTask(
	privategkemaster_contract.ControllerManagerGrouperTaskID,
	privategkemaster_contract.ControllerManagerLogFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, log *log.Log) string {
		return "" // No grouping needed
	},
)

var controllerManagerLogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask[struct{}](
	privategkemaster_contract.ControllerManagerLogToTimelineMapperTaskID,
	&ControllerManagerTimelineMapper{
		uidPrefixTokenCandidates: []rune{
			'"', ' ', '\'', '=',
		},
	},
)

// ControllerManagerTimelineMapper maps controller manager logs to timeline paths.
type ControllerManagerTimelineMapper struct {
	inspectiontaskbase.StatelessMapperBase
	uidPrefixTokenCandidates []rune
}

// Dependencies implements inspectiontaskbase.LogToTimelineMapper.
func (m *ControllerManagerTimelineMapper) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		commonlogk8saudit_contract.ResourceUIDPatternFinderTaskID.Ref(),
		googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(),
	}
}

// GroupedLogTask implements inspectiontaskbase.LogToTimelineMapper.
func (m *ControllerManagerTimelineMapper) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privategkemaster_contract.ControllerManagerGrouperTaskID.Ref()
}

// LogIngesterTask implements inspectiontaskbase.LogToTimelineMapper.
func (m *ControllerManagerTimelineMapper) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.LogIngesterTaskID.Ref()
}

// ProcessLogByGroup implements inspectiontaskbase.LogToTimelineMapper.
func (m *ControllerManagerTimelineMapper) ProcessLogByGroup(ctx context.Context, l *log.Log, _ struct{}) (*khifilev6.TimelineChangeSet, struct{}, error) {
	finder := coretask.GetTaskResult(ctx, commonlogk8saudit_contract.ResourceUIDPatternFinderTaskID.Ref())
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref())
	masterFieldSet, err := log.GetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
	if err != nil {
		return nil, struct{}{}, err
	}
	commonMainMessage, err := log.GetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{})
	if err != nil {
		return nil, struct{}{}, err
	}
	controllerManagerFieldSet, err := log.GetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sControllerManagerComponentFieldSet{})
	if err != nil {
		return nil, struct{}{}, err
	}

	cs := khifilev6.NewTimelineChangeSet(l)
	resources := patternfinder.FindAllWithStarterRunes(commonMainMessage.Message, finder, false, m.uidPrefixTokenCandidates...)
	writtenResourcePaths := map[uint32]struct{}{}

	for _, tPath := range masterFieldSet.ResourceTimelines(ctx, clusterIdentity.ClusterName) {
		cs.AddEvent(tPath)
		writtenResourcePaths[tPath.ID] = struct{}{}
	}

	for _, tPath := range controllerManagerFieldSet.AssociatedResourceTimelines(ctx, clusterIdentity.ClusterName) {
		cs.AddEvent(tPath)
		writtenResourcePaths[tPath.ID] = struct{}{}
	}
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

var _ inspectiontaskbase.LogToTimelineMapper[struct{}] = (*ControllerManagerTimelineMapper)(nil)
