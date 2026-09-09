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
		name            string
		initialTasks    []UntypedTask
		availableTasks  []UntypedTask
		wantTaskIDs     []string
		wantBoundRefIDs []string
	}{
		{
			name: "higher priority fan-in edge is kept and lower priority edge creating cycle is pruned",
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
			wantTaskIDs:     []string{"prod-high#default", "consumer#default", "prod-low#default"},
			wantBoundRefIDs: []string{"prod-high"},
		},
		{
			name: "higher priority external producer is kept and lower priority self-loop edge is pruned",
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
			wantTaskIDs:     []string{"prod-high#default", "consumer#default"},
			wantBoundRefIDs: []string{"prod-high"},
		},
		{
			name: "multi-hop transitive cycle with lower priority producer is pruned",
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
			wantTaskIDs:     []string{"prod-high#default", "consumer#default", "task-mid#default", "prod-low#default"},
			wantBoundRefIDs: []string{"prod-high"},
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

			gotBoundRefIDs := taskSet.BoundReferenceIDsWithTag(tagA.ID())
			if diff := cmp.Diff(tc.wantBoundRefIDs, gotBoundRefIDs); diff != "" {
				t.Errorf("BoundReferenceIDsWithTag() mismatch (-want +got):\n%s", diff)
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
