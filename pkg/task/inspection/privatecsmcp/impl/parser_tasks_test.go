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
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	core_contract "github.com/GoogleCloudPlatform/khi/pkg/task/core/contract"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testchangeset"
)

func TestCSMCPTimelineMapper_ProcessLogByGroup(t *testing.T) {
	ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())
	builder := khifilev6.NewBuilder()
	ctx = khictx.WithValue(ctx, inspectioncore_contract.Builder, builder)

	tenantProjectPath := googlecloudcommon_contract.MustGCPProjectTimeline(ctx, "test-tenant")
	servicePath := privatecsmcp_contract.MustCloudRunServiceTimeline(ctx, tenantProjectPath, "test-service", "unknown")

	clusterPath := commonlogk8saudit_contract.MustK8sClusterTimeline(ctx, "test-cluster")
	apiPath := commonlogk8saudit_contract.MustK8sAPIVersionTimeline(ctx, clusterPath, "core/v1")
	kindPath := commonlogk8saudit_contract.MustK8sKindTimeline(ctx, apiPath, "pod")
	nsPath := commonlogk8saudit_contract.MustK8sNamespaceTimeline(ctx, kindPath, "default")
	podPath := commonlogk8saudit_contract.MustK8sNamespacedResourceTimeline(ctx, nsPath, "test-pod")
	csmcpPodPath := privatecsmcp_contract.MustCSMCPPodLogTimeline(ctx, podPath)
	connPath := privatecsmcp_contract.MustCSMCPConnectionTimeline(ctx, podPath, "1")

	testCases := []struct {
		name          string
		logMessage    string
		prevGroupData *csmcpTimelineState
		assert        func(*testing.T, *khifilev6.TimelineChangeSet)
	}{
		{
			name:       "new connection log",
			logMessage: "ADS: new connection for node:test-pod.default-1",
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(servicePath).
					HasEvent(csmcpPodPath).
					HasRevision(connPath, &khifilev6.StagingRevision{
						VerbType:     inspectioncore_contract.VerbUnknown,
						StateType:    privatecsmcp_contract.RevisionStateCSMCPConnectionConnected,
						ResourceBody: nil,
						Principal:    "csm-cp",
					})
			},
		},
		{
			name:       "terminated log with previous connected state",
			logMessage: "ADS: test-pod.default-1 terminated",
			prevGroupData: &csmcpTimelineState{
				ConnectedConns: map[string]bool{"1": true},
			},
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(servicePath).
					HasEvent(csmcpPodPath).
					HasRevision(connPath, &khifilev6.StagingRevision{
						VerbType:     inspectioncore_contract.VerbUnknown,
						StateType:    privatecsmcp_contract.RevisionStateCSMCPConnectionTerminated,
						ResourceBody: nil,
						Principal:    "csm-cp",
					})
			},
		},
		{
			name:       "terminated log without previous connected state",
			logMessage: "ADS: test-pod.default-1 terminated",
			prevGroupData: &csmcpTimelineState{
				ConnectedConns: map[string]bool{},
			},
			assert: func(t *testing.T, cs *khifilev6.TimelineChangeSet) {
				testchangeset.AssertTimeline(t, cs).
					HasEvent(servicePath).
					HasEvent(csmcpPodPath).
					HasRevision(connPath, &khifilev6.StagingRevision{
						ChangedTime:  time.Unix(0, 0),
						VerbType:     inspectioncore_contract.VerbUnknown,
						StateType:    privatecsmcp_contract.RevisionStateCSMCPConnectionConnectedLogNotFound,
						ResourceBody: nil,
						Principal:    "csm-cp",
					}).
					HasRevision(connPath, &khifilev6.StagingRevision{
						VerbType:     inspectioncore_contract.VerbUnknown,
						StateType:    privatecsmcp_contract.RevisionStateCSMCPConnectionTerminated,
						ResourceBody: nil,
						Principal:    "csm-cp",
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Re-use outer ctx so that paths match builder
			tm := khictx.MustGetValue(ctx, core_contract.TaskResultMapContextKey)
			typedmap.Set(tm, typedmap.NewTypedKey[string](privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref().ReferenceIDString()), "test-tenant")
			typedmap.Set(tm, typedmap.NewTypedKey[string](privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID.Ref().ReferenceIDString()), "test-service")
			typedmap.Set(tm, typedmap.NewTypedKey[googlecloudk8scommon_contract.GoogleCloudClusterIdentity](googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref().ReferenceIDString()), googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
				ClusterName: "test-cluster",
				Location:    "test-location",
			})

			fs := &privatecsmcp_contract.CSMCPFieldSet{
				InstanceID: "unknown",
				Message:    tc.logMessage,
				Pods: []privatecsmcp_contract.PodIdentifier{
					{
						Name:         "test-pod",
						Namespace:    "default",
						ConnectionID: "1",
					},
				},
			}
			logObj := log.NewLogWithFieldSetsForTest(fs, &log.CommonFieldSet{Timestamp: time.Now()})

			mapper := &csmcpTimelineMapper{}
			cs, _, err := mapper.ProcessLogByGroup(ctx, logObj, tc.prevGroupData)
			if err != nil {
				t.Fatalf("ProcessLogByGroup() failed: %v", err)
			}

			if tc.assert != nil {
				tc.assert(t, cs)
			}
		})
	}
}

func TestCSMCPTimelineMapper_PreProcessLogByGroup(t *testing.T) {
	testCases := []struct {
		name       string
		logMessage string
		wantConns  map[string]bool
	}{
		{
			name:       "connected log",
			logMessage: "ADS: new connection for node:test-pod.default-1",
			wantConns: map[string]bool{
				"default/test-pod/1": true,
			},
		},
		{
			name:       "terminated log",
			logMessage: "ADS: test-pod.default-1 terminated",
			wantConns:  map[string]bool{},
		},
		{
			name:       "unrelated log",
			logMessage: "Some other message",
			wantConns:  map[string]bool{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := &privatecsmcp_contract.CSMCPFieldSet{
				InstanceID: "unknown",
				Message:    tc.logMessage,
				Pods: []privatecsmcp_contract.PodIdentifier{
					{
						Name:         "test-pod",
						Namespace:    "default",
						ConnectionID: "1",
					},
				},
			}
			logObj := log.NewLogWithFieldSetsForTest(fs)

			mapper := &csmcpTimelineMapper{}
			state, err := mapper.PreProcessLogByGroup(t.Context(), 0, logObj, nil)
			if err != nil {
				t.Fatalf("PreProcessLogByGroup() failed: %v", err)
			}

			for id, expected := range tc.wantConns {
				if state.ConnectedConns[id] != expected {
					t.Errorf("PreProcessLogByGroup() returned connection state %v for %s, want %v", state.ConnectedConns[id], id, expected)
				}
			}
			for id, got := range state.ConnectedConns {
				if _, ok := tc.wantConns[id]; !ok && got {
					t.Errorf("PreProcessLogByGroup() unexpectedly returned true for %s", id)
				}
			}
		})
	}
}
