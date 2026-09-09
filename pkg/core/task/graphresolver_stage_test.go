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

package coretask

import (
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/google/go-cmp/cmp"
)

func TestResolveGraph_FanInCycle_DownstreamTaskPtPRouting(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name                 string
		initialTasks         []UntypedTask
		availableTasks       []UntypedTask
		wantTaskIDs          []string
		wantDownstreamSource string
	}{
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   prod_high["prod-high (priority=10)"] -->|"FanIn (tag-a)"| consumer_s1["consumer#default-stage-1"]
		//   consumer_s1 -->|"PointToPoint"| prod_low["prod-low (priority=100)"]
		//   prod_low -->|"FanIn (tag-a)"| consumer_s2["consumer#default-stage-2"]
		//   prod_high -->|"FanIn (tag-a)"| consumer_s2
		//   consumer_s1 -.->|"OrderOnly"| consumer_s2
		//   consumer_s2 -->|"PointToPoint"| downstream["downstream"]
		// ```
		{
			name: "downstream task depending on multi-stage consumer is routed from stage-2 when not leading to feedback",
			initialTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("downstream", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("downstream", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}),
			},
			wantTaskIDs: []string{
				"prod-high#default",
				"consumer#default-stage-1",
				"prod-low#default",
				"consumer#default-stage-2",
				"downstream#default",
			},
			wantDownstreamSource: "consumer#default-stage-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err != nil {
				t.Fatalf("ResolveGraph() error = %v, want nil", err)
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}
			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("task order mismatch (-want +got):\n%s", diff)
			}

			incoming := taskSet.IncomingEdges("downstream#default")
			if len(incoming) != 1 {
				t.Fatalf("downstream incoming edges count = %d, want 1", len(incoming))
			}
			if gotSource := incoming[0].SourceImplID; gotSource != tc.wantDownstreamSource {
				t.Errorf("downstream incoming edge source = %q, want %q", gotSource, tc.wantDownstreamSource)
			}
			if incoming[0].Kind != taskid.EdgeKindData {
				t.Errorf("downstream incoming edge kind = %v, want %v", incoming[0].Kind, taskid.EdgeKindData)
			}

			incomingS2 := taskSet.IncomingEdges("consumer#default-stage-2")
			hasInterStageOrderEdge := false
			for _, e := range incomingS2 {
				if e.SourceImplID == "consumer#default-stage-1" && e.Kind == taskid.EdgeKindOrderOnly {
					hasInterStageOrderEdge = true
					break
				}
			}
			if !hasInterStageOrderEdge {
				t.Errorf("expected inter-stage OrderOnly edge from stage-1 to stage-2, but not found")
			}
		})
	}
}

func TestResolveGraph_FanInCycle_UpstreamTaskPtPRouting(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name               string
		initialTasks       []UntypedTask
		availableTasks     []UntypedTask
		wantTaskIDs        []string
		wantUpstreamSource string
	}{
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   upstream["upstream"] -->|"PointToPoint"| consumer_s1["consumer#default-stage-1"]
		//   upstream -->|"PointToPoint"| consumer_s2["consumer#default-stage-2"]
		//   prod_high["prod-high (priority=10)"] -->|"FanIn (tag-a)"| consumer_s1
		//   consumer_s1 -->|"PointToPoint"| prod_low["prod-low (priority=100)"]
		//   prod_low -->|"FanIn (tag-a)"| consumer_s2
		//   prod_high -->|"FanIn (tag-a)"| consumer_s2
		//   consumer_s1 -.->|"OrderOnly"| consumer_s2
		// ```
		{
			name: "upstream task required by multi-stage consumer is routed to both stage-1 and stage-2",
			initialTasks: []UntypedTask{
				createMockTask("upstream", "default", nil),
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("upstream"),
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("upstream", "default", nil),
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("upstream"),
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantTaskIDs: []string{
				"prod-high#default",
				"upstream#default",
				"consumer#default-stage-1",
				"prod-low#default",
				"consumer#default-stage-2",
			},
			wantUpstreamSource: "upstream#default",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err != nil {
				t.Fatalf("ResolveGraph() error = %v, want nil", err)
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}
			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("task order mismatch (-want +got):\n%s", diff)
			}

			incomingS1 := taskSet.IncomingEdges("consumer#default-stage-1")
			hasUpstreamS1 := false
			for _, e := range incomingS1 {
				if e.SourceImplID == tc.wantUpstreamSource {
					hasUpstreamS1 = true
					if e.Kind != taskid.EdgeKindData {
						t.Errorf("consumer stage-1 incoming edge kind = %v, want %v", e.Kind, taskid.EdgeKindData)
					}
					break
				}
			}
			if !hasUpstreamS1 {
				t.Errorf("expected incoming edge from %q to stage-1, but not found", tc.wantUpstreamSource)
			}

			incomingS2 := taskSet.IncomingEdges("consumer#default-stage-2")
			hasUpstreamS2 := false
			for _, e := range incomingS2 {
				if e.SourceImplID == tc.wantUpstreamSource {
					hasUpstreamS2 = true
					if e.Kind != taskid.EdgeKindData {
						t.Errorf("consumer stage-2 incoming edge kind = %v, want %v", e.Kind, taskid.EdgeKindData)
					}
					break
				}
			}
			if !hasUpstreamS2 {
				t.Errorf("expected incoming edge from %q to stage-2, but not found", tc.wantUpstreamSource)
			}
		})
	}
}

func TestResolveGraph_FanInCycle_SplitProducerFanInRouting(t *testing.T) {
	tagA := NewTag[any]("tag-a")
	tagB := NewTag[any]("tag-b")

	testCases := []struct {
		name                 string
		initialTasks         []UntypedTask
		availableTasks       []UntypedTask
		wantTaskIDs          []string
		wantDownstreamSource string
	}{
		{
			name: "downstream task consumes tag provided by split multi-stage consumer",
			initialTasks: []UntypedTask{
				createMockTask("prod-bootstrap", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, ProvidesTag(tagB), AllowMultiStageExecution()),
				createMockTask("prod-feedback", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("downstream", "default", []Dependency{
					tagB.Ref(),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-bootstrap", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, ProvidesTag(tagB), AllowMultiStageExecution()),
				createMockTask("prod-feedback", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("downstream", "default", []Dependency{
					tagB.Ref(),
				}),
			},
			wantTaskIDs: []string{
				"prod-bootstrap#default",
				"consumer#default-stage-1",
				"prod-feedback#default",
				"consumer#default-stage-2",
				"downstream#default",
			},
			wantDownstreamSource: "consumer#default-stage-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}

			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("task order mismatch (-want +got):\n%s", diff)
			}

			incoming := taskSet.IncomingEdges("downstream#default")
			if len(incoming) != 1 {
				t.Fatalf("downstream incoming edges count = %d, want 1", len(incoming))
			}
			if gotSource := incoming[0].SourceImplID; gotSource != tc.wantDownstreamSource {
				t.Errorf("downstream incoming edge source = %q, want %q", gotSource, tc.wantDownstreamSource)
			}
			if incoming[0].Kind != taskid.EdgeKindData {
				t.Errorf("downstream incoming edge kind = %v, want %v", incoming[0].Kind, taskid.EdgeKindData)
			}

			gotBoundRefs := taskSet.BoundReferenceIDsForTaskImplWithTag("downstream#default", tagB.ID())
			if diff := cmp.Diff([]string{"consumer"}, gotBoundRefs); diff != "" {
				t.Errorf("BoundReferenceIDsForTaskImplWithTag mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveGraph_FanInCycle_ChainedSplitConsumersPtPRouting(t *testing.T) {
	tagA := NewTag[any]("tag-a")
	tagB := NewTag[any]("tag-b")

	testCases := []struct {
		name                      string
		initialTasks              []UntypedTask
		availableTasks            []UntypedTask
		wantTaskIDs               []string
		wantConsumerBStage1Source string
		wantConsumerBStage2Source string
	}{
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   prod_high_a["prod-high-a (priority=10)"] -->|"FanIn (tag-a)"| consumer_a_s1["consumer-a#default-stage-1"]
		//   consumer_a_s1 -->|"PointToPoint"| prod_low_a["prod-low-a (priority=100)"]
		//   prod_low_a -->|"FanIn (tag-a)"| consumer_a_s2["consumer-a#default-stage-2"]
		//   prod_high_a -->|"FanIn (tag-a)"| consumer_a_s2
		//   consumer_a_s1 -.->|"OrderOnly"| consumer_a_s2
		//   consumer_a_s2 -->|"PointToPoint"| consumer_b_s1["consumer-b#default-stage-1"]
		//   consumer_a_s2 -->|"PointToPoint"| consumer_b_s2["consumer-b#default-stage-2"]
		//   prod_high_b["prod-high-b (priority=10)"] -->|"FanIn (tag-b)"| consumer_b_s1
		//   consumer_b_s1 -->|"PointToPoint"| prod_low_b["prod-low-b (priority=100)"]
		//   prod_low_b -->|"FanIn (tag-b)"| consumer_b_s2
		//   prod_high_b -->|"FanIn (tag-b)"| consumer_b_s2
		//   consumer_b_s1 -.->|"OrderOnly"| consumer_b_s2
		// ```
		{
			name: "downstream split consumer depending on upstream split consumer receives PtP edges from stage-2",
			initialTasks: []UntypedTask{
				createMockTask("prod-high-a", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer-a", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low-a", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer-a"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("prod-high-b", "default", nil, ProvidesTag(tagB, WithTagPriority(10))),
				createMockTask("consumer-b", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer-a"),
					tagB.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low-b", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer-b"),
				}, ProvidesTag(tagB, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-high-a", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer-a", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low-a", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer-a"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("prod-high-b", "default", nil, ProvidesTag(tagB, WithTagPriority(10))),
				createMockTask("consumer-b", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer-a"),
					tagB.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low-b", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer-b"),
				}, ProvidesTag(tagB, WithTagPriority(100))),
			},
			wantTaskIDs: []string{
				"prod-high-a#default",
				"consumer-a#default-stage-1",
				"prod-high-b#default",
				"prod-low-a#default",
				"consumer-a#default-stage-2",
				"consumer-b#default-stage-1",
				"prod-low-b#default",
				"consumer-b#default-stage-2",
			},
			wantConsumerBStage1Source: "consumer-a#default-stage-2",
			wantConsumerBStage2Source: "consumer-a#default-stage-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}

			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("task order mismatch (-want +got):\n%s", diff)
			}

			incomingS1 := taskSet.IncomingEdges("consumer-b#default-stage-1")
			hasS1Edge := false
			for _, e := range incomingS1 {
				if e.SourceImplID == tc.wantConsumerBStage1Source {
					hasS1Edge = true
					if e.Kind != taskid.EdgeKindData {
						t.Errorf("consumer-b stage-1 incoming edge kind = %v, want %v", e.Kind, taskid.EdgeKindData)
					}
					break
				}
			}
			if !hasS1Edge {
				t.Errorf("expected incoming edge from %q to consumer-b stage-1, but not found", tc.wantConsumerBStage1Source)
			}

			incomingS2 := taskSet.IncomingEdges("consumer-b#default-stage-2")
			hasS2Edge := false
			for _, e := range incomingS2 {
				if e.SourceImplID == tc.wantConsumerBStage2Source {
					hasS2Edge = true
					if e.Kind != taskid.EdgeKindData {
						t.Errorf("consumer-b stage-2 incoming edge kind = %v, want %v", e.Kind, taskid.EdgeKindData)
					}
					break
				}
			}
			if !hasS2Edge {
				t.Errorf("expected incoming edge from %q to consumer-b stage-2, but not found", tc.wantConsumerBStage2Source)
			}
		})
	}
}
