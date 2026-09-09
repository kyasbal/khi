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

package inspectionmetadata

import (
	"fmt"
	"sync"
	"testing"

	pb "github.com/GoogleCloudPlatform/khi/pkg/generated/khifile/v6"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAIContextMetadata_BasicRecording(t *testing.T) {
	testCases := []struct {
		name         string
		operations   func(m *AIContextMetadata)
		wantSections []*AIContextSection
	}{
		{
			name: "records properties, sets, and markdown in order of priority",
			operations: func(m *AIContextMetadata) {
				m.SetPriority("Section B", 200)
				m.SetProperty("Section B", "Key1", "Value1")

				m.SetPriority("Section A", 100)
				m.AddToSetProperty("Section A", "Namespaces", "kube-system")
				m.AddToSetProperty("Section A", "Namespaces", "default")
				m.AddToSetProperty("Section A", "Namespaces", "default") // Duplicate
				m.AppendSummaryMarkdown("Section A", "Summary line 1")
				m.AppendSummaryMarkdown("Section A", "Summary line 2")
			},
			wantSections: []*AIContextSection{
				{
					Title:      "Section A",
					Priority:   100,
					Properties: map[string]string{},
					SetProperties: map[string][]string{
						"Namespaces": {"default", "kube-system"},
					},
					SummaryMarkdown: "Summary line 1\n\nSummary line 2",
				},
				{
					Title:    "Section B",
					Priority: 200,
					Properties: map[string]string{
						"Key1": "Value1",
					},
					SetProperties:   map[string][]string{},
					SummaryMarkdown: "",
				},
			},
		},
		{
			name: "ties in priority are broken by alphabetical title",
			operations: func(m *AIContextMetadata) {
				m.SetProperty("Zeta", "K", "V")
				m.SetProperty("Alpha", "K", "V")
			},
			wantSections: []*AIContextSection{
				{
					Title:           "Alpha",
					Priority:        DefaultSectionPriority,
					Properties:      map[string]string{"K": "V"},
					SetProperties:   map[string][]string{},
					SummaryMarkdown: "",
				},
				{
					Title:           "Zeta",
					Priority:        DefaultSectionPriority,
					Properties:      map[string]string{"K": "V"},
					SetProperties:   map[string][]string{},
					SummaryMarkdown: "",
				},
			},
		},
		{
			name: "ignores empty inputs and normalizes markdown newlines",
			operations: func(m *AIContextMetadata) {
				m.AddToSetProperty("Section A", "EmptyKey", "")
				m.AppendSummaryMarkdown("Section A", "")
				m.AppendSummaryMarkdown("Section A", "Line 1\n\n\n")
				m.AppendSummaryMarkdown("Section A", "\n\n\nLine 2")
			},
			wantSections: []*AIContextSection{
				{
					Title:           "Section A",
					Priority:        DefaultSectionPriority,
					Properties:      map[string]string{},
					SetProperties:   map[string][]string{},
					SummaryMarkdown: "Line 1\n\nLine 2",
				},
			},
		},
		{
			name: "overwrites scalar property on repeated key",
			operations: func(m *AIContextMetadata) {
				m.SetProperty("Section A", "Key", "Initial")
				m.SetProperty("Section A", "Key", "Updated")
			},
			wantSections: []*AIContextSection{
				{
					Title:           "Section A",
					Priority:        DefaultSectionPriority,
					Properties:      map[string]string{"Key": "Updated"},
					SetProperties:   map[string][]string{},
					SummaryMarkdown: "",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewAIContextMetadata()
			tc.operations(m)

			got := m.Sections()
			if diff := cmp.Diff(tc.wantSections, got); diff != "" {
				t.Errorf("Sections() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAIContextMetadata_ToSerializable(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(m *AIContextMetadata)
		wantProto *pb.AIContextMetadata
	}{
		{
			name:  "empty metadata serializes cleanly",
			setup: func(m *AIContextMetadata) {},
			wantProto: &pb.AIContextMetadata{
				Sections: []*pb.AIContextSection{},
			},
		},
		{
			name: "preserves default section priority through proto",
			setup: func(m *AIContextMetadata) {
				m.SetProperty("DefaultSec", "Version", "1.0")
			},
			wantProto: &pb.AIContextMetadata{
				Sections: []*pb.AIContextSection{
					{
						Title:           proto.String("DefaultSec"),
						Priority:        proto.Int32(DefaultSectionPriority),
						Properties:      map[string]string{"Version": "1.0"},
						SetProperties:   map[string]*pb.StringList{},
						SummaryMarkdown: proto.String(""),
					},
				},
			},
		},
		{
			name: "serializes custom sections cleanly",
			setup: func(m *AIContextMetadata) {
				m.SetPriority("K8s", 10)
				m.SetProperty("K8s", "Version", "v1.30")
				m.AddToSetProperty("K8s", "Nodes", "node-1")
				m.AddToSetProperty("K8s", "Nodes", "node-2")
				m.AppendSummaryMarkdown("K8s", "Everything healthy.")
			},
			wantProto: &pb.AIContextMetadata{
				Sections: []*pb.AIContextSection{
					{
						Title:    proto.String("K8s"),
						Priority: proto.Int32(10),
						Properties: map[string]string{
							"Version": "v1.30",
						},
						SetProperties: map[string]*pb.StringList{
							"Nodes": {
								Values: []string{"node-1", "node-2"},
							},
						},
						SummaryMarkdown: proto.String("Everything healthy."),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewAIContextMetadata()
			tc.setup(m)

			serializable := m.ToSerializable()
			pbMeta, ok := serializable.(*pb.AIContextMetadata)
			if !ok {
				t.Fatalf("ToSerializable() returned unexpected type %T", serializable)
			}

			if diff := cmp.Diff(tc.wantProto, pbMeta, protocmp.Transform()); diff != "" {
				t.Errorf("ToSerializable() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAIContextMetadata_ConcurrentAccess(t *testing.T) {
	testCases := []struct {
		name          string
		goroutines    int
		opsPerRoutine int
	}{
		{
			name:          "concurrent writes and reads without data races",
			goroutines:    16,
			opsPerRoutine: 50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewAIContextMetadata()
			var wg sync.WaitGroup

			for i := 0; i < tc.goroutines; i++ {
				wg.Add(1)
				workerID := i
				go func() {
					defer wg.Done()
					for j := 0; j < tc.opsPerRoutine; j++ {
						secName := fmt.Sprintf("Section-%d", j%4)
						m.SetProperty(secName, fmt.Sprintf("Key-%d", workerID), fmt.Sprintf("Val-%d", j))
						m.AddToSetProperty(secName, "SetKey", fmt.Sprintf("Item-%d", (workerID*100)+j))
						m.AppendSummaryMarkdown(secName, "worker progress")
						_ = m.Sections()
					}
				}()
			}

			wg.Wait()

			sections := m.Sections()
			if got, want := len(sections), 4; got != want {
				t.Errorf("got %d sections, want %d", got, want)
			}
		})
	}
}
