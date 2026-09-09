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
	"fmt"
	"slices"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
)

type consumerTagKey struct {
	consumerImplID string
	tag            string
}

// resolveFanInEdgesAndCycles resolves candidate fan-in edges against the fixed DAG formed by pointToPointEdges.
// It detects cycles caused by fan-in dependencies and applies priority-based pruning or returns descriptive fail-fast errors.
func resolveFanInEdgesAndCycles(
	graphTaskMap map[string]UntypedTask,
	pointToPointEdges []taskid.TaskEdge,
	candidateFanInEdges []taskid.TaskEdge,
) ([]taskid.TaskEdge, map[string][]string, map[string]map[string][]string, error) {
	refToTask, refToImplID, implToTask := buildTaskMaps(graphTaskMap)

	outgoing := make(map[string][]string, len(graphTaskMap))
	inDegree := make(map[string]int, len(graphTaskMap))
	for _, task := range graphTaskMap {
		implID := task.UntypedID().String()
		outgoing[implID] = nil
		inDegree[implID] = 0
	}

	for _, e := range pointToPointEdges {
		sourceImplID, ok := refToImplID[e.SourceRefID]
		if !ok {
			continue
		}
		outgoing[sourceImplID] = append(outgoing[sourceImplID], e.TargetID)
		inDegree[e.TargetID]++
	}

	if err := verifyAcyclic(outgoing, inDegree, graphTaskMap); err != nil {
		return nil, nil, nil, err
	}

	sortedKeys, edgesByKey := groupCandidateFanInEdgesByConsumerTag(candidateFanInEdges)
	var acceptedFanInEdges []taskid.TaskEdge

	for _, key := range sortedKeys {
		consumerTask := implToTask[key.consumerImplID]
		candidates := edgesByKey[key]

		accepted, err := resolveCandidateFanInForKey(
			consumerTask,
			key,
			candidates,
			outgoing,
			refToTask,
			refToImplID,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		acceptedFanInEdges = append(acceptedFanInEdges, accepted...)
	}

	allEdges := make([]taskid.TaskEdge, 0, len(pointToPointEdges)+len(acceptedFanInEdges))
	allEdges = append(allEdges, pointToPointEdges...)
	allEdges = append(allEdges, acceptedFanInEdges...)
	dedupedEdges := deduplicateAndNormalizeEdges(allEdges)

	boundFanInRefIDs, boundFanInRefIDsByTask := collectBoundFanInRefIDs(acceptedFanInEdges)
	return dedupedEdges, boundFanInRefIDs, boundFanInRefIDsByTask, nil
}

// buildTaskMaps constructs lookup maps for tasks and their identifiers.
func buildTaskMaps(graphTaskMap map[string]UntypedTask) (map[string]UntypedTask, map[string]string, map[string]UntypedTask) {
	refToTask := make(map[string]UntypedTask, len(graphTaskMap))
	refToImplID := make(map[string]string, len(graphTaskMap))
	implToTask := make(map[string]UntypedTask, len(graphTaskMap))
	for _, task := range graphTaskMap {
		refID := task.UntypedID().ReferenceIDString()
		implID := task.UntypedID().String()
		refToTask[refID] = task
		refToImplID[refID] = implID
		implToTask[implID] = task
	}
	return refToTask, refToImplID, implToTask
}

// groupCandidateFanInEdgesByConsumerTag groups candidate fan-in edges by (consumerImplID, tag) and sorts keys deterministically.
func groupCandidateFanInEdgesByConsumerTag(candidateFanInEdges []taskid.TaskEdge) ([]consumerTagKey, map[consumerTagKey][]taskid.TaskEdge) {
	edgesByKey := make(map[consumerTagKey][]taskid.TaskEdge)
	for _, e := range candidateFanInEdges {
		key := consumerTagKey{
			consumerImplID: e.TargetID,
			tag:            e.Tag,
		}
		edgesByKey[key] = append(edgesByKey[key], e)
	}

	keys := make([]consumerTagKey, 0, len(edgesByKey))
	for k := range edgesByKey {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b consumerTagKey) int {
		if a.consumerImplID != b.consumerImplID {
			return strings.Compare(a.consumerImplID, b.consumerImplID)
		}
		return strings.Compare(a.tag, b.tag)
	})
	return keys, edgesByKey
}

// resolveCandidateFanInForKey resolves candidate fan-in edges for a specific consumer and tag key.
func resolveCandidateFanInForKey(
	consumerTask UntypedTask,
	key consumerTagKey,
	candidates []taskid.TaskEdge,
	outgoing map[string][]string,
	refToTask map[string]UntypedTask,
	refToImplID map[string]string,
) ([]taskid.TaskEdge, error) {
	candidatesByPriority := make(map[int][]taskid.TaskEdge)
	for _, e := range candidates {
		candidatesByPriority[e.Priority] = append(candidatesByPriority[e.Priority], e)
	}

	priorities := make([]int, 0, len(candidatesByPriority))
	for p := range candidatesByPriority {
		priorities = append(priorities, p)
	}
	slices.Sort(priorities)

	var acceptedEdges []taskid.TaskEdge

	for _, priority := range priorities {
		candidateEdgesAtPriority := deduplicateCandidateEdgesAtPriority(candidatesByPriority[priority])

		for _, candidateEdge := range candidateEdgesAtPriority {
			producerImplID := refToImplID[candidateEdge.SourceRefID]
			createsCycle := (producerImplID == key.consumerImplID) || isReachable(key.consumerImplID, producerImplID, outgoing)

			if !createsCycle {
				outgoing[producerImplID] = append(outgoing[producerImplID], key.consumerImplID)
				acceptedEdges = append(acceptedEdges, candidateEdge)
				continue
			}

			hasAllowMultiStage := typedmap.GetOrDefault(consumerTask.Labels(), LabelKeyAllowMultiStageExecution, false)
			if !hasAllowMultiStage {
				return nil, fmt.Errorf("task %s requires multi-stage execution but lacks AllowMultiStageExecution label", consumerTask.UntypedID())
			}

			if len(candidateEdgesAtPriority) > 1 {
				return nil, buildAmbiguousPriorityError(candidateEdge, candidateEdgesAtPriority, refToTask, key.tag)
			}

			if len(acceptedEdges) == 0 {
				return nil, fmt.Errorf("task %s has circular fan-in dependency on tag %s with no higher-priority producer to bootstrap execution", consumerTask.UntypedID(), key.tag)
			}
			// Strictly lower priority producer with AllowMultiStageExecution and bootstrap producer: prune edge.
		}
	}
	return acceptedEdges, nil
}

// deduplicateCandidateEdgesAtPriority sorts and deduplicates candidate edges by SourceRefID within the same priority.
func deduplicateCandidateEdgesAtPriority(candidateEdges []taskid.TaskEdge) []taskid.TaskEdge {
	slices.SortFunc(candidateEdges, func(a, b taskid.TaskEdge) int {
		return strings.Compare(a.SourceRefID, b.SourceRefID)
	})
	uniqueCandidateEdges := make([]taskid.TaskEdge, 0, len(candidateEdges))
	seenSources := make(map[string]bool)
	for _, e := range candidateEdges {
		if !seenSources[e.SourceRefID] {
			seenSources[e.SourceRefID] = true
			uniqueCandidateEdges = append(uniqueCandidateEdges, e)
		}
	}
	return uniqueCandidateEdges
}

// buildAmbiguousPriorityError formats a descriptive error when competing producers have equal priority in a cycle.
func buildAmbiguousPriorityError(
	candidateEdge taskid.TaskEdge,
	candidateEdgesAtPriority []taskid.TaskEdge,
	refToTask map[string]UntypedTask,
	tag string,
) error {
	var otherEdge taskid.TaskEdge
	for _, other := range candidateEdgesAtPriority {
		if other.SourceRefID != candidateEdge.SourceRefID {
			otherEdge = other
			break
		}
	}
	producerImplIDA := candidateEdge.SourceRefID
	if producerTask, ok := refToTask[candidateEdge.SourceRefID]; ok {
		producerImplIDA = producerTask.UntypedID().String()
	}
	producerImplIDB := otherEdge.SourceRefID
	if producerTask, ok := refToTask[otherEdge.SourceRefID]; ok {
		producerImplIDB = producerTask.UntypedID().String()
	}
	if producerImplIDA > producerImplIDB {
		producerImplIDA, producerImplIDB = producerImplIDB, producerImplIDA
	}
	return fmt.Errorf("ambiguous FanIn priority between %s and %s for tag %s", producerImplIDA, producerImplIDB, tag)
}

// collectBoundFanInRefIDs aggregates unique source reference IDs for accepted fan-in edges globally and per consumer task.
func collectBoundFanInRefIDs(acceptedFanInEdges []taskid.TaskEdge) (map[string][]string, map[string]map[string][]string) {
	boundFanInRefIDs := make(map[string][]string)
	boundFanInRefIDsByTask := make(map[string]map[string][]string)
	for _, e := range acceptedFanInEdges {
		if e.Tag == "" {
			continue
		}
		boundFanInRefIDs[e.Tag] = append(boundFanInRefIDs[e.Tag], e.SourceRefID)
		if boundFanInRefIDsByTask[e.TargetID] == nil {
			boundFanInRefIDsByTask[e.TargetID] = make(map[string][]string)
		}
		boundFanInRefIDsByTask[e.TargetID][e.Tag] = append(boundFanInRefIDsByTask[e.TargetID][e.Tag], e.SourceRefID)
	}
	for tag, refIDs := range boundFanInRefIDs {
		slices.Sort(refIDs)
		boundFanInRefIDs[tag] = slices.Compact(refIDs)
	}
	for targetID, byTag := range boundFanInRefIDsByTask {
		for tag, refIDs := range byTag {
			slices.Sort(refIDs)
			boundFanInRefIDsByTask[targetID][tag] = slices.Compact(refIDs)
		}
	}
	return boundFanInRefIDs, boundFanInRefIDsByTask
}

// verifyAcyclic checks if the graph formed by outgoing contains any cycle.
func verifyAcyclic(outgoing map[string][]string, inDegree map[string]int, graphTaskMap map[string]UntypedTask) error {
	inDegreeCopy := make(map[string]int, len(inDegree))
	for k, v := range inDegree {
		inDegreeCopy[k] = v
	}
	queue := make([]string, 0, len(graphTaskMap))
	for implID, deg := range inDegreeCopy {
		if deg == 0 {
			queue = append(queue, implID)
		}
	}
	visitedCount := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visitedCount++
		for _, target := range outgoing[curr] {
			inDegreeCopy[target]--
			if inDegreeCopy[target] == 0 {
				queue = append(queue, target)
			}
		}
	}
	if visitedCount < len(graphTaskMap) {
		cyclicPath := extractCyclicDependencyPath(outgoing, inDegreeCopy)
		return fmt.Errorf("failed to sort as a runnable task graph. \n The graph contains cyclic dependency\n%s", cyclicPath)
	}
	return nil
}

// isReachable checks if there is a directed path of length >= 1 from startImplID to targetImplID.
func isReachable(startImplID, targetImplID string, outgoing map[string][]string) bool {
	visited := make(map[string]bool)
	queue := make([]string, 0, len(outgoing[startImplID]))
	for _, next := range outgoing[startImplID] {
		if next == targetImplID {
			return true
		}
		visited[next] = true
		queue = append(queue, next)
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, next := range outgoing[curr] {
			if next == targetImplID {
				return true
			}
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}

// extractCyclicDependencyPath identifies and formats the cyclic path in the graph.
func extractCyclicDependencyPath(
	outgoing map[string][]string,
	inDegree map[string]int,
) string {
	unresolved := make([]string, 0)
	for implID, deg := range inDegree {
		if deg > 0 {
			unresolved = append(unresolved, implID)
		}
	}
	slices.Sort(unresolved)
	if len(unresolved) == 0 {
		return ""
	}

	cycle := findCycle(outgoing, inDegree, unresolved)
	if len(cycle) == 0 {
		return fmt.Sprintf("... -> %s -> ...", unresolved[0])
	}

	// Format matching: "... -> tail] -> [cycle -> ...] -> [head -> ..."
	return fmt.Sprintf("... -> %s] -> [%s] -> [%s -> ...", cycle[len(cycle)-1], strings.Join(cycle, " -> "), cycle[0])
}

// findCycle performs DFS over unresolved nodes to detect and return a cycle path.
func findCycle(outgoing map[string][]string, inDegree map[string]int, unresolved []string) []string {
	visitState := make(map[string]int) // 0: unvisited, 1: visiting, 2: visited
	parent := make(map[string]string)
	var cycle []string

	var dfs func(currImplID string) bool
	dfs = func(currImplID string) bool {
		visitState[currImplID] = 1
		for _, nextImplID := range outgoing[currImplID] {
			if inDegree[nextImplID] <= 0 {
				continue
			}
			if visitState[nextImplID] == 1 {
				for curr := currImplID; curr != nextImplID && curr != ""; curr = parent[curr] {
					cycle = append(cycle, curr)
				}
				cycle = append(cycle, nextImplID)
				slices.Reverse(cycle)
				return true
			}
			if visitState[nextImplID] == 0 {
				parent[nextImplID] = currImplID
				if dfs(nextImplID) {
					return true
				}
			}
		}
		visitState[currImplID] = 2
		return false
	}

	for _, start := range unresolved {
		if visitState[start] == 0 && dfs(start) {
			break
		}
	}
	return cycle
}
