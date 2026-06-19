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
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/patternfinder"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
)

// KubeletLogFilterTask filters logs for kubelet.
var KubeletLogFilterTask = inspectiontaskbase.NewLogFilterTask(
	privategkemaster_contract.KubeletLogFilterTaskID,
	privategkemaster_contract.CommonFieldSetReaderTaskID.Ref(),
	func(ctx context.Context, l *log.Log) bool {
		componentFieldSet, err := log.GetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
		if err != nil {
			return false
		}
		return componentFieldSet.PrivateGKEMasterParserType() == privategkemaster_contract.PrivateGKEMasterParserTypeKubelet
	},
)

// KubeletLogGroupTask groups kubelet logs.
var KubeletLogGroupTask = inspectiontaskbase.NewLogGrouperTask(
	privategkemaster_contract.KubeletLogGroupTaskID,
	privategkemaster_contract.KubeletLogFilterTaskID.Ref(),
	func(ctx context.Context, l *log.Log) string {
		masterLogField := log.MustGetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
		return fmt.Sprintf("%s#%s", masterLogField.HostName, masterLogField.PodID)
	},
)

// KubeletTimelineMapper maps kubelet logs to timeline paths.
type KubeletTimelineMapper struct {
	inspectiontaskbase.StatelessMapperBase
}

// Dependencies implements inspectiontaskbase.LogToTimelineMapperV2.
func (k *KubeletTimelineMapper) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privategkemaster_contract.PodSandboxIDDiscoveryTaskID.Ref(),
		commonlogk8saudit_contract.ContainerIDPatternFinderTaskID.Ref(),
		commonlogk8saudit_contract.ResourceUIDPatternFinderTaskID.Ref(),
		googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(),
	}
}

// GroupedLogTask implements inspectiontaskbase.LogToTimelineMapperV2.
func (k *KubeletTimelineMapper) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privategkemaster_contract.KubeletLogGroupTaskID.Ref()
}

// LogIngesterTask implements inspectiontaskbase.LogToTimelineMapperV2.
func (k *KubeletTimelineMapper) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.LogIngesterTaskID.Ref()
}

// ProcessLogByGroup implements inspectiontaskbase.LogToTimelineMapperV2.
func (k *KubeletTimelineMapper) ProcessLogByGroup(ctx context.Context, l *log.Log, _ struct{}) (*khifilev6.TimelineChangeSet, struct{}, error) {
	masterFieldSet := log.MustGetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
	containerIDPatternFinder := coretask.GetTaskResult(ctx, commonlogk8saudit_contract.ContainerIDPatternFinderTaskID.Ref())
	podIDFinder := coretask.GetTaskResult(ctx, privategkemaster_contract.PodSandboxIDDiscoveryTaskID.Ref())
	resourceUIDPatternFinder := coretask.GetTaskResult(ctx, commonlogk8saudit_contract.ResourceUIDPatternFinderTaskID.Ref())
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref())

	cs := khifilev6.NewTimelineChangeSet(l)

	for _, tPath := range masterFieldSet.ResourceTimelines(ctx, clusterIdentity.ClusterName) {
		cs.AddEvent(tPath)
	}

	original := masterFieldSet.StructuredBody.Raw()

	foundPods := map[string]struct{}{}
	podFindResults := patternfinder.FindAllWithStarterRunes(original, podIDFinder, false, '"')

	for _, result := range podFindResults {
		podTimelinePath := mustK8sPodTimeline(ctx, clusterIdentity.ClusterName, result.Value.PodNamespace, result.Value.PodName)
		cs.AddEvent(podTimelinePath)
		foundPods[fmt.Sprintf("%s/%s", result.Value.PodNamespace, result.Value.PodName)] = struct{}{}
	}

	containerFindResults := patternfinder.FindAllWithStarterRunes(original, containerIDPatternFinder, false, '"')
	for _, result := range containerFindResults {
		podSandboxID := result.Value.PodSandboxID
		foundPod := patternfinder.FindAllWithStarterRunes(podSandboxID, podIDFinder, true)
		if len(foundPod) == 0 {
			continue
		}
		pod := foundPod[0].Value
		podTimelinePath := mustK8sPodTimeline(ctx, clusterIdentity.ClusterName, pod.PodNamespace, pod.PodName)
		containerTimelinePath := commonlogk8saudit_contract.MustK8sContainerTimeline(ctx, podTimelinePath, result.Value.ContainerName)
		cs.AddEvent(containerTimelinePath)
	}

	resourceFindResults := patternfinder.FindAllWithStarterRunes(original, resourceUIDPatternFinder, false, '"')
	for _, result := range resourceFindResults {
		res := result.Value
		if res.APIVersion == "core/v1" && res.Kind == "pod" {
			if _, ok := foundPods[fmt.Sprintf("%s/%s", res.Namespace, res.Name)]; ok {
				continue
			}
		}
		resTimelinePath := commonlogk8saudit_contract.MustResourceTimeline(ctx, clusterIdentity.ClusterName, res)
		cs.AddEvent(resTimelinePath)
	}

	// Kubelet specific resource bindings
	podNameWithNamespace, err := masterFieldSet.StructuredBody.StringField("pod")
	if err == nil && podNameWithNamespace != "" {
		podNamespace, podName, err := slashSplittedPodNameToNamespaceAndName(podNameWithNamespace)
		if err == nil {
			podTimelinePath := mustK8sPodTimeline(ctx, clusterIdentity.ClusterName, podNamespace, podName)
			containerName, err := masterFieldSet.StructuredBody.StringField("containerName")
			if err == nil && containerName != "" {
				containerTimelinePath := commonlogk8saudit_contract.MustK8sContainerTimeline(ctx, podTimelinePath, containerName)
				cs.AddEvent(containerTimelinePath)
			} else {
				cs.AddEvent(podTimelinePath)
			}
		}
	} else {
		podNames, err := masterFieldSet.StructuredBody.StringField("pods")
		if err == nil && podNames != "" {
			podNames = strings.Trim(podNames, "[]")
			podNamesSplitted := strings.Split(podNames, ",")
			for _, podNamespaceAndNameWithSlash := range podNamesSplitted {
				podNamespaceAndNameWithSlash = strings.Trim(podNamespaceAndNameWithSlash, `"`)
				podNamespace, podName, err := slashSplittedPodNameToNamespaceAndName(podNamespaceAndNameWithSlash)
				if err == nil {
					podTimelinePath := mustK8sPodTimeline(ctx, clusterIdentity.ClusterName, podNamespace, podName)
					cs.AddEvent(podTimelinePath)
				}
			}
		}
	}

	return cs, struct{}{}, nil
}

func mustK8sPodTimeline(ctx context.Context, clusterName string, namespace string, podName string) *khifilev6.TimelinePath {
	clusterPath := commonlogk8saudit_contract.MustK8sClusterTimeline(ctx, clusterName)
	apiVersionPath := commonlogk8saudit_contract.MustK8sAPIVersionTimeline(ctx, clusterPath, "core/v1")
	kindPath := commonlogk8saudit_contract.MustK8sKindTimeline(ctx, apiVersionPath, "pod")
	namespacePath := commonlogk8saudit_contract.MustK8sNamespaceTimeline(ctx, kindPath, namespace)
	return commonlogk8saudit_contract.MustK8sNamespacedResourceTimeline(ctx, namespacePath, podName)
}

var _ inspectiontaskbase.LogToTimelineMapperV2[struct{}] = (*KubeletTimelineMapper)(nil)

// KubeletLogLogToTimelineMapperTask maps kubelet logs to the timeline.
var KubeletLogLogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTaskV2(
	privategkemaster_contract.KubeletLogLogToTimelineMapperTaskID,
	&KubeletTimelineMapper{},
)
