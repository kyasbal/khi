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
	"slices"

	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
)

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

	sourceID := sourcePair.stage2.UntypedID().String()
	if leadsToFeedback {
		sourceID = sourcePair.stage1.UntypedID().String()
	}

	if !targetIsSplit {
		rewrittenEdge := e
		rewrittenEdge.SourceImplID = sourceID
		return []taskid.TaskEdge{rewrittenEdge}
	}

	targetPair := stageTasks[e.TargetImplID]
	e1 := e
	e1.SourceImplID = sourceID
	e1.TargetImplID = targetPair.stage1.UntypedID().String()
	e2 := e
	e2.SourceImplID = sourceID
	e2.TargetImplID = targetPair.stage2.UntypedID().String()
	return []taskid.TaskEdge{e1, e2}
}

// leadsToFeedbackProducer checks if targetID is a feedback producer or has a directed path leading to one.
func leadsToFeedbackProducer(
	targetID string,
	feedbackProducers map[string]bool,
	pointToPointOutgoing map[string][]string,
) bool {
	for feedbackProducerID := range feedbackProducers {
		if targetID == feedbackProducerID || isReachable(targetID, feedbackProducerID, pointToPointOutgoing) {
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
			for _, be := range bootstrap {
				rewritten := be
				if sourcePair, sourceIsSplit := stageTasks[be.SourceImplID]; sourceIsSplit {
					rewritten.SourceImplID = sourcePair.stage2.UntypedID().String()
				}
				resolvedFanInEdges = append(resolvedFanInEdges, rewritten)
			}
			continue
		}

		stage1ID := pair.stage1.UntypedID().String()
		stage2ID := pair.stage2.UntypedID().String()

		for _, be := range bootstrap {
			e1 := be
			e2 := be
			if sourcePair, sourceIsSplit := stageTasks[be.SourceImplID]; sourceIsSplit {
				e1.SourceImplID = sourcePair.stage2.UntypedID().String()
				e2.SourceImplID = sourcePair.stage2.UntypedID().String()
			}
			e1.TargetImplID = stage1ID
			e2.TargetImplID = stage2ID
			resolvedFanInEdges = append(resolvedFanInEdges, e1, e2)
		}

		for _, fe := range feedback {
			rewrittenEdge := fe
			if fe.SourceImplID == key.consumerImplID {
				rewrittenEdge.SourceImplID = stage1ID
			} else if sourcePair, sourceIsSplit := stageTasks[fe.SourceImplID]; sourceIsSplit {
				rewrittenEdge.SourceImplID = sourcePair.stage2.UntypedID().String()
			}
			rewrittenEdge.TargetImplID = stage2ID
			resolvedFanInEdges = append(resolvedFanInEdges, rewrittenEdge)
		}
	}
	return resolvedFanInEdges
}
