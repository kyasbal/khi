// Copyright 2024 Google LLC
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
	"context"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/google/go-cmp/cmp"
)

type mockUntypedTask struct {
	id           taskid.UntypedTaskImplementationID
	labels       *typedmap.ReadonlyTypedMap
	dependencies []Dependency
}

func (m *mockUntypedTask) UntypedID() taskid.UntypedTaskImplementationID { return m.id }
func (m *mockUntypedTask) Labels() *typedmap.ReadonlyTypedMap            { return m.labels }
func (m *mockUntypedTask) Dependencies() []Dependency                    { return m.dependencies }
func (m *mockUntypedTask) UntypedRun(ctx context.Context) (any, error)   { return nil, nil }

var _ UntypedTask = (*mockUntypedTask)(nil)

func createMockTask(refID, implHash string, deps []Dependency, labelOpts ...LabelOpt) UntypedTask {
	id := taskid.NewImplementationID(taskid.NewTaskReference[any](refID), implHash)
	return &mockUntypedTask{
		id:           id,
		labels:       NewLabelSet(labelOpts...),
		dependencies: deps,
	}
}

func TestResolveGraph_MandatoryClosure(t *testing.T) {
	testCases := []struct {
		name           string
		initialTasks   []UntypedTask
		availableTasks []UntypedTask
		wantTaskIDs    []string
		wantErr        bool
	}{
		{
			name: "single task without dependencies",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", nil),
			},
			availableTasks: []UntypedTask{
				createMockTask("task-a", "default", nil),
			},
			wantTaskIDs: []string{"task-a#default"},
			wantErr:     false,
		},
		{
			name: "transitive mandatory dependencies are resolved",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", []Dependency{
					taskid.NewTaskReference[any]("task-b"),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("task-a", "default", []Dependency{
					taskid.NewTaskReference[any]("task-b"),
				}),
				createMockTask("task-b", "default", []Dependency{
					taskid.NewTaskReference[any]("task-c"),
				}),
				createMockTask("task-c", "default", nil),
			},
			wantTaskIDs: []string{"task-c#default", "task-b#default", "task-a#default"},
			wantErr:     false,
		},
		{
			name: "required task label automatically pulls task into graph",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", nil),
			},
			availableTasks: []UntypedTask{
				createMockTask("task-a", "default", nil),
				createMockTask("required-task", "default", nil, WithLabelValue(LabelKeyRequiredTask, true)),
			},
			wantTaskIDs: []string{"required-task#default", "task-a#default"},
			wantErr:     false,
		},
		{
			name: "missing mandatory dependency returns error",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", []Dependency{
					taskid.NewTaskReference[any]("non-existent"),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("task-a", "default", []Dependency{
					taskid.NewTaskReference[any]("non-existent"),
				}),
			},
			wantTaskIDs: nil,
			wantErr:     true,
		},
		{
			name: "highest priority task implementation is selected",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider"),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider"),
				}),
				createMockTask("provider", "low-priority", nil, WithLabelValue(LabelKeyTaskSelectionPriority, 10)),
				createMockTask("provider", "high-priority", nil, WithLabelValue(LabelKeyTaskSelectionPriority, 100)),
			},
			wantTaskIDs: []string{"provider#high-priority", "consumer#default"},
			wantErr:     false,
		},
		{
			name: "same priority task implementation tie-breaks deterministically by implementation ID",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider"),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider"),
				}),
				createMockTask("provider", "impl-z", nil, WithLabelValue(LabelKeyTaskSelectionPriority, 50)),
				createMockTask("provider", "impl-a", nil, WithLabelValue(LabelKeyTaskSelectionPriority, 50)),
			},
			wantTaskIDs: []string{"provider#impl-a", "consumer#default"},
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, nil)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ResolveGraph() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}

			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("ResolveGraph() task IDs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveGraph_OptionalDependencies(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name           string
		initialTasks   []UntypedTask
		availableTasks []UntypedTask
		wantTaskIDs    []string
		wantEdgeCount  int
	}{
		{
			name: "optional dependency is activated when target exists in graph",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider", taskid.Optional),
				}),
				createMockTask("provider", "default", nil),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider", taskid.Optional),
				}),
				createMockTask("provider", "default", nil),
			},
			wantTaskIDs:   []string{"provider#default", "consumer#default"},
			wantEdgeCount: 1,
		},
		{
			name: "optional dependency is NOT activated when target is not in graph",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider", taskid.Optional),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider", taskid.Optional),
				}),
				createMockTask("provider", "default", nil),
			},
			wantTaskIDs:   []string{"consumer#default"},
			wantEdgeCount: 0,
		},
		{
			name: "fan-in with ScopeActiveGraph only includes existing tasks",
			initialTasks: []UntypedTask{
				createMockTask("collector", "default", []Dependency{
					tagA.Ref(),
				}),
				createMockTask("prod1", "default", nil, ProvidesTag(tagA)),
			},
			availableTasks: []UntypedTask{
				createMockTask("collector", "default", []Dependency{
					tagA.Ref(),
				}),
				createMockTask("prod1", "default", nil, ProvidesTag(tagA)),
				createMockTask("prod2", "default", nil, ProvidesTag(tagA)), // not in graph
			},
			wantTaskIDs:   []string{"prod1#default", "collector#default"},
			wantEdgeCount: 1,
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
				t.Errorf("ResolveGraph() task IDs mismatch (-want +got):\n%s", diff)
			}
			if len(taskSet.Edges()) != tc.wantEdgeCount {
				t.Errorf("got %d edges, want %d", len(taskSet.Edges()), tc.wantEdgeCount)
			}
		})
	}
}

func TestResolveGraph_FanInScopeAll(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name           string
		initialTasks   []UntypedTask
		availableTasks []UntypedTask
		disabledTasks  []UntypedTask
		wantTaskIDs    []string
		wantErr        bool
	}{
		{
			name: "ScopeAll pulls in all available providers and their transitive dependencies",
			initialTasks: []UntypedTask{
				createMockTask("collector", "default", []Dependency{
					tagA.Ref(FromAll),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("collector", "default", []Dependency{
					tagA.Ref(FromAll),
				}),
				createMockTask("dep1", "default", nil),
				createMockTask("prod1", "default", []Dependency{
					taskid.NewTaskReference[any]("dep1"),
				}, ProvidesTag(tagA)),
				createMockTask("prod2", "default", nil, ProvidesTag(tagA)),
			},
			wantTaskIDs: []string{"dep1#default", "prod1#default", "prod2#default", "collector#default"},
			wantErr:     false,
		},
		{
			name: "ScopeAll fails when a provider has a disabled dependency",
			initialTasks: []UntypedTask{
				createMockTask("collector", "default", []Dependency{
					tagA.Ref(FromAll),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("collector", "default", []Dependency{
					tagA.Ref(FromAll),
				}),
				createMockTask("dep1", "default", nil),
				createMockTask("prod1", "default", []Dependency{
					taskid.NewTaskReference[any]("dep1"),
				}, ProvidesTag(tagA)),
			},
			disabledTasks: []UntypedTask{
				createMockTask("dep1", "default", nil),
			},
			wantTaskIDs: nil,
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, tc.disabledTasks)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ResolveGraph() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}

			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("ResolveGraph() task IDs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveGraph_FanInScopeActiveFeatures(t *testing.T) {
	tagDiscovery := NewTag[any]("discovery")

	testCases := []struct {
		name        string
		setup       func() (initialTasks, availableTasks, disabledTasks []UntypedTask)
		wantTaskIDs []string
		wantErr     bool
	}{
		{
			name: "anchor check includes anchored provider but excludes unselected feature provider",
			setup: func() ([]UntypedTask, []UntypedTask, []UntypedTask) {
				commonAncestor := createMockTask("common-ancestor", "default", nil)
				prod1 := createMockTask("prod1", "default", []Dependency{
					taskid.NewTaskReference[any]("common-ancestor"),
				}, ProvidesTag(tagDiscovery))
				unselectedFeature := createMockTask("unselected-feature", "default", nil)
				prod2 := createMockTask("prod2", "default", []Dependency{
					taskid.NewTaskReference[any]("unselected-feature"),
				}, ProvidesTag(tagDiscovery))
				collector := createMockTask("collector", "default", []Dependency{
					tagDiscovery.Ref(FromActiveFeatures),
				})
				return []UntypedTask{collector, commonAncestor},
					[]UntypedTask{collector, commonAncestor, prod1, unselectedFeature, prod2},
					[]UntypedTask{unselectedFeature}
			},
			wantTaskIDs: []string{"common-ancestor#default", "prod1#default", "collector#default"},
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initial, available, disabled := tc.setup()
			taskSet, err := ResolveGraph(initial, available, disabled)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ResolveGraph() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}

			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("ResolveGraph() task IDs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveGraph_EdgeDeduplication(t *testing.T) {
	tagA := NewTag[any]("tag-a")

	testCases := []struct {
		name         string
		initialTasks []UntypedTask
		wantEdges    []taskid.TaskEdge
	}{
		{
			name: "point-to-point data dependency and fan-in order-only deduplicate to data edge with tag preserved",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", nil, ProvidesTag(tagA)),
				createMockTask("task-b", "default", []Dependency{
					taskid.NewTaskReference[any]("task-a"),
					tagA.Ref(taskid.OrderOnly),
				}),
			},
			wantEdges: []taskid.TaskEdge{
				{
					SourceRefID: "task-a",
					SourceID:    "task-a#default",
					TargetID:    "task-b#default",
					Kind:        taskid.EdgeKindData,
					Condition:   taskid.ConditionRequired,
					Cardinality: taskid.CardinalityPointToPoint,
					Tag:         "tag-a",
					Priority:    100,
				},
			},
		},
		{
			name: "order-only upgraded to data edge with priority and condition merging",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", nil, ProvidesTag(tagA, WithTagPriority(20))),
				createMockTask("task-b", "default", []Dependency{
					taskid.NewTaskReference[any]("task-a", taskid.OrderOnly, taskid.Optional),
					tagA.Ref(),
				}),
			},
			wantEdges: []taskid.TaskEdge{
				{
					SourceRefID: "task-a",
					SourceID:    "task-a#default",
					TargetID:    "task-b#default",
					Kind:        taskid.EdgeKindData,
					Condition:   taskid.ConditionRequired,
					Cardinality: taskid.CardinalityPointToPoint,
					Tag:         "tag-a",
					Priority:    20,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.initialTasks, nil)
			if err != nil {
				t.Fatalf("ResolveGraph() unexpected error = %v", err)
			}

			edges := taskSet.Edges()
			if diff := cmp.Diff(tc.wantEdges, edges); diff != "" {
				t.Errorf("Edges() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveGraph_CyclicDependency(t *testing.T) {
	testCases := []struct {
		name         string
		initialTasks []UntypedTask
		wantErrMsg   string
	}{
		{
			name: "three-task cycle returns cyclic dependency error",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", []Dependency{
					taskid.NewTaskReference[any]("task-b"),
				}),
				createMockTask("task-b", "default", []Dependency{
					taskid.NewTaskReference[any]("task-c"),
				}),
				createMockTask("task-c", "default", []Dependency{
					taskid.NewTaskReference[any]("task-a"),
				}),
			},
			wantErrMsg: "cyclic dependency",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ResolveGraph(tc.initialTasks, tc.initialTasks, nil)
			if err == nil {
				t.Fatal("expected error for cyclic dependency, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantErrMsg)
			}
		})
	}
}

func TestResolveGraph_DisabledTasks(t *testing.T) {
	testCases := []struct {
		name           string
		initialTasks   []UntypedTask
		availableTasks []UntypedTask
		disabledTasks  []UntypedTask
		wantTaskIDs    []string
		wantErr        bool
	}{
		{
			name: "initial task is disabled returns error",
			initialTasks: []UntypedTask{
				createMockTask("task-a", "default", nil),
			},
			availableTasks: []UntypedTask{
				createMockTask("task-a", "default", nil),
			},
			disabledTasks: []UntypedTask{
				createMockTask("task-a", "default", nil),
			},
			wantErr: true,
		},
		{
			name: "mandatory dependency of initial task is disabled returns error",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider"),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider"),
				}),
				createMockTask("provider", "default", nil),
			},
			disabledTasks: []UntypedTask{
				createMockTask("provider", "default", nil),
			},
			wantErr: true,
		},
		{
			name: "optional dependency is disabled is safely ignored",
			initialTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider", taskid.Optional),
				}),
			},
			availableTasks: []UntypedTask{
				createMockTask("consumer", "default", []Dependency{
					taskid.NewTaskReference[any]("provider", taskid.Optional),
				}),
				createMockTask("provider", "default", nil),
			},
			disabledTasks: []UntypedTask{
				createMockTask("provider", "default", nil),
			},
			wantTaskIDs: []string{"consumer#default"},
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskSet, err := ResolveGraph(tc.initialTasks, tc.availableTasks, tc.disabledTasks)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ResolveGraph() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			var gotTaskIDs []string
			for _, task := range taskSet.GetAll() {
				gotTaskIDs = append(gotTaskIDs, task.UntypedID().String())
			}

			if diff := cmp.Diff(tc.wantTaskIDs, gotTaskIDs); diff != "" {
				t.Errorf("ResolveGraph() task IDs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type mockUnknownScopeFanIn struct {
	tag string
}

var _ taskid.FanInDescriptor = (*mockUnknownScopeFanIn)(nil)

func (m mockUnknownScopeFanIn) DescriptorKind() taskid.EdgeKind {
	return taskid.EdgeKindData
}

func (m mockUnknownScopeFanIn) DescriptorCondition() taskid.EdgeCondition {
	return taskid.ConditionRequired
}

func (m mockUnknownScopeFanIn) DescriptorCardinality() taskid.EdgeCardinality {
	return taskid.CardinalityFanIn
}

func (m mockUnknownScopeFanIn) DescriptorScope() taskid.DependencyScope {
	return taskid.DependencyScope(999)
}

func (m mockUnknownScopeFanIn) Tag() string {
	return m.tag
}

func TestResolveGraph_UnknownScope(t *testing.T) {
	testCases := []struct {
		name       string
		task       UntypedTask
		wantErrMsg string
	}{
		{
			name: "unknown DependencyScope returns error",
			task: createMockTask("collector", "default", []Dependency{
				mockUnknownScopeFanIn{tag: "some-tag"},
			}),
			wantErrMsg: "unsupported dependency scope",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ResolveGraph([]UntypedTask{tc.task}, []UntypedTask{tc.task}, nil)
			if err == nil {
				t.Fatal("expected error for unknown DependencyScope, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantErrMsg)
			}
		})
	}
}
