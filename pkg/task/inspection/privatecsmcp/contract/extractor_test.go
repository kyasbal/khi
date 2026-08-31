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

package privatecsmcp_contract

import (
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
	"github.com/google/go-cmp/cmp"
)

func TestExtractCSMCP(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    CSMCPFieldSet
		wantErr bool
	}{
		{
			name:  "no jsonPayload.message",
			input: `{}`,
			want: CSMCPFieldSet{
				Timestamp:  nil,
				Message:    "",
				Pods:       nil,
				InstanceID: "unknown",
			},
		},
		{
			name: "message without pods",
			input: `jsonPayload:
  message: "just some log"`,
			want: CSMCPFieldSet{
				Timestamp:  nil,
				Message:    "just some log",
				Pods:       nil,
				InstanceID: "unknown",
			},
		},
		{
			name: "message with pod only",
			input: `jsonPayload:
  message: "some text my-pod-1.default some other text"`,
			want: CSMCPFieldSet{
				Timestamp: nil,
				Message:   "some text my-pod-1.default some other text",
				Pods: []PodIdentifier{
					{Name: "my-pod-1", Namespace: "default"},
				},
				InstanceID: "unknown",
			},
		},
		{
			name: "message with node and pod",
			input: `jsonPayload:
  message: "node:my-pod-2.kube-system another-pod.default"`,
			want: CSMCPFieldSet{
				Timestamp: nil,
				Message:   "node:my-pod-2.kube-system another-pod.default",
				Pods: []PodIdentifier{
					{Name: "my-pod-2", Namespace: "kube-system"},
					{Name: "another-pod", Namespace: "default"},
				},
				InstanceID: "unknown",
			},
		},
		{
			name: "message with invalid pod formats",
			input: `jsonPayload:
  message: "notapod. node:invalid. pod.name.too.long node:my-pod.namespace text-node:pod.ns 100.200ms"`,
			// node:my-pod.namespace will match. 100.200ms will be ignored.
			want: CSMCPFieldSet{
				Timestamp: nil,
				Message:   "notapod. node:invalid. pod.name.too.long node:my-pod.namespace text-node:pod.ns 100.200ms",
				Pods: []PodIdentifier{
					{Name: "my-pod", Namespace: "namespace"},
				},
				InstanceID: "unknown",
			},
		},
		{
			name: "message with connection id in namespace",
			input: `jsonPayload:
  message: "another-pod.default-12345 pod.my-ns-456-789 node:p3.n3-12"`,
			want: CSMCPFieldSet{
				Timestamp: nil,
				Message:   "another-pod.default-12345 pod.my-ns-456-789 node:p3.n3-12",
				Pods: []PodIdentifier{
					{Name: "another-pod", Namespace: "default", ConnectionID: "12345"},
					{Name: "pod", Namespace: "my-ns-456", ConnectionID: "789"},
					{Name: "p3", Namespace: "n3", ConnectionID: "12"},
				},
				InstanceID: "unknown",
			},
		},
		{
			name: "with timestamp",
			input: `jsonPayload:
  message: "msg"
  timestamp: "2026-07-03T15:46:12Z"`,
			want: CSMCPFieldSet{
				Timestamp: func() *time.Time {
					t, _ := time.Parse(time.RFC3339, "2026-07-03T15:46:12Z")
					return &t
				}(),
				Message:    "msg",
				Pods:       nil,
				InstanceID: "unknown",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, err := structured.FromYAML(tc.input)
			if err != nil {
				t.Fatalf("failed to parse yaml: %v", err)
			}
			reader := structured.NewNodeReader(node)
			got, err := ExtractCSMCP(reader)
			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractCSMCP() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ExtractCSMCP() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("mock node returns mock values", func(t *testing.T) {
		mockFS := CSMCPFieldSet{
			Message:    "mock msg",
			InstanceID: "mock-inst",
		}
		reader := structured.NewNodeReader(structured.NewMockNode(mockFS))
		got, err := ExtractCSMCP(reader)
		if err != nil {
			t.Fatalf("ExtractCSMCP() error = %v", err)
		}
		if diff := cmp.Diff(mockFS, got); diff != "" {
			t.Errorf("ExtractCSMCP() mismatch (-want +got):\n%s", diff)
		}
	})
}
