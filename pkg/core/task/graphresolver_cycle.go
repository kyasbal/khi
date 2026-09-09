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
// When a candidate fan-in edge creates a cycle on a consumer task annotated with AllowMultiStageExecution,
// it splits the consumer task into multiple execution stages (stage-1 and stage-2) and re-routes edges.
func resolveFanInEdgesAndCycles(
	graphTaskMap map[string]UntypedTask,
	pointToPointEdges []taskid.TaskEdge,
	candidateFanInEdges []taskid.TaskEdge,
) ([]UntypedTask, []taskid.TaskEdge, map[string]map[string][]string, error) {
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

	if err := verifyAcyclic(outgoing, inDegree, len(graphTaskMap)); err != nil {
		return nil, nil, nil, err
	}

	// Keep a copy of the PtP outgoing graph for reachability checks during edge re-routing.
	ptpOutgoing := cloneOutgoingGraph(outgoing)

	sortedKeys, edgesByKey := groupCandidateFanInEdgesByConsumerTag(candidateFanInEdges)
	bootstrapEdgesByKey := make(map[consumerTagKey][]taskid.TaskEdge)
	feedbackEdgesByKey := make(map[consumerTagKey][]taskid.TaskEdge)
	feedbackProducersByConsumer := make(map[string]map[string]bool)

	for _, key := range sortedKeys {
		consumerTask := implToTask[key.consumerImplID]
		candidates := edgesByKey[key]

		bootstrap, feedback, err := resolveCandidateFanInForKey(
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
		bootstrapEdgesByKey[key] = bootstrap
		feedbackEdgesByKey[key] = feedback

		if len(feedback) > 0 {
			if feedbackProducersByConsumer[key.consumerImplID] == nil {
				feedbackProducersByConsumer[key.consumerImplID] = make(map[string]bool)
			}
			for _, fe := range feedback {
				producerImplID := refToImplID[fe.SourceRefID]
				feedbackProducersByConsumer[key.consumerImplID][producerImplID] = true
			}
		}
	}

	if len(feedbackProducersByConsumer) == 0 {
		var allFanInEdges []taskid.TaskEdge
		for _, key := range sortedKeys {
			allFanInEdges = append(allFanInEdges, bootstrapEdgesByKey[key]...)
		}
		allEdges := make([]taskid.TaskEdge, 0, len(pointToPointEdges)+len(allFanInEdges))
		allEdges = append(allEdges, pointToPointEdges...)
		allEdges = append(allEdges, allFanInEdges...)
		dedupedEdges := deduplicateAndNormalizeEdges(allEdges)

		boundFanInRefIDsByTask := collectBoundFanInRefIDs(allFanInEdges)
		tasks := make([]UntypedTask, 0, len(graphTaskMap))
		for _, t := range graphTaskMap {
			tasks = append(tasks, t)
		}
		slices.SortFunc(tasks, compareTaskByImplementationID)
		return tasks, dedupedEdges, boundFanInRefIDsByTask, nil
	}

	stageTasks := make(map[string]stageTaskPair, len(feedbackProducersByConsumer))
	for consumerImplID := range feedbackProducersByConsumer {
		consumerTask := implToTask[consumerImplID]
		stage1 := &stageTask{
			originalTask: consumerTask,
			stageID:      taskid.NewStageImplementationID(consumerTask.UntypedID(), 1),
		}
		stage2 := &stageTask{
			originalTask: consumerTask,
			stageID:      taskid.NewStageImplementationID(consumerTask.UntypedID(), 2),
		}
		stageTasks[consumerImplID] = stageTaskPair{stage1: stage1, stage2: stage2}
	}

	allTasks := make([]UntypedTask, 0, len(graphTaskMap)+len(stageTasks))
	for _, task := range graphTaskMap {
		implID := task.UntypedID().String()
		if pair, ok := stageTasks[implID]; ok {
			allTasks = append(allTasks, pair.stage1, pair.stage2)
		} else {
			allTasks = append(allTasks, task)
		}
	}
	slices.SortFunc(allTasks, compareTaskByImplementationID)

	resolvedPtPEdges := reroutePointToPointEdges(
		pointToPointEdges,
		stageTasks,
		feedbackProducersByConsumer,
		ptpOutgoing,
		implToTask,
	)

	resolvedFanInEdges := rerouteFanInEdges(
		sortedKeys,
		bootstrapEdgesByKey,
		feedbackEdgesByKey,
		stageTasks,
	)

	allEdges := make([]taskid.TaskEdge, 0, len(resolvedPtPEdges)+len(resolvedFanInEdges))
	allEdges = append(allEdges, resolvedPtPEdges...)
	allEdges = append(allEdges, resolvedFanInEdges...)
	dedupedEdges := deduplicateAndNormalizeEdges(allEdges)

	boundFanInRefIDsByTask := collectBoundFanInRefIDs(resolvedFanInEdges)

	finalOutgoing := make(map[string][]string, len(allTasks))
	finalInDegree := make(map[string]int, len(allTasks))
	for _, t := range allTasks {
		implID := t.UntypedID().String()
		finalOutgoing[implID] = nil
		finalInDegree[implID] = 0
	}
	for _, e := range dedupedEdges {
		finalOutgoing[e.SourceID] = append(finalOutgoing[e.SourceID], e.TargetID)
		finalInDegree[e.TargetID]++
	}
	if err := verifyAcyclic(finalOutgoing, finalInDegree, len(allTasks)); err != nil {
		return nil, nil, nil, err
	}

	return allTasks, dedupedEdges, boundFanInRefIDsByTask, nil
}

type stageTaskPair struct {
	stage1 UntypedTask
	stage2 UntypedTask
}

// reroutePointToPointEdges rewrites Point-to-Point edges when consumer tasks are split into stages.
func reroutePointToPointEdges(
	pointToPointEdges []taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
	feedbackProducersByConsumer map[string]map[string]bool,
	ptpOutgoing map[string][]string,
	implToTask map[string]UntypedTask,
) []taskid.TaskEdge {
	var resolvedPtPEdges []taskid.TaskEdge
	for _, e := range pointToPointEdges {
		_, sourceIsSplit := stageTasks[e.SourceID]
		_, targetIsSplit := stageTasks[e.TargetID]

		if !sourceIsSplit && !targetIsSplit {
			resolvedPtPEdges = append(resolvedPtPEdges, e)
			continue
		}

		if !sourceIsSplit && targetIsSplit {
			targetPair := stageTasks[e.TargetID]
			e1 := e
			e1.TargetID = targetPair.stage1.UntypedID().String()
			e2 := e
			e2.TargetID = targetPair.stage2.UntypedID().String()
			resolvedPtPEdges = append(resolvedPtPEdges, e1, e2)
			continue
		}

		sourcePair := stageTasks[e.SourceID]
		feedbackProducers := feedbackProducersByConsumer[e.SourceID]

		leadsToFeedback := false
		for feedbackProducerID := range feedbackProducers {
			if e.TargetID == feedbackProducerID || isReachable(e.TargetID, feedbackProducerID, ptpOutgoing) {
				leadsToFeedback = true
				break
			}
		}

		sourceID := sourcePair.stage2.UntypedID().String()
		if leadsToFeedback {
			sourceID = sourcePair.stage1.UntypedID().String()
		}

		if !targetIsSplit {
			rewrittenEdge := e
			rewrittenEdge.SourceID = sourceID
			resolvedPtPEdges = append(resolvedPtPEdges, rewrittenEdge)
		} else {
			targetPair := stageTasks[e.TargetID]
			e1 := e
			e1.SourceID = sourceID
			e1.TargetID = targetPair.stage1.UntypedID().String()
			e2 := e
			e2.SourceID = sourceID
			e2.TargetID = targetPair.stage2.UntypedID().String()
			resolvedPtPEdges = append(resolvedPtPEdges, e1, e2)
		}
	}

	for consumerImplID, pair := range stageTasks {
		consumerTask := implToTask[consumerImplID]
		resolvedPtPEdges = append(resolvedPtPEdges, taskid.TaskEdge{
			SourceRefID: consumerTask.UntypedID().ReferenceIDString(),
			SourceID:    pair.stage1.UntypedID().String(),
			TargetID:    pair.stage2.UntypedID().String(),
			Kind:        taskid.EdgeKindOrderOnly,
			Condition:   taskid.ConditionRequired,
			Cardinality: taskid.CardinalityPointToPoint,
		})
	}

	return resolvedPtPEdges
}

// rerouteFanInEdges duplicates bootstrap FanIn edges across stages and routes feedback edges to stage-2.
func rerouteFanInEdges(
	sortedKeys []consumerTagKey,
	bootstrapEdgesByKey map[consumerTagKey][]taskid.TaskEdge,
	feedbackEdgesByKey map[consumerTagKey][]taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
) []taskid.TaskEdge {
	var resolvedFanInEdges []taskid.TaskEdge
	for _, key := range sortedKeys {
		bootstrap := bootstrapEdgesByKey[key]
		feedback := feedbackEdgesByKey[key]

		pair, isSplit := stageTasks[key.consumerImplID]
		if !isSplit {
			resolvedFanInEdges = append(resolvedFanInEdges, bootstrap...)
			continue
		}

		stage1ID := pair.stage1.UntypedID().String()
		stage2ID := pair.stage2.UntypedID().String()

		for _, be := range bootstrap {
			e1 := be
			e1.TargetID = stage1ID
			e2 := be
			e2.TargetID = stage2ID
			resolvedFanInEdges = append(resolvedFanInEdges, e1, e2)
		}

		for _, fe := range feedback {
			rewrittenEdge := fe
			if fe.SourceID == key.consumerImplID {
				rewrittenEdge.SourceID = stage1ID
			}
			rewrittenEdge.TargetID = stage2ID
			resolvedFanInEdges = append(resolvedFanInEdges, rewrittenEdge)
		}
	}
	return resolvedFanInEdges
}

// cloneOutgoingGraph creates a deep copy of an outgoing adjacency list.
func cloneOutgoingGraph(outgoing map[string][]string) map[string][]string {
	cloned := make(map[string][]string, len(outgoing))
	for k, v := range outgoing {
		cloned[k] = slices.Clone(v)
	}
	return cloned
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

// resolveCandidateFanInForKey resolves candidate fan-in edges for a specific consumer and tag key,
// partitioning them into bootstrap edges (acyclic) and feedback edges (cyclic with AllowMultiStageExecution).
func resolveCandidateFanInForKey(
	consumerTask UntypedTask,
	key consumerTagKey,
	candidates []taskid.TaskEdge,
	outgoing map[string][]string,
	refToTask map[string]UntypedTask,
	refToImplID map[string]string,
) ([]taskid.TaskEdge, []taskid.TaskEdge, error) {
	candidatesByPriority := make(map[int][]taskid.TaskEdge)
	for _, e := range candidates {
		candidatesByPriority[e.Priority] = append(candidatesByPriority[e.Priority], e)
	}

	priorities := make([]int, 0, len(candidatesByPriority))
	for p := range candidatesByPriority {
		priorities = append(priorities, p)
	}
	slices.Sort(priorities)

	var bootstrapEdges []taskid.TaskEdge
	var feedbackEdges []taskid.TaskEdge

	for _, priority := range priorities {
		candidateEdgesAtPriority := deduplicateCandidateEdgesAtPriority(candidatesByPriority[priority])

		for _, candidateEdge := range candidateEdgesAtPriority {
			producerImplID := refToImplID[candidateEdge.SourceRefID]
			createsCycle := (producerImplID == key.consumerImplID) || isReachable(key.consumerImplID, producerImplID, outgoing)

			if !createsCycle {
				outgoing[producerImplID] = append(outgoing[producerImplID], key.consumerImplID)
				bootstrapEdges = append(bootstrapEdges, candidateEdge)
				continue
			}

			hasAllowMultiStage := typedmap.GetOrDefault(consumerTask.Labels(), LabelKeyAllowMultiStageExecution, false)
			if !hasAllowMultiStage {
				return nil, nil, fmt.Errorf("task %s requires multi-stage execution but lacks AllowMultiStageExecution label", consumerTask.UntypedID())
			}

			if len(candidateEdgesAtPriority) > 1 {
				return nil, nil, buildAmbiguousPriorityError(candidateEdge, candidateEdgesAtPriority, refToTask, key.tag)
			}

			if len(bootstrapEdges) == 0 {
				return nil, nil, fmt.Errorf("task %s has circular fan-in dependency on tag %s with no higher-priority producer to bootstrap execution", consumerTask.UntypedID(), key.tag)
			}
			feedbackEdges = append(feedbackEdges, candidateEdge)
		}
	}
	return bootstrapEdges, feedbackEdges, nil
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

// collectBoundFanInRefIDs aggregates unique source reference IDs for accepted fan-in edges per consumer task implementation ID.
func collectBoundFanInRefIDs(acceptedFanInEdges []taskid.TaskEdge) map[string]map[string][]string {
	boundFanInRefIDsByTask := make(map[string]map[string][]string)
	for _, e := range acceptedFanInEdges {
		if e.Tag == "" {
			continue
		}
		if boundFanInRefIDsByTask[e.TargetID] == nil {
			boundFanInRefIDsByTask[e.TargetID] = make(map[string][]string)
		}
		boundFanInRefIDsByTask[e.TargetID][e.Tag] = append(boundFanInRefIDsByTask[e.TargetID][e.Tag], e.SourceRefID)
	}
	for targetID, byTag := range boundFanInRefIDsByTask {
		for tag, refIDs := range byTag {
			slices.Sort(refIDs)
			boundFanInRefIDsByTask[targetID][tag] = slices.Compact(refIDs)
		}
	}
	return boundFanInRefIDsByTask
}

// verifyAcyclic checks if the graph formed by outgoing contains any cycle.
func verifyAcyclic(outgoing map[string][]string, inDegree map[string]int, nodeCount int) error {
	inDegreeCopy := make(map[string]int, len(inDegree))
	for k, v := range inDegree {
		inDegreeCopy[k] = v
	}
	queue := make([]string, 0, nodeCount)
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
	if visitedCount < nodeCount {
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
