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
	implToTask := buildImplToTaskMap(graphTaskMap)

	outgoing, err := buildPtPOutgoingGraph(graphTaskMap, pointToPointEdges)
	if err != nil {
		return nil, nil, nil, err
	}

	// Keep a copy of the PtP outgoing graph for reachability checks during edge re-routing.
	ptpOutgoing := cloneOutgoingGraph(outgoing)

	sortedKeys, edgesByKey := groupCandidateFanInEdgesByConsumerTag(candidateFanInEdges)
	bootstrapEdgesByKey, feedbackEdgesByKey, feedbackProducersByConsumer, err := partitionCandidateFanInEdges(
		sortedKeys,
		edgesByKey,
		implToTask,
		outgoing,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(feedbackProducersByConsumer) == 0 {
		return assembleSingleStageResult(graphTaskMap, pointToPointEdges, sortedKeys, bootstrapEdgesByKey)
	}

	return expandMultiStageResult(
		graphTaskMap,
		pointToPointEdges,
		sortedKeys,
		bootstrapEdgesByKey,
		feedbackEdgesByKey,
		feedbackProducersByConsumer,
		ptpOutgoing,
		implToTask,
	)
}

// buildPtPOutgoingGraph constructs the adjacency and in-degree maps from point-to-point edges and validates acyclicity.
func buildPtPOutgoingGraph(
	graphTaskMap map[string]UntypedTask,
	pointToPointEdges []taskid.TaskEdge,
) (map[string][]string, error) {
	outgoing := make(map[string][]string, len(graphTaskMap))
	inDegree := make(map[string]int, len(graphTaskMap))
	for _, task := range graphTaskMap {
		implID := task.UntypedID().String()
		outgoing[implID] = nil
		inDegree[implID] = 0
	}

	for _, e := range pointToPointEdges {
		outgoing[e.SourceID] = append(outgoing[e.SourceID], e.TargetID)
		inDegree[e.TargetID]++
	}

	if err := verifyAcyclic(outgoing, inDegree, len(graphTaskMap)); err != nil {
		return nil, err
	}
	return outgoing, nil
}

// partitionCandidateFanInEdges partitions candidate fan-in edges into bootstrap edges and feedback edges per consumer and tag.
func partitionCandidateFanInEdges(
	sortedKeys []consumerTagKey,
	edgesByKey map[consumerTagKey][]taskid.TaskEdge,
	implToTask map[string]UntypedTask,
	outgoing map[string][]string,
) (
	bootstrapEdgesByKey map[consumerTagKey][]taskid.TaskEdge,
	feedbackEdgesByKey map[consumerTagKey][]taskid.TaskEdge,
	feedbackProducersByConsumer map[string]map[string]bool,
	err error,
) {
	bootstrapEdgesByKey = make(map[consumerTagKey][]taskid.TaskEdge)
	feedbackEdgesByKey = make(map[consumerTagKey][]taskid.TaskEdge)
	feedbackProducersByConsumer = make(map[string]map[string]bool)

	for _, key := range sortedKeys {
		consumerTask := implToTask[key.consumerImplID]
		candidates := edgesByKey[key]

		bootstrap, feedback, err := resolveCandidateFanInForKey(
			consumerTask,
			key,
			candidates,
			outgoing,
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
				feedbackProducersByConsumer[key.consumerImplID][fe.SourceID] = true
			}
		}
	}
	return bootstrapEdgesByKey, feedbackEdgesByKey, feedbackProducersByConsumer, nil
}

// assembleSingleStageResult builds the resolved tasks, deduplicated edges, and bound fan-in refs
// when no cycle is detected and no tasks need to be split into stages.
func assembleSingleStageResult(
	graphTaskMap map[string]UntypedTask,
	pointToPointEdges []taskid.TaskEdge,
	sortedKeys []consumerTagKey,
	bootstrapEdgesByKey map[consumerTagKey][]taskid.TaskEdge,
) ([]UntypedTask, []taskid.TaskEdge, map[string]map[string][]string, error) {
	var allFanInEdges []taskid.TaskEdge
	for _, key := range sortedKeys {
		allFanInEdges = append(allFanInEdges, bootstrapEdgesByKey[key]...)
	}
	allEdges := make([]taskid.TaskEdge, 0, len(pointToPointEdges)+len(allFanInEdges))
	allEdges = append(allEdges, pointToPointEdges...)
	allEdges = append(allEdges, allFanInEdges...)
	dedupedEdges := deduplicateAndNormalizeEdges(allEdges)

	boundFanInRefIDsByTask := collectBoundFanInRefIDsByTask(allFanInEdges)
	tasks := make([]UntypedTask, 0, len(graphTaskMap))
	for _, t := range graphTaskMap {
		tasks = append(tasks, t)
	}
	slices.SortFunc(tasks, compareTaskByImplementationID)
	return tasks, dedupedEdges, boundFanInRefIDsByTask, nil
}

// expandMultiStageResult expands cyclic consumers into stage-1 and stage-2 tasks,
// reroutes point-to-point and fan-in edges, and verifies acyclicity of the final stage DAG.
func expandMultiStageResult(
	graphTaskMap map[string]UntypedTask,
	pointToPointEdges []taskid.TaskEdge,
	sortedKeys []consumerTagKey,
	bootstrapEdgesByKey map[consumerTagKey][]taskid.TaskEdge,
	feedbackEdgesByKey map[consumerTagKey][]taskid.TaskEdge,
	feedbackProducersByConsumer map[string]map[string]bool,
	ptpOutgoing map[string][]string,
	implToTask map[string]UntypedTask,
) ([]UntypedTask, []taskid.TaskEdge, map[string]map[string][]string, error) {
	stageTasks := createStageTasks(feedbackProducersByConsumer, implToTask)
	allTasks := expandTasksWithStages(graphTaskMap, stageTasks)

	resolvedPtPEdges := reroutePointToPointEdges(
		pointToPointEdges,
		stageTasks,
		feedbackProducersByConsumer,
		ptpOutgoing,
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

	boundFanInRefIDsByTask := collectBoundFanInRefIDsByTask(resolvedFanInEdges)

	if err := verifyFinalDAG(allTasks, dedupedEdges); err != nil {
		return nil, nil, nil, err
	}

	return allTasks, dedupedEdges, boundFanInRefIDsByTask, nil
}

// createStageTasks creates stage-1 and stage-2 task pairs for each cyclic consumer.
func createStageTasks(
	feedbackProducersByConsumer map[string]map[string]bool,
	implToTask map[string]UntypedTask,
) map[string]stageTaskPair {
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
	return stageTasks
}

// expandTasksWithStages replaces split consumers with their stage-1 and stage-2 pairs, returning a sorted task list.
func expandTasksWithStages(
	graphTaskMap map[string]UntypedTask,
	stageTasks map[string]stageTaskPair,
) []UntypedTask {
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
	return allTasks
}

// verifyFinalDAG builds the final adjacency graph and validates acyclicity.
func verifyFinalDAG(allTasks []UntypedTask, dedupedEdges []taskid.TaskEdge) error {
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
	return verifyAcyclic(finalOutgoing, finalInDegree, len(allTasks))
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
) []taskid.TaskEdge {
	var resolvedPtPEdges []taskid.TaskEdge
	for _, e := range pointToPointEdges {
		resolvedPtPEdges = append(resolvedPtPEdges, rerouteSinglePtPEdge(e, stageTasks, feedbackProducersByConsumer, ptpOutgoing)...)
	}

	for _, pair := range stageTasks {
		resolvedPtPEdges = append(resolvedPtPEdges, taskid.TaskEdge{
			SourceRefID: pair.stage1.UntypedID().ReferenceIDString(),
			SourceID:    pair.stage1.UntypedID().String(),
			TargetID:    pair.stage2.UntypedID().String(),
			Kind:        taskid.EdgeKindOrderOnly,
			Condition:   taskid.ConditionRequired,
			Cardinality: taskid.CardinalityPointToPoint,
		})
	}

	return resolvedPtPEdges
}

// rerouteSinglePtPEdge reroutes a point-to-point edge based on whether its source and target are split into stages.
func rerouteSinglePtPEdge(
	e taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
	feedbackProducersByConsumer map[string]map[string]bool,
	ptpOutgoing map[string][]string,
) []taskid.TaskEdge {
	_, sourceIsSplit := stageTasks[e.SourceID]
	_, targetIsSplit := stageTasks[e.TargetID]

	if !sourceIsSplit && !targetIsSplit {
		return []taskid.TaskEdge{e}
	}

	if !sourceIsSplit && targetIsSplit {
		targetPair := stageTasks[e.TargetID]
		e1 := e
		e1.TargetID = targetPair.stage1.UntypedID().String()
		e2 := e
		e2.TargetID = targetPair.stage2.UntypedID().String()
		return []taskid.TaskEdge{e1, e2}
	}

	sourcePair := stageTasks[e.SourceID]
	feedbackProducers := feedbackProducersByConsumer[e.SourceID]

	leadsToFeedback := leadsToFeedbackProducer(e.TargetID, feedbackProducers, ptpOutgoing)

	sourceID := sourcePair.stage2.UntypedID().String()
	if leadsToFeedback {
		sourceID = sourcePair.stage1.UntypedID().String()
	}

	if !targetIsSplit {
		rewrittenEdge := e
		rewrittenEdge.SourceID = sourceID
		return []taskid.TaskEdge{rewrittenEdge}
	}

	targetPair := stageTasks[e.TargetID]
	e1 := e
	e1.SourceID = sourceID
	e1.TargetID = targetPair.stage1.UntypedID().String()
	e2 := e
	e2.SourceID = sourceID
	e2.TargetID = targetPair.stage2.UntypedID().String()
	return []taskid.TaskEdge{e1, e2}
}

// leadsToFeedbackProducer checks if targetID is a feedback producer or has a directed path leading to one.
func leadsToFeedbackProducer(
	targetID string,
	feedbackProducers map[string]bool,
	ptpOutgoing map[string][]string,
) bool {
	for feedbackProducerID := range feedbackProducers {
		if targetID == feedbackProducerID || isReachable(targetID, feedbackProducerID, ptpOutgoing) {
			return true
		}
	}
	return false
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

// buildImplToTaskMap constructs a lookup map from task implementation ID to task.
func buildImplToTaskMap(graphTaskMap map[string]UntypedTask) map[string]UntypedTask {
	implToTask := make(map[string]UntypedTask, len(graphTaskMap))
	for _, task := range graphTaskMap {
		implToTask[task.UntypedID().String()] = task
	}
	return implToTask
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
) (bootstrapEdges, feedbackEdges []taskid.TaskEdge, err error) {
	candidatesByPriority := make(map[int][]taskid.TaskEdge)
	for _, e := range candidates {
		candidatesByPriority[e.Priority] = append(candidatesByPriority[e.Priority], e)
	}

	priorities := make([]int, 0, len(candidatesByPriority))
	for p := range candidatesByPriority {
		priorities = append(priorities, p)
	}
	slices.Sort(priorities)

	for _, priority := range priorities {
		candidateEdgesAtPriority := deduplicateCandidateEdgesAtPriority(candidatesByPriority[priority])

		for _, candidateEdge := range candidateEdgesAtPriority {
			producerImplID := candidateEdge.SourceID
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
				return nil, nil, buildAmbiguousPriorityError(candidateEdge, candidateEdgesAtPriority, key.tag)
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
	tag string,
) error {
	var otherEdge taskid.TaskEdge
	for _, other := range candidateEdgesAtPriority {
		if other.SourceRefID != candidateEdge.SourceRefID {
			otherEdge = other
			break
		}
	}
	producerImplIDA := candidateEdge.SourceID
	producerImplIDB := otherEdge.SourceID
	if producerImplIDA > producerImplIDB {
		producerImplIDA, producerImplIDB = producerImplIDB, producerImplIDA
	}
	return fmt.Errorf("ambiguous FanIn priority between %s and %s for tag %s", producerImplIDA, producerImplIDB, tag)
}

// collectBoundFanInRefIDsByTask aggregates unique source reference IDs for accepted fan-in edges per consumer task implementation ID.
func collectBoundFanInRefIDsByTask(acceptedFanInEdges []taskid.TaskEdge) map[string]map[string][]string {
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
