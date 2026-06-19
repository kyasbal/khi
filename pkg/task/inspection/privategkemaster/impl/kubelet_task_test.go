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
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/patternfinder"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/logutil"
	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	googlecloudlogk8snode_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8snode/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testchangeset"
)

func TestKubeletLogLogToTimelineMapper(t *testing.T) {
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		desc                 string
		inputMessage         string
		inputNodeLogFieldSet *privategkemaster_contract.GKEMasterLogFieldSet
		inputPodIDInfo       map[string]*googlecloudlogk8snode_contract.PodSandboxIDInfo
		inputContainerIDInfo map[string]*commonlogk8saudit_contract.ContainerIdentity
		inputResourceUIDInfo map[string]*commonlogk8saudit_contract.ResourceIdentity
		assert               func(t *testing.T, ctx context.Context, cs *khifilev6.TimelineChangeSet)
	}{
		{
			desc:         "log with pod sandbox id",
			inputMessage: `I0929 08:30:43.794472    1949 generic.go:334] "Generic (PLEG): container finished" podID="4cba26fb-f074-44fe-9afa-5195e903c337" podID="6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1"`,
			inputNodeLogFieldSet: &privategkemaster_contract.GKEMasterLogFieldSet{
				ComponentName: "kubelet",
				HostName:      "node-1",
			},
			inputPodIDInfo: map[string]*googlecloudlogk8snode_contract.PodSandboxIDInfo{
				"6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1": {
					PodName:      "podname",
					PodNamespace: "kube-system",
					PodSandboxID: "6123c6aacf0c78dc38ec4f0ff72edd3cf04eb82ca0e3e7dddd3950ea9753bdf1",
				},
			},
			assert: func(t *testing.T, ctx context.Context, cs *khifilev6.TimelineChangeSet) {
				wantPodPath := mustK8sPodTimeline(ctx, "cluster", "kube-system", "podname")
				testchangeset.AssertTimeline(t, cs).
					HasEvent(wantPodPath)
			},
		},
		{
			desc:         "log with pod and container name from klog fields and exitCode",
			inputMessage: `I0929 08:30:43.794472    1949 generic.go:334] "Killing container" podID="4cba26fb-f074-44fe-9afa-5195e903c337" msg="Syncing pod" pod="kube-system/podname" containerName="containername" exitCode=137`,
			inputNodeLogFieldSet: &privategkemaster_contract.GKEMasterLogFieldSet{
				ComponentName: "kubelet",
				HostName:      "node-1",
			},
			assert: func(t *testing.T, ctx context.Context, cs *khifilev6.TimelineChangeSet) {
				wantPodPath := mustK8sPodTimeline(ctx, "cluster", "kube-system", "podname")
				wantContainerPath := commonlogk8saudit_contract.MustK8sContainerTimeline(ctx, wantPodPath, "containername")
				testchangeset.AssertTimeline(t, cs).
					HasEvent(wantContainerPath)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Mock the task results for dependencies
			podIDFinder := patternfinder.NewNaivePatternFinder[*googlecloudlogk8snode_contract.PodSandboxIDInfo]()
			if tc.inputPodIDInfo != nil {
				for k, v := range tc.inputPodIDInfo {
					podIDFinder.AddPattern(k, v)
				}
			}
			containerIDFinder := patternfinder.NewNaivePatternFinder[*commonlogk8saudit_contract.ContainerIdentity]()
			if tc.inputContainerIDInfo != nil {
				for k, v := range tc.inputContainerIDInfo {
					containerIDFinder.AddPattern(k, v)
				}
			}
			finder := patternfinder.NewNaivePatternFinder[*commonlogk8saudit_contract.ResourceIdentity]()
			if tc.inputResourceUIDInfo != nil {
				for k, v := range tc.inputResourceUIDInfo {
					finder.AddPattern(k, v)
				}
			}

			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())
			ctx = tasktest.WithTaskResult(ctx, privategkemaster_contract.PodSandboxIDDiscoveryTaskID.Ref(), podIDFinder)
			ctx = tasktest.WithTaskResult(ctx, commonlogk8saudit_contract.ContainerIDPatternFinderTaskID.Ref(), containerIDFinder)
			ctx = tasktest.WithTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(), googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
				ClusterName: "cluster",
				ProjectID:   "project",
				Location:    "location",
			})

			ctx = tasktest.WithTaskResult(ctx, commonlogk8saudit_contract.ResourceUIDPatternFinderTaskID.Ref(), finder)
			klogParser := logutil.NewMultiTextLogParser(
				logutil.NewKLogTextParser(true),
				&logutil.FallbackRawTextLogParser{},
			)
			message := klogParser.TryParse(tc.inputMessage)
			tc.inputNodeLogFieldSet.StructuredBody = message
			l := log.NewLogWithFieldSetsForTest(
				&log.CommonFieldSet{Timestamp: testTime},
				tc.inputNodeLogFieldSet,
			)
			mapper := &KubeletTimelineMapper{}
			cs, _, err := mapper.ProcessLogByGroup(ctx, l, struct{}{})
			if err != nil {
				t.Fatalf("ProcessLogByGroup() error = %v", err)
			}
			tc.assert(t, ctx, cs)
		})
	}
}
