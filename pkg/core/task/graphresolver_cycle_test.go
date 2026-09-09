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
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/google/go-cmp/cmp"
)

func TestResolveGraph_FanInCycle_PriorityDifference(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name                  string
		initialTasks          []UntypedTask
		availableTasks        []UntypedTask
		wantTaskIDs           []string
		wantBoundRefIDsByTask map[string][]string
	}{
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   prod_high["prod-high (priority=10)"] -->|"FanIn (tag-a)"| consumer_s1["consumer#stage-1"]
		//   consumer_s1 -->|"PointToPoint"| prod_low["prod-low (priority=100)"]
		//   prod_low -->|"FanIn (tag-a)"| consumer_s2["consumer#stage-2"]
		//   prod_high -->|"FanIn (tag-a)"| consumer_s2
		//   consumer_s1 -.->|"OrderOnly"| consumer_s2
		// ```
		{
			name: "higher priority fan-in edge boots stage-1 and cyclic lower priority edge is resolved in stage-2",
			initialTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantTaskIDs: []string{"prod-high#default", "consumer#default-stage-1", "prod-low#default", "consumer#default-stage-2"},
			wantBoundRefIDsByTask: map[string][]string{
				"consumer#default-stage-1": {"prod-high"},
				"consumer#default-stage-2": {"prod-high", "prod-low"},
			},
		},
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   prod_high["prod-high (priority=10)"] -->|"FanIn (tag-a)"| consumer_s1["consumer#stage-1"]
		//   prod_high -->|"FanIn (tag-a)"| consumer_s2["consumer#stage-2"]
		//   consumer_s1 -->|"FanIn (tag-a, self-loop)"| consumer_s2
		//   consumer_s1 -.->|"OrderOnly"| consumer_s2
		// ```
		{
			name: "higher priority external producer boots stage-1 and lower priority self-loop edge is resolved in stage-2",
			initialTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution(), ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution(), ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantTaskIDs: []string{"prod-high#default", "consumer#default-stage-1", "consumer#default-stage-2"},
			wantBoundRefIDsByTask: map[string][]string{
				"consumer#default-stage-1": {"prod-high"},
				"consumer#default-stage-2": {"consumer", "prod-high"},
			},
		},
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   prod_high["prod-high (priority=10)"] -->|"FanIn (tag-a)"| consumer_s1["consumer#stage-1"]
		//   consumer_s1 -->|"PointToPoint"| task_mid["task-mid"]
		//   task_mid -->|"PointToPoint"| prod_low["prod-low (priority=100)"]
		//   prod_low -->|"FanIn (tag-a)"| consumer_s2["consumer#stage-2"]
		//   prod_high -->|"FanIn (tag-a)"| consumer_s2
		//   consumer_s1 -.->|"OrderOnly"| consumer_s2
		// ```
		{
			name: "multi-hop transitive cycle with lower priority producer is resolved across stage-1 and stage-2",
			initialTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("task-mid", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("task-mid"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("task-mid", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("task-mid"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantTaskIDs: []string{"prod-high#default", "consumer#default-stage-1", "task-mid#default", "prod-low#default", "consumer#default-stage-2"},
			wantBoundRefIDsByTask: map[string][]string{
				"consumer#default-stage-1": {"prod-high"},
				"consumer#default-stage-2": {"prod-high", "prod-low"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err != nil {
				t.Fatalf("ResolveGraph() unexpected error = %v", err)
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}

			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("ResolveGraph() task IDs mismatch (-want +got):\n%s", diff)
			}

			for taskImplID, wantRefs := range tc.wantBoundRefIDsByTask {
				gotRefs := taskSet.BoundReferenceIDsForTaskWithTag(taskImplID, tagA.ID())
				if diff := cmp.Diff(wantRefs, gotRefs); diff != "" {
					t.Errorf("BoundReferenceIDsForTaskWithTag(%q) mismatch (-want +got):\n%s", taskImplID, diff)
				}
			}
		})
	}
}

func TestResolveGraph_FanInCycle_AmbiguousPriorityTie(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name           string
		initialTasks   []UntypedTask
		availableTasks []UntypedTask
		wantErrMsg     string
	}{
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   prod_a["prod-a (priority=100)"] -->|"FanIn (tag-a)"| consumer["consumer"]
		//   consumer -->|"PointToPoint"| prod_b["prod-b (priority=100)"]
		//   prod_b -.->|"FanIn (tag-a, ambiguous tie)"| consumer
		// ```
		{
			name: "equal priorities in cycle reject resolution with ambiguous priority tie error",
			initialTasks: []UntypedTask{
				createMockTask("prod-a", "default", nil, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-b", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-a", "default", nil, ProvidesTag(tagA, WithTagPriority(100))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-b", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantErrMsg: "ambiguous FanIn priority between prod-a#default and prod-b#default for tag tag-a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err == nil {
				t.Fatal("expected error for ambiguous priority tie, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantErrMsg)
			}
		})
	}
}

func TestResolveGraph_FanInCycle_MissingAllowMultiStageExecution(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name           string
		initialTasks   []UntypedTask
		availableTasks []UntypedTask
		wantErrMsg     string
	}{
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   prod_high["prod-high (priority=10)"] -->|"FanIn (tag-a)"| consumer["consumer (no AllowMultiStageExecution)"]
		//   consumer -->|"PointToPoint"| prod_low["prod-low (priority=100)"]
		//   prod_low -.->|"FanIn (tag-a, cycle rejected)"| consumer
		// ```
		{
			name: "cycle involving consumer without AllowMultiStageExecution rejects resolution",
			initialTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantErrMsg: "task consumer#default requires multi-stage execution but lacks AllowMultiStageExecution label",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err == nil {
				t.Fatal("expected error for missing AllowMultiStageExecution, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantErrMsg)
			}
		})
	}
}

func TestResolveGraph_FanInCycle_SoleProducerNoBootstrap(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name           string
		initialTasks   []UntypedTask
		availableTasks []UntypedTask
		wantErrMsg     string
	}{
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   consumer["consumer"] -->|"PointToPoint"| prod_low["prod-low (priority=100)"]
		//   prod_low -.->|"FanIn (tag-a, cycle with no bootstrap)"| consumer
		// ```
		{
			name: "sole cyclic producer with AllowMultiStageExecution rejects resolution because no higher priority producer exists",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-low", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantErrMsg: "task consumer#default has circular fan-in dependency on tag tag-a with no higher-priority producer to bootstrap execution",
		},
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   consumer["consumer"] -.->|"FanIn (tag-a, self-loop with no bootstrap)"| consumer
		// ```
		{
			name: "self-loop sole producer with AllowMultiStageExecution rejects resolution because no higher priority producer exists",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution(), ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution(), ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantErrMsg: "task consumer#default has circular fan-in dependency on tag tag-a with no higher-priority producer to bootstrap execution",
		},
		// Mermaid task graph:
		// ```mermaid
		// graph TD
		//   consumer["consumer"] -->|"PointToPoint"| prod_cyclic_high["prod-cyclic-high (priority=10)"]
		//   prod_cyclic_high -.->|"FanIn (tag-a, higher priority cyclic)"| consumer
		//   prod_acyclic_low["prod-acyclic-low (priority=100)"] -->|"FanIn (tag-a, lower priority acyclic)"| consumer
		// ```
		{
			name: "cyclic producer having higher priority than acyclic producer rejects resolution because bootstrap must be higher priority",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-cyclic-high", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("prod-acyclic-low", "default", nil, ProvidesTag(tagA, WithTagPriority(100))),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					tagA.Ref(),
				}, AllowMultiStageExecution()),
				createMockTask("prod-cyclic-high", "default", []Dependency{
					taskid.NewTaskReference[any]("consumer"),
				}, ProvidesTag(tagA, WithTagPriority(10))),
				createMockTask("prod-acyclic-low", "default", nil, ProvidesTag(tagA, WithTagPriority(100))),
			},
			wantErrMsg: "task consumer#default has circular fan-in dependency on tag tag-a with no higher-priority producer to bootstrap execution",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if err == nil {
				t.Fatal("expected error for circular fan-in without bootstrap, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantErrMsg)
			}
		})
	}
}

func TestResolveGraph_ShuffleInvariance(t *testing.T) {
	tagA := NewTag[any]("tag-a")
	tagB := NewTag[any]("tag-b")

	baseInitialTasks := []UntypedTask{
		createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
		createMockTask("consumer", "default", []Dependency{
			tagA.Ref(),
			tagB.Ref(),
		}, AllowMultiStageExecution()),
		createMockTask("prod-low", "default", []Dependency{
			taskid.NewTaskReference[any]("consumer"),
		}, ProvidesTag(tagA, WithTagPriority(100))),
		createMockTask("tag-b-prod", "default", nil, ProvidesTag(tagB, WithTagPriority(50))),
	}
	baseAvailableTasks := []UntypedTask{
		createMockTask("prod-high", "default", nil, ProvidesTag(tagA, WithTagPriority(10))),
		createMockTask("consumer", "default", []Dependency{
			tagA.Ref(),
			tagB.Ref(),
		}, AllowMultiStageExecution()),
		createMockTask("prod-low", "default", []Dependency{
			taskid.NewTaskReference[any]("consumer"),
		}, ProvidesTag(tagA, WithTagPriority(100))),
		createMockTask("tag-b-prod", "default", nil, ProvidesTag(tagB, WithTagPriority(50))),
	}

	baselineSet, err := ResolveGraph(baseInitialTasks, baseAvailableTasks, nil)
	if err != nil {
		t.Fatalf("baseline ResolveGraph() failed: %v", err)
	}

	var wantTaskIDs []string
	for _, task := range baselineSet.GetAll() {
		wantTaskIDs = append(wantTaskIDs, task.UntypedID().String())
	}

	r := rand.New(rand.NewPCG(42, 1024))
	for i := range 100 {
		shuffledInitialTasks := slices.Clone(baseInitialTasks)
		r.Shuffle(len(shuffledInitialTasks), func(i, j int) {
			shuffledInitialTasks[i], shuffledInitialTasks[j] = shuffledInitialTasks[j], shuffledInitialTasks[i]
		})

		shuffledAvailableTasks := slices.Clone(baseAvailableTasks)
		r.Shuffle(len(shuffledAvailableTasks), func(i, j int) {
			shuffledAvailableTasks[i], shuffledAvailableTasks[j] = shuffledAvailableTasks[j], shuffledAvailableTasks[i]
		})

		shuffledSet, err := ResolveGraph(shuffledInitialTasks, shuffledAvailableTasks, nil)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}

		var gotTaskIDs []string
		for _, task := range shuffledSet.GetAll() {
			gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
		}

		if diff := cmp.Diff(wantTaskIDs, gotTaskIDs); diff != "" {
			t.Fatalf("iteration %d: task order mismatch (-want +got):\n%s", i, diff)
		}
	}
}

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
		})
	}
}
