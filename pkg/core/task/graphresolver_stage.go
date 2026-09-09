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
	"context"
	"slices"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
)

// stageTask wraps an UntypedTask to provide a distinct stage implementation ID for multi-stage execution.
type stageTask struct {
	originalTask UntypedTask
	stageID      taskid.UntypedTaskImplementationID
}

var _ UntypedTask = (*stageTask)(nil)

// UntypedID implements UntypedTask.
func (s *stageTask) UntypedID() taskid.UntypedTaskImplementationID {
	return s.stageID
}

// Labels implements UntypedTask.
func (s *stageTask) Labels() *typedmap.ReadonlyTypedMap {
	return s.originalTask.Labels()
}

// Dependencies implements UntypedTask.
func (s *stageTask) Dependencies() []Dependency {
	return s.originalTask.Dependencies()
}

// UntypedRun implements UntypedTask.
func (s *stageTask) UntypedRun(ctx context.Context) (any, error) {
	return s.originalTask.UntypedRun(ctx)
}

type stageTaskPair struct {
	stage1 UntypedTask
	stage2 UntypedTask
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
	pointToPointOutgoing map[string][]string,
	implToTask map[string]UntypedTask,
) ([]UntypedTask, []taskid.TaskEdge, map[string]map[string][]string, error) {
	stageTasks := createStageTasks(feedbackProducersByConsumer, implToTask)
	allTasks := expandTasksWithStages(graphTaskMap, stageTasks)

	resolvedPointToPointEdges := reroutePointToPointEdges(
		pointToPointEdges,
		stageTasks,
		feedbackProducersByConsumer,
		pointToPointOutgoing,
	)

	resolvedFanInEdges := rerouteFanInEdges(
		sortedKeys,
		bootstrapEdgesByKey,
		feedbackEdgesByKey,
		stageTasks,
	)

	allEdges := make([]taskid.TaskEdge, 0, len(resolvedPointToPointEdges)+len(resolvedFanInEdges))
	allEdges = append(allEdges, resolvedPointToPointEdges...)
	allEdges = append(allEdges, resolvedFanInEdges...)
	dedupedEdges := deduplicateAndNormalizeEdges(allEdges)

	boundFanInRefIDsByTaskImpl := collectBoundFanInRefIDsByTaskImpl(resolvedFanInEdges)

	if err := verifyFinalDAG(allTasks, dedupedEdges); err != nil {
		return nil, nil, nil, err
	}

	return allTasks, dedupedEdges, boundFanInRefIDsByTaskImpl, nil
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
		finalOutgoing[e.SourceImplID] = append(finalOutgoing[e.SourceImplID], e.TargetImplID)
		finalInDegree[e.TargetImplID]++
	}
	return verifyAcyclic(finalOutgoing, finalInDegree, len(allTasks))
}

// reroutePointToPointEdges rewrites Point-to-Point edges when consumer tasks are split into stages.
func reroutePointToPointEdges(
	pointToPointEdges []taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
	feedbackProducersByConsumer map[string]map[string]bool,
	pointToPointOutgoing map[string][]string,
) []taskid.TaskEdge {
	var resolvedPointToPointEdges []taskid.TaskEdge
	for _, e := range pointToPointEdges {
		resolvedPointToPointEdges = append(resolvedPointToPointEdges, rerouteSinglePointToPointEdge(e, stageTasks, feedbackProducersByConsumer, pointToPointOutgoing)...)
	}

	splitConsumerIDs := make([]string, 0, len(stageTasks))
	for id := range stageTasks {
		splitConsumerIDs = append(splitConsumerIDs, id)
	}
	slices.Sort(splitConsumerIDs)

	for _, id := range splitConsumerIDs {
		pair := stageTasks[id]
		resolvedPointToPointEdges = append(resolvedPointToPointEdges, taskid.TaskEdge{
			SourceRefID:  pair.stage1.UntypedID().ReferenceIDString(),
			SourceImplID: pair.stage1.UntypedID().String(),
			TargetImplID: pair.stage2.UntypedID().String(),
			Kind:         taskid.EdgeKindOrderOnly,
			Condition:    taskid.ConditionRequired,
			Cardinality:  taskid.CardinalityPointToPoint,
		})
	}

	return resolvedPointToPointEdges
}

// rerouteSinglePointToPointEdge reroutes a point-to-point edge based on whether its source and target are split into stages.
func rerouteSinglePointToPointEdge(
	e taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
	feedbackProducersByConsumer map[string]map[string]bool,
	pointToPointOutgoing map[string][]string,
) []taskid.TaskEdge {
	_, sourceIsSplit := stageTasks[e.SourceImplID]
	_, targetIsSplit := stageTasks[e.TargetImplID]

	if !sourceIsSplit && !targetIsSplit {
		return []taskid.TaskEdge{e}
	}

	if !sourceIsSplit && targetIsSplit {
		targetPair := stageTasks[e.TargetImplID]
		e1 := e
		e1.TargetImplID = targetPair.stage1.UntypedID().String()
		e2 := e
		e2.TargetImplID = targetPair.stage2.UntypedID().String()
		return []taskid.TaskEdge{e1, e2}
	}

	sourcePair := stageTasks[e.SourceImplID]
	feedbackProducers := feedbackProducersByConsumer[e.SourceImplID]

	leadsToFeedback := leadsToFeedbackProducer(e.TargetImplID, feedbackProducers, pointToPointOutgoing)

	sourceImplID := sourcePair.stage2.UntypedID().String()
	if leadsToFeedback {
		sourceImplID = sourcePair.stage1.UntypedID().String()
	}

	if !targetIsSplit {
		rewrittenEdge := e
		rewrittenEdge.SourceImplID = sourceImplID
		return []taskid.TaskEdge{rewrittenEdge}
	}

	targetPair := stageTasks[e.TargetImplID]
	e1 := e
	e1.SourceImplID = sourceImplID
	e1.TargetImplID = targetPair.stage1.UntypedID().String()
	e2 := e
	e2.SourceImplID = sourceImplID
	e2.TargetImplID = targetPair.stage2.UntypedID().String()
	return []taskid.TaskEdge{e1, e2}
}

// leadsToFeedbackProducer checks if targetImplID is a feedback producer or has a directed path leading to one.
func leadsToFeedbackProducer(
	targetImplID string,
	feedbackProducers map[string]bool,
	pointToPointOutgoing map[string][]string,
) bool {
	for feedbackProducerImplID := range feedbackProducers {
		if targetImplID == feedbackProducerImplID || isReachable(targetImplID, feedbackProducerImplID, pointToPointOutgoing) {
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
		edges := rerouteFanInEdgesForConsumer(
			key,
			bootstrapEdgesByKey[key],
			feedbackEdgesByKey[key],
			stageTasks,
		)
		resolvedFanInEdges = append(resolvedFanInEdges, edges...)
	}
	return resolvedFanInEdges
}

// rerouteFanInEdgesForConsumer reroutes fan-in edges for a single consumer key based on whether it is split into stages.
func rerouteFanInEdgesForConsumer(
	key consumerTagKey,
	bootstrap []taskid.TaskEdge,
	feedback []taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
) []taskid.TaskEdge {
	pair, isSplit := stageTasks[key.consumerImplID]
	if !isSplit {
		return rerouteUnsplitFanInEdges(bootstrap, stageTasks)
	}

	return rerouteSplitFanInEdges(key.consumerImplID, pair, bootstrap, feedback, stageTasks)
}

// rerouteUnsplitFanInEdges routes bootstrap edges for a non-split consumer, rewiring split producers to stage-2.
func rerouteUnsplitFanInEdges(
	bootstrap []taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
) []taskid.TaskEdge {
	var result []taskid.TaskEdge
	for _, be := range bootstrap {
		rewritten := be
		if sourcePair, sourceIsSplit := stageTasks[be.SourceImplID]; sourceIsSplit {
			rewritten.SourceImplID = sourcePair.stage2.UntypedID().String()
		}
		result = append(result, rewritten)
	}
	return result
}

// rerouteSplitFanInEdges duplicates bootstrap edges to stage-1/stage-2 and routes feedback edges to stage-2.
func rerouteSplitFanInEdges(
	consumerImplID string,
	pair stageTaskPair,
	bootstrap []taskid.TaskEdge,
	feedback []taskid.TaskEdge,
	stageTasks map[string]stageTaskPair,
) []taskid.TaskEdge {
	stage1ImplID := pair.stage1.UntypedID().String()
	stage2ImplID := pair.stage2.UntypedID().String()

	var result []taskid.TaskEdge
	for _, be := range bootstrap {
		e1 := be
		e2 := be
		if sourcePair, sourceIsSplit := stageTasks[be.SourceImplID]; sourceIsSplit {
			e1.SourceImplID = sourcePair.stage2.UntypedID().String()
			e2.SourceImplID = sourcePair.stage2.UntypedID().String()
		}
		e1.TargetImplID = stage1ImplID
		e2.TargetImplID = stage2ImplID
		result = append(result, e1, e2)
	}

	for _, fe := range feedback {
		rewrittenEdge := fe
		if fe.SourceImplID == consumerImplID {
			rewrittenEdge.SourceImplID = stage1ImplID
		} else if sourcePair, sourceIsSplit := stageTasks[fe.SourceImplID]; sourceIsSplit {
			rewrittenEdge.SourceImplID = sourcePair.stage2.UntypedID().String()
		}
		rewrittenEdge.TargetImplID = stage2ImplID
		result = append(result, rewrittenEdge)
	}
	return result
}
