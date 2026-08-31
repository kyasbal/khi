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
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
	googlecloudlogk8scontrolplane_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8scontrolplane/contract"
	"github.com/google/go-cmp/cmp"
)

func TestGKEMasterLogFieldSet_PrivateGKEMasterParserType(t *testing.T) {
	testCases := []struct {
		podID    string
		expected PrivateGKEMasterParserType
	}{
		{
			podID:    "kube-scheduler",
			expected: PrivateGKEMasterParserTypeScheduler,
		},
		{
			podID:    "kube-controller-manager",
			expected: PrivateGKEMasterParserTypeControllerManager,
		},
		{
			podID:    "kubelet",
			expected: PrivateGKEMasterParserTypeKubelet,
		},
		{
			podID:    "containerd",
			expected: PrivateGKEMasterParserTypeContainerRuntime,
		},
		{
			podID:    "kube-apiserver",
			expected: PrivateGKEMasterParserTypeOther,
		},
		{
			podID:    "random-component",
			expected: PrivateGKEMasterParserTypeOther,
		},
		{
			podID:    "",
			expected: PrivateGKEMasterParserTypeOther,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.podID, func(t *testing.T) {
			fs := &GKEMasterLogFieldSet{
				ComponentName: tc.podID,
			}
			got := fs.PrivateGKEMasterParserType()
			if got != tc.expected {
				t.Errorf("PrivateGKEMasterParserType() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestExtractGKEMasterLog(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		wantProject string
		wantHost    string
		wantComp    string
		wantPodID   string
		wantNSID    string
		wantCont    string
	}{
		{
			desc: "container resource log",
			input: `
logName: "projects/test-project/logs/kube-controller-manager"
resource:
  type: "container"
  labels:
    project_id: "test-project"
    pod_id: "pod-1"
    namespace_id: "kube-system"
    container_name: "kube-controller-manager"
labels:
  "compute.googleapis.com/resource_name": "master-node-0"
textPayload: "I0831 12:00:00.000000 1 leaderelection.go:248] successfully acquired lease"
`,
			wantProject: "test-project",
			wantHost:    "master-node-0",
			wantComp:    "kube-controller-manager",
			wantPodID:   "pod-1",
			wantNSID:    "kube-system",
			wantCont:    "kube-controller-manager",
		},
		{
			desc: "non-container resource log with jsonPayload",
			input: `
logName: "projects/test-project/logs/kubelet"
resource:
  type: "gke_node"
  labels:
    project_id: "test-project"
labels:
  "compute.googleapis.com/resource_name": "master-node-1"
jsonPayload:
  message: "I0831 12:00:00.000000 1 kubelet.go:100] Starting kubelet"
`,
			wantProject: "test-project",
			wantHost:    "master-node-1",
			wantComp:    "kubelet",
			wantPodID:   "",
			wantNSID:    "",
			wantCont:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			node, err := structured.FromYAML(tc.input)
			if err != nil {
				t.Fatalf("failed to parse YAML test input: %v", err)
			}
			reader := structured.NewNodeReader(node)
			got, err := ExtractGKEMasterLog(reader)
			if err != nil {
				t.Fatalf("ExtractGKEMasterLog() unexpected error: %v", err)
			}
			if got.ProjectID != tc.wantProject {
				t.Errorf("ProjectID = %q, want %q", got.ProjectID, tc.wantProject)
			}
			if got.HostName != tc.wantHost {
				t.Errorf("HostName = %q, want %q", got.HostName, tc.wantHost)
			}
			if got.ComponentName != tc.wantComp {
				t.Errorf("ComponentName = %q, want %q", got.ComponentName, tc.wantComp)
			}
			if got.PodID != tc.wantPodID {
				t.Errorf("PodID = %q, want %q", got.PodID, tc.wantPodID)
			}
			if got.NamespaceID != tc.wantNSID {
				t.Errorf("NamespaceID = %q, want %q", got.NamespaceID, tc.wantNSID)
			}
			if got.ContainerName != tc.wantCont {
				t.Errorf("ContainerName = %q, want %q", got.ContainerName, tc.wantCont)
			}
			if got.StructuredBody == nil {
				t.Errorf("StructuredBody is nil")
			}
		})
	}

	t.Run("mock node returns mock values", func(t *testing.T) {
		mockFS := GKEMasterLogFieldSet{
			ProjectID:     "mock-proj",
			HostName:      "mock-host",
			ComponentName: "kubelet",
		}
		reader := structured.NewNodeReader(structured.NewMockNode(mockFS))
		got, err := ExtractGKEMasterLog(reader)
		if err != nil {
			t.Fatalf("ExtractGKEMasterLog() error = %v", err)
		}
		if diff := cmp.Diff(mockFS, got); diff != "" {
			t.Errorf("ExtractGKEMasterLog mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestExtractGKEMasterCommonMessage(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  string
	}{
		{
			desc: "from textPayload",
			input: `
textPayload: "sample text payload"
`,
			want: "sample text payload",
		},
		{
			desc: "from jsonPayload.message",
			input: `
jsonPayload:
  message: "sample json message"
`,
			want: "sample json message",
		},
		{
			desc: "from jsonPayload.MESSAGE",
			input: `
jsonPayload:
  MESSAGE: "sample json MESSAGE"
`,
			want: "sample json MESSAGE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			node, err := structured.FromYAML(tc.input)
			if err != nil {
				t.Fatalf("failed to parse YAML: %v", err)
			}
			reader := structured.NewNodeReader(node)
			got, err := ExtractGKEMasterCommonMessage(reader)
			if err != nil {
				t.Fatalf("ExtractGKEMasterCommonMessage() error = %v", err)
			}
			if got != tc.want {
				t.Errorf("ExtractGKEMasterCommonMessage() = %q, want %q", got, tc.want)
			}
		})
	}

	t.Run("mock node returns mock message", func(t *testing.T) {
		mockFS := googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{
			Message: "mock common message",
		}
		reader := structured.NewNodeReader(structured.NewMockNode(mockFS))
		got, err := ExtractGKEMasterCommonMessage(reader)
		if err != nil {
			t.Fatalf("ExtractGKEMasterCommonMessage() error = %v", err)
		}
		if got != mockFS.Message {
			t.Errorf("ExtractGKEMasterCommonMessage() = %q, want %q", got, mockFS.Message)
		}
	})
}
