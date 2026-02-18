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
	"github.com/GoogleCloudPlatform/khi/pkg/model/enum"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8sauditv2_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8sauditv2/contract"
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

// KubeletLogLogToTimelineMapperTask maps kubelet logs to the timeline.
var KubeletLogLogToTimelineMapperTask = inspectiontaskbase.NewLogToTimelineMapperTask[struct{}](
	privategkemaster_contract.KubeletLogLogToTimelineMapperTaskID,
	&kubeletNodeLogLogToTimelineMapperSetting{},
)

type kubeletNodeLogLogToTimelineMapperSetting struct{}

// Dependencies implements inspectiontaskbase.LogToTimelineMapper.
func (k *kubeletNodeLogLogToTimelineMapperSetting) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{
		privategkemaster_contract.PodSandboxIDDiscoveryTaskID.Ref(),
		commonlogk8sauditv2_contract.ContainerIDPatternFinderTaskID.Ref(),
		commonlogk8sauditv2_contract.ResourceUIDPatternFinderTaskID.Ref(),
		googlecloudk8scommon_contract.ClusterIndentityTaskID.Ref(),
	}
}

// GroupedLogTask implements inspectiontaskbase.LogToTimelineMapper.
func (k *kubeletNodeLogLogToTimelineMapperSetting) GroupedLogTask() taskid.TaskReference[inspectiontaskbase.LogGroupMap] {
	return privategkemaster_contract.KubeletLogGroupTaskID.Ref()
}

// LogIngesterTask implements inspectiontaskbase.LogToTimelineMapper.
func (k *kubeletNodeLogLogToTimelineMapperSetting) LogIngesterTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.LogIngesterTaskID.Ref()
}

// ProcessLogByGroup implements inspectiontaskbase.LogToTimelineMapper.
func (k *kubeletNodeLogLogToTimelineMapperSetting) ProcessLogByGroup(ctx context.Context, l *log.Log, cs *history.ChangeSet, builder *history.Builder, prevGroupData struct{}) (struct{}, error) {
	masterFieldSet := log.MustGetFieldSet(l, &privategkemaster_contract.GKEMasterLogFieldSet{})
	containerIDPatternFinder := coretask.GetTaskResult(ctx, commonlogk8sauditv2_contract.ContainerIDPatternFinderTaskID.Ref())
	podIDFinder := coretask.GetTaskResult(ctx, privategkemaster_contract.PodSandboxIDDiscoveryTaskID.Ref())
	resourceUIDPatternFinder := coretask.GetTaskResult(ctx, commonlogk8sauditv2_contract.ResourceUIDPatternFinderTaskID.Ref())
	clusterIdentity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIndentityTaskID.Ref())

	for _, path := range masterFieldSet.ResourcePaths(clusterIdentity.ClusterName) {
		cs.AddEvent(path)
	}

	original := masterFieldSet.StructuredBody.Raw()

	severity, err := masterFieldSet.StructuredBody.Severity()
	if err == nil {
		cs.SetLogSeverity(severity)
	}

	foundPods := map[string]struct{}{}
	summaryReplaceMap := map[string]string{}
	podFindResults := patternfinder.FindAllWithStarterRunes(original, podIDFinder, false, '"')

	for _, result := range podFindResults {
		cs.AddEvent(result.Value.ResourcePath())
		summaryReplaceMap[result.Value.PodSandboxID] = toReadablePodSandboxName(result.Value.PodNamespace, result.Value.PodName)
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
		cs.AddEvent(result.Value.ResourcePath(pod.PodNamespace, pod.PodName))
		summaryReplaceMap[result.Value.ContainerID] = toReadableContainerName(pod.PodNamespace, pod.PodName, result.Value.ContainerName)
	}

	resourceFindResults := patternfinder.FindAllWithStarterRunes(original, resourceUIDPatternFinder, false, '"')
	for _, result := range resourceFindResults {
		res := result.Value
		if res.APIVersion == "core/v1" && res.Kind == "pod" {
			if _, ok := foundPods[fmt.Sprintf("%s/%s", res.Namespace, res.Name)]; ok {
				continue
			}
		}
		cs.AddEvent(resourcepath.ResourcePath{
			Path:               res.ResourcePathString(),
			ParentRelationship: enum.RelationshipChild,
		})
		uid, err := result.GetMatchedString(original)
		if err != nil {
			continue
		}
		summaryReplaceMap[uid] = toReadableResourceName(result.Value.APIVersion, result.Value.Kind, result.Value.Namespace, result.Value.Name)
	}

	// Kubelet specific severity adjustments
	klogExitCode, err := masterFieldSet.StructuredBody.StringField("exitCode")
	if err == nil && klogExitCode != "" && klogExitCode != "0" {
		if klogExitCode == "137" {
			cs.SetLogSeverity(enum.SeverityError)
		} else {
			cs.SetLogSeverity(enum.SeverityWarning)
		}
	}
	summary, err := parseDefaultSummary(masterFieldSet.StructuredBody)
	if err != nil {
		summary = original
	}
	for k, v := range summaryReplaceMap {
		i := strings.Index(summary, k)
		if i == -1 {
			summary = fmt.Sprintf("%s %s", summary, v)
		} else {
			summary = strings.ReplaceAll(summary, k, v)
		}
	}

	// Kubelet specific resource bindings
	podNameWithNamespace, err := masterFieldSet.StructuredBody.StringField("pod")
	if err == nil && podNameWithNamespace != "" {
		podNamespace, podName, err := slashSplittedPodNameToNamespaceAndName(podNameWithNamespace)
		if err == nil {
			containerName, err := masterFieldSet.StructuredBody.StringField("containerName")
			if err == nil && containerName != "" {
				cs.AddEvent(resourcepath.Container(podNamespace, podName, containerName))
				cs.SetLogSummary(fmt.Sprintf("%s %s", summary, toReadableContainerName(podNamespace, podName, containerName)))
			} else {
				cs.AddEvent(resourcepath.Pod(podNamespace, podName))
				cs.SetLogSummary(fmt.Sprintf("%s %s", summary, toReadablePodSandboxName(podNamespace, podName)))
			}
		} else {
			cs.SetLogSummary(summary)
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
					cs.AddEvent(resourcepath.Pod(podNamespace, podName))
					summary = fmt.Sprintf("%s %s", summary, toReadablePodSandboxName(podNamespace, podName))
				}
			}
		}
		cs.SetLogSummary(summary)
	}

	return struct{}{}, nil
}
