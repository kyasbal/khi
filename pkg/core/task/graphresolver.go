// Copyright 2025 Google LLC
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
	"container/heap"
	"fmt"
	"slices"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
)

// ResolveGraph resolves the task graph starting from initialTasks, drawing dependencies from availableTasks,
// and excluding disabledTasks.
// It applies a deterministic 4.5-phase resolution algorithm and returns a runnable TaskSet containing
// topologically sorted tasks and concrete TaskEdges.
func ResolveGraph(
	initialTasks []UntypedTask,
	availableTasks []UntypedTask,
	disabledTasks []UntypedTask,
) (*TaskSet, error) {
	disabledTaskMap := createDisabledTaskMap(disabledTasks)

	// --- Phase 1: Mandatory Closure ---
	graphTaskMap, err := resolveMandatoryClosure(initialTasks, availableTasks, disabledTaskMap)
	if err != nil {
		return nil, err
	}

	// --- Phase 2: Fan-In Candidate Binding ---
	candidateFanInEdges, err := bindFanInDependencies(graphTaskMap, availableTasks, disabledTaskMap)
	if err != nil {
		return nil, err
	}

	// --- Phase 3: Optional Binding & Mandatory 1:1 Edges ---
	pointToPointEdges := bindPointToPointDependencies(graphTaskMap)

	// --- Phase 3.5: Fan-In Cycle Resolution, Pruning & Edge Deduplication ---
	resolvedTasks, resolvedEdges, boundFanInRefIDs, boundFanInRefIDsByTask, err := resolveFanInEdgesAndCycles(graphTaskMap, pointToPointEdges, candidateFanInEdges)
	if err != nil {
		return nil, err
	}

	// --- Phase 4: Kahn's Algorithm on E = E_data U E_order ---
	return buildAndSortTaskSet(resolvedTasks, resolvedEdges, boundFanInRefIDs, boundFanInRefIDsByTask), nil
}

// createDisabledTaskMap creates a lookup set of reference IDs and implementation IDs for disabled tasks.
func createDisabledTaskMap(disabledTasks []UntypedTask) map[string]struct{} {
	disabledMap := make(map[string]struct{})
	for _, t := range disabledTasks {
		disabledMap[t.UntypedID().ReferenceIDString()] = struct{}{}
		disabledMap[t.UntypedID().String()] = struct{}{}
	}
	return disabledMap
}

// resolveMandatoryClosure computes the initial graph tasks by expanding mandatory point-to-point dependencies
// starting from initialTasks and system-required tasks in availableTasks.
func resolveMandatoryClosure(
	initialTasks []UntypedTask,
	availableTasks []UntypedTask,
	disabledTaskMap map[string]struct{},
) (map[string]UntypedTask, error) {
	graphTaskMap := make(map[string]UntypedTask)
	queue := make([]UntypedTask, 0)

	// Collect initial tasks.
	for _, t := range initialTasks {
		refID := t.UntypedID().ReferenceIDString()
		if _, isDisabled := disabledTaskMap[refID]; isDisabled {
			return nil, fmt.Errorf("initial task %q is explicitly disabled", t.UntypedID())
		}
		if existing, exists := graphTaskMap[refID]; exists {
			if existing.UntypedID().String() != t.UntypedID().String() {
				return nil, fmt.Errorf("conflicting task implementations for reference %q: %s and %s", refID, existing.UntypedID(), t.UntypedID())
			}
		} else {
			graphTaskMap[refID] = t
			queue = append(queue, t)
		}
	}

	// Collect system-required tasks (LabelKeyRequiredTask).
	for _, t := range availableTasks {
		if req, found := typedmap.Get(t.Labels(), LabelKeyRequiredTask); found && req {
			refID := t.UntypedID().ReferenceIDString()
			if _, isDisabled := disabledTaskMap[refID]; isDisabled {
				continue
			}
			if existing, exists := graphTaskMap[refID]; exists {
				if existing.UntypedID().String() != t.UntypedID().String() {
					return nil, fmt.Errorf("conflicting task implementations for reference %q: %s and %s", refID, existing.UntypedID(), t.UntypedID())
				}
			} else {
				graphTaskMap[refID] = t
				queue = append(queue, t)
			}
		}
	}

	// Expand mandatory point-to-point dependency closure.
	if err := expandMandatoryDependencies(queue, graphTaskMap, availableTasks, disabledTaskMap); err != nil {
		return nil, err
	}

	return graphTaskMap, nil
}

// expandMandatoryDependencies expands the mandatory point-to-point dependency closure for tasks in queue.
func expandMandatoryDependencies(
	queue []UntypedTask,
	graphTaskMap map[string]UntypedTask,
	availableTasks []UntypedTask,
	disabledTaskMap map[string]struct{},
) error {
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, dep := range curr.Dependencies() {
			if dep.DescriptorCondition() == taskid.ConditionRequired && dep.DescriptorCardinality() == taskid.CardinalityPointToPoint {
				ptp, ok := dep.(taskid.PointToPointDescriptor)
				if !ok {
					continue
				}
				refID := ptp.ReferenceID()
				if _, exists := graphTaskMap[refID]; !exists {
					if _, isDisabled := disabledTaskMap[refID]; isDisabled {
						return fmt.Errorf("required dependency %q required by %q is disabled", refID, curr.UntypedID())
					}
					targetTask, err := findBestTaskImplementation(refID, availableTasks)
					if err != nil {
						return fmt.Errorf("required dependency %q required by %q not found: %w", refID, curr.UntypedID(), err)
					}
					graphTaskMap[refID] = targetTask
					queue = append(queue, targetTask)
				}
			}
		}
	}
	return nil
}

// bindFanInDependencies resolves candidate fan-in edges for all tasks in graphTaskMap,
// expanding candidate producers according to their DependencyScope.
func bindFanInDependencies(
	graphTaskMap map[string]UntypedTask,
	availableTasks []UntypedTask,
	disabledTaskMap map[string]struct{},
) ([]taskid.TaskEdge, error) {
	tagToGraphTasks := make(map[string][]UntypedTask)
	for _, t := range graphTaskMap {
		for _, tag := range getProvidedTags(t) {
			tagToGraphTasks[tag] = append(tagToGraphTasks[tag], t)
		}
	}
	for tag := range tagToGraphTasks {
		slices.SortFunc(tagToGraphTasks[tag], compareTaskByImplementationID)
	}

	var rawEdges []taskid.TaskEdge

	for _, t := range graphTaskMap {
		for _, dep := range t.Dependencies() {
			if dep.DescriptorCardinality() == taskid.CardinalityFanIn {
				fanInDep, ok := dep.(taskid.FanInDescriptor)
				if !ok {
					continue
				}
				tag := fanInDep.Tag()
				var matchingProducers []UntypedTask

				switch fanInDep.DescriptorScope() {
				case taskid.ScopeAll:
					var err error
					matchingProducers, err = findAndExpandAllTasksByTag(tag, availableTasks, graphTaskMap, disabledTaskMap)
					if err != nil {
						return nil, err
					}
				case taskid.ScopeActiveFeatures:
					matchingProducers = findAndConditionallyExpandActiveFeatureTasks(tag, availableTasks, graphTaskMap, disabledTaskMap)
				case taskid.ScopeActiveGraph:
					matchingProducers = tagToGraphTasks[tag]
				default:
					return nil, fmt.Errorf("unknown or unsupported dependency scope: %v", fanInDep.DescriptorScope())
				}

				for _, p := range matchingProducers {
					priority := typedmap.GetOrDefault(p.Labels(), LabelKeyProvidedTagPriority(tag), DefaultTagPriority)
					rawEdges = append(rawEdges, taskid.TaskEdge{
						SourceRefID: p.UntypedID().ReferenceIDString(),
						SourceID:    p.UntypedID().String(),
						TargetID:    t.UntypedID().String(),
						Kind:        dep.DescriptorKind(),
						Condition:   taskid.ConditionRequired,
						Cardinality: taskid.CardinalityFanIn,
						Tag:         tag,
						Priority:    priority,
					})
				}
			}
		}
	}

	return rawEdges, nil
}

// bindPointToPointDependencies creates point-to-point edges for all tasks in graphTaskMap
// whose target exists in the graph (both required and optional).
func bindPointToPointDependencies(graphTaskMap map[string]UntypedTask) []taskid.TaskEdge {
	var rawEdges []taskid.TaskEdge
	for _, t := range graphTaskMap {
		for _, dep := range t.Dependencies() {
			if dep.DescriptorCardinality() == taskid.CardinalityPointToPoint {
				ptp, ok := dep.(taskid.PointToPointDescriptor)
				if !ok {
					continue
				}
				refID := ptp.ReferenceID()
				if sourceTask, exists := graphTaskMap[refID]; exists {
					rawEdges = append(rawEdges, taskid.TaskEdge{
						SourceRefID: refID,
						SourceID:    sourceTask.UntypedID().String(),
						TargetID:    t.UntypedID().String(),
						Kind:        dep.DescriptorKind(),
						Condition:   dep.DescriptorCondition(),
						Cardinality: taskid.CardinalityPointToPoint,
					})
				}
			}
		}
	}
	return rawEdges
}

// deduplicateAndNormalizeEdges merges duplicate edges between the same source and target.
// EdgeKindData takes precedence over EdgeKindOrderOnly.
// ConditionRequired takes precedence over ConditionOptional.
// Minimum Priority takes precedence.
func deduplicateAndNormalizeEdges(rawEdges []taskid.TaskEdge) []taskid.TaskEdge {
	type edgeKey struct {
		sourceID string
		targetID string
	}
	edgeMap := make(map[edgeKey]taskid.TaskEdge)
	order := make([]edgeKey, 0, len(rawEdges))

	for _, e := range rawEdges {
		key := edgeKey{sourceID: e.SourceID, targetID: e.TargetID}
		if existing, exists := edgeMap[key]; exists {
			if e.Kind == taskid.EdgeKindData {
				existing.Kind = taskid.EdgeKindData
			}
			if e.Condition == taskid.ConditionRequired {
				existing.Condition = taskid.ConditionRequired
			}
			if existing.Priority == 0 || (e.Priority > 0 && e.Priority < existing.Priority) {
				existing.Priority = e.Priority
			}
			if existing.Tag == "" && e.Tag != "" {
				existing.Tag = e.Tag
			}
			if existing.Cardinality == taskid.CardinalityPointToPoint || e.Cardinality == taskid.CardinalityPointToPoint {
				existing.Cardinality = taskid.CardinalityPointToPoint
			}
			if existing.SourceRefID == "" && e.SourceRefID != "" {
				existing.SourceRefID = e.SourceRefID
			}
			edgeMap[key] = existing
		} else {
			edgeMap[key] = e
			order = append(order, key)
		}
	}

	deduped := make([]taskid.TaskEdge, 0, len(order))
	for _, k := range order {
		deduped = append(deduped, edgeMap[k])
	}
	return deduped
}

// buildAndSortTaskSet executes Kahn's algorithm with a min-heap for deterministic ordering.
func buildAndSortTaskSet(
	tasks []UntypedTask,
	edges []taskid.TaskEdge,
	boundFanInRefIDs map[string][]string,
	boundFanInRefIDsByTask map[string]map[string][]string,
) *TaskSet {
	inDegree := make(map[string]int, len(tasks))
	outgoing := make(map[string][]string) // key: source task ID string -> []target task ID string

	// Index tasks by ImplementationID for O(1) ready queue push.
	implToTask := make(map[string]UntypedTask, len(tasks))
	for _, task := range tasks {
		implID := task.UntypedID().String()
		inDegree[implID] = 0
		implToTask[implID] = task
	}

	for _, e := range edges {
		outgoing[e.SourceID] = append(outgoing[e.SourceID], e.TargetID)
		inDegree[e.TargetID]++
	}

	h := &taskMinHeap{}
	heap.Init(h)
	for _, task := range tasks {
		if inDegree[task.UntypedID().String()] == 0 {
			heap.Push(h, task)
		}
	}

	sortedTasks := make([]UntypedTask, 0, len(tasks))
	for h.Len() > 0 {
		curr := heap.Pop(h).(UntypedTask)
		sortedTasks = append(sortedTasks, curr)

		for _, targetID := range outgoing[curr.UntypedID().String()] {
			inDegree[targetID]--
			if inDegree[targetID] == 0 {
				if task, ok := implToTask[targetID]; ok {
					heap.Push(h, task)
				}
			}
		}
	}

	return NewResolvedTaskSet(sortedTasks, edges, boundFanInRefIDs, boundFanInRefIDsByTask)
}

// findBestTaskImplementation finds the task implementation for a reference ID with highest priority.
func findBestTaskImplementation(refID string, availableTasks []UntypedTask) (UntypedTask, error) {
	var candidates []UntypedTask
	for _, task := range availableTasks {
		if task.UntypedID().ReferenceIDString() == refID {
			candidates = append(candidates, task)
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no available task can be referenced as %q", refID)
	}

	// Sort candidates: priority descending, then implementation ID ascending.
	slices.SortFunc(candidates, func(a, b UntypedTask) int {
		pA := typedmap.GetOrDefault(a.Labels(), LabelKeyTaskSelectionPriority, 0)
		pB := typedmap.GetOrDefault(b.Labels(), LabelKeyTaskSelectionPriority, 0)
		if pA != pB {
			return pB - pA // higher priority first
		}
		return strings.Compare(a.UntypedID().String(), b.UntypedID().String())
	})

	return candidates[0], nil
}

// getProvidedTags extracts all tag IDs that the task provides.
func getProvidedTags(task UntypedTask) []string {
	var tags []string
	for _, key := range task.Labels().Keys() {
		if strings.HasPrefix(key, LabelKeyProvidedTagPrefix) {
			provided := typedmap.GetOrDefault(task.Labels(), typedmap.NewTypedKey[bool](key), false)
			if provided {
				tags = append(tags, strings.TrimPrefix(key, LabelKeyProvidedTagPrefix))
			}
		}
	}
	slices.Sort(tags)
	return tags
}

// findAndExpandAllTasksByTag finds all non-disabled tasks providing the given tag and expands their mandatory dependencies.
func findAndExpandAllTasksByTag(
	tag string,
	availableTasks []UntypedTask,
	graphTaskMap map[string]UntypedTask,
	disabledTaskMap map[string]struct{},
) ([]UntypedTask, error) {
	var matching []UntypedTask
	for _, t := range availableTasks {
		refID := t.UntypedID().ReferenceIDString()
		if _, isDisabled := disabledTaskMap[refID]; isDisabled {
			continue
		}
		if slices.Contains(getProvidedTags(t), tag) {
			matching = append(matching, t)
			if _, exists := graphTaskMap[refID]; !exists {
				graphTaskMap[refID] = t
				if err := expandMandatoryDependencies([]UntypedTask{t}, graphTaskMap, availableTasks, disabledTaskMap); err != nil {
					return nil, err
				}
			}
		}
	}
	slices.SortFunc(matching, compareTaskByImplementationID)
	return matching, nil
}

// findAndConditionallyExpandActiveFeatureTasks performs the anchor check on candidate producers.
func findAndConditionallyExpandActiveFeatureTasks(
	tag string,
	availableTasks []UntypedTask,
	graphTaskMap map[string]UntypedTask,
	disabledTaskMap map[string]struct{},
) []UntypedTask {
	var matching []UntypedTask
	for _, t := range availableTasks {
		refID := t.UntypedID().ReferenceIDString()
		if _, isDisabled := disabledTaskMap[refID]; isDisabled {
			continue
		}
		if slices.Contains(getProvidedTags(t), tag) {
			if _, inGraph := graphTaskMap[refID]; inGraph {
				matching = append(matching, t)
				continue
			}

			// Perform transaction-like anchor closure check
			subgraph, ok := tryAnchorClosure(t, availableTasks, graphTaskMap, disabledTaskMap)
			if ok {
				for _, subTask := range subgraph {
					graphTaskMap[subTask.UntypedID().ReferenceIDString()] = subTask
				}
				matching = append(matching, t)
			}
		}
	}
	slices.SortFunc(matching, compareTaskByImplementationID)
	return matching
}

// tryAnchorClosure verifies that all mandatory dependency branches of candidate reach existing tasks in graphTaskMap
// without passing through any disabled tasks.
func tryAnchorClosure(
	candidate UntypedTask,
	availableTasks []UntypedTask,
	graphTaskMap map[string]UntypedTask,
	disabledTaskMap map[string]struct{},
) ([]UntypedTask, bool) {
	subgraphMap := make(map[string]UntypedTask)
	queue := []UntypedTask{candidate}
	subgraphMap[candidate.UntypedID().ReferenceIDString()] = candidate

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// If curr is in disabledTaskMap, anchor miss!
		if _, isDisabled := disabledTaskMap[curr.UntypedID().ReferenceIDString()]; isDisabled {
			return nil, false
		}

		hasMandatoryDeps := false
		for _, dep := range curr.Dependencies() {
			if dep.DescriptorCondition() == taskid.ConditionRequired && dep.DescriptorCardinality() == taskid.CardinalityPointToPoint {
				hasMandatoryDeps = true
				ptp, ok := dep.(taskid.PointToPointDescriptor)
				if !ok {
					continue
				}
				depRefID := ptp.ReferenceID()

				if _, isDisabled := disabledTaskMap[depRefID]; isDisabled {
					return nil, false // Branch targets a disabled task! Anchor miss!
				}

				if _, inGraph := graphTaskMap[depRefID]; inGraph {
					// Anchor Hit! Branch successfully anchors to graph.
					continue
				}

				if _, inSub := subgraphMap[depRefID]; !inSub {
					nextTask, err := findBestTaskImplementation(depRefID, availableTasks)
					if err != nil {
						return nil, false // Missing dependency
					}
					subgraphMap[depRefID] = nextTask
					queue = append(queue, nextTask)
				}
			}
		}

		// If a task has no mandatory dependencies and is not in graphTaskMap, it's an unanchored root task!
		if !hasMandatoryDeps {
			if _, inGraph := graphTaskMap[curr.UntypedID().ReferenceIDString()]; !inGraph {
				return nil, false
			}
		}
	}

	result := make([]UntypedTask, 0, len(subgraphMap))
	for _, task := range subgraphMap {
		result = append(result, task)
	}
	return result, true
}

// compareTaskByImplementationID compares two UntypedTasks lexicographically by their implementation ID string.
func compareTaskByImplementationID(a, b UntypedTask) int {
	return strings.Compare(a.UntypedID().String(), b.UntypedID().String())
}
