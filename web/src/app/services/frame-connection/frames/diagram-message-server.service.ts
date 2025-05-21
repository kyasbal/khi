/**
 * Copyright 2025 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { inject, Injectable } from '@angular/core';
import { WindowConnectorService } from '../window-connector.service';
import { map } from 'rxjs';
import {
  QUERY_CURRENT_INSPECTION_METADATA,
  QUERY_DIAGRAM_DATA,
} from 'src/app/common/schema/inter-window-messages';
import { InspectionDataStoreService } from '../../inspection-data-store.service';
import { InspectionData } from 'src/app/store/inspection-data';
import {
  BasicDiagramNamespacedElement,
  ContainerDiagramElement,
  ContainerStatus,
  ContainerType,
  DiagramElementType,
  DiagramFrame,
  DiagramPathType,
  NodeDiagramElement,
  PodDiagramElement,
} from '../../diagram/diagram-model-types';
import { ResourceTimeline, TimelineLayer } from 'src/app/store/timeline';
import { ResourceRevision } from 'src/app/store/revision';
import { RevisionVerb } from 'src/app/generated';

/**
 * Service that handles RPC communication for diagram data
 *
 * Runs on the main application page and responds to requests from diagram frames.
 * Provides inspection metadata and diagram model data to visualization components.
 */
@Injectable({ providedIn: 'root' })
export class DiagramMessageServer {
  /**
   * Communication service for handling inter-frame RPC
   */
  private readonly connector = inject(WindowConnectorService);

  /**
   * Data store containing the current inspection information
   */
  private readonly dataStore = inject(InspectionDataStoreService);

  /**
   * Registers RPC handlers for diagram-related queries
   *
   * Sets up handlers to respond to:
   * - Inspection metadata requests (time range)
   * - Diagram data requests for specific timestamps
   */
  public start() {
    this.connector.serveRPC(QUERY_CURRENT_INSPECTION_METADATA, () => {
      return this.dataStore.inspectionData.pipe(
        map((i) => {
          return i === null
            ? null
            : {
                // when no inspection data opened, return null as the response.
                startTime: i.range.begin,
                endTime: i.range.end,
              };
        }),
      );
    });

    this.connector.serveRPC(QUERY_DIAGRAM_DATA, (req) => {
      return this.dataStore.inspectionData.pipe(
        map((i) => {
          return i === null
            ? null
            : {
                // when no inspection data opened, return null as the response.
                model: getDiagramFrame(i, req.ts),
              };
        }),
      );
    });
  }
}

/**
 * Generates a diagram frame for the specified timestamp
 *
 * @param inspectionData - The current inspection data
 * @param time - Target timestamp for the diagram frame
 * @returns DiagramFrame representing the state at the specified time
 */
function getDiagramFrame(
  inspectionData: InspectionData,
  time: number,
): DiagramFrame {
  const result: DiagramFrame = {
    ts: time,
    nodes: [],
    lower: [],
    lowerPaths: [],
    upper: [],
    upperPaths: [],
  };
  generateNodeDiagramElementForFrame(inspectionData, time, result);
  return result;
}

function generateNodeDiagramElementForFrame(
  inspectionData: InspectionData,
  time: number,
  output: DiagramFrame,
) {
  const nodeTimelines = inspectionData.timelines.filter(
    (timeline) =>
      timeline.layer === TimelineLayer.Name &&
      timeline.getNameOfLayer(TimelineLayer.Namespace) === 'cluster-scope' &&
      timeline.getNameOfLayer(TimelineLayer.Kind) === 'node',
  );

  const nodeNameToPodTimelines: { [key: string]: ResourceTimeline[] } = {};
  const podTimelines = inspectionData.timelines.filter(
    (timeline) =>
      timeline.layer === TimelineLayer.Name &&
      timeline.getNameOfLayer(TimelineLayer.Kind) === 'pod',
  );

  podTimelines.forEach((pod) => {
    const currentRevision = pod.getLatestRevisionOfTime(time);
    if (currentRevision && !isDeletedPod(currentRevision)) {
      let nodeName = getNodeNameFromPodManifest(currentRevision.parsedManifest);
      if (!nodeName) {
        // try getting node name from the binding resource instead
        const bindingTimeline = inspectionData.getTimelineByResourcePath(
          pod.resourcePath + '#binding',
        );
        if (bindingTimeline) {
          const currentBindingTimelineRevision =
            bindingTimeline.getLatestRevisionOfTime(time);
          if (currentBindingTimelineRevision) {
            nodeName = getNodeNameFromPodBindingManifest(
              currentBindingTimelineRevision.parsedManifest,
            );
          }
        }
      }
      if (nodeName) {
        if (!nodeNameToPodTimelines[nodeName]) {
          nodeNameToPodTimelines[nodeName] = [];
        }
        nodeNameToPodTimelines[nodeName].push(pod);
      }
    }
  });

  output.nodes = nodeTimelines
    .map((timeline) => {
      const nodeName = timeline.getNameOfLayer(TimelineLayer.Name);
      const podTimelines = nodeNameToPodTimelines[nodeName] ?? [];
      const currentNodeManifest = timeline.getLatestRevisionOfTime(time);
      if (!currentNodeManifest || isDeletedResource(currentNodeManifest)) {
        return null;
      }

      podTimelines.forEach((podTimeline) => {
        const podTimeineRevision = podTimeline.getLatestRevisionOfTime(time);
        if (podTimeineRevision && !isDeletedPod(podTimeineRevision)) {
          recursivelyTrackOwnerResources(
            inspectionData,
            podTimeline.resourcePath,
            time,
            podTimeineRevision.parsedManifest,
            false,
            0,
            output,
          );
        }
      });

      return {
        id: timeline.resourcePath,
        type: DiagramElementType.Node,
        name: nodeName,
        pods: podTimelines
          .map((timeline) => generatePodDiagramElementForFrame(timeline, time))
          .filter((pod) => !!pod),
      } as NodeDiagramElement;
    })
    .filter((timeline) => timeline !== null);
}

function generatePodDiagramElementForFrame(
  podTimeline: ResourceTimeline,
  time: number,
): PodDiagramElement | null {
  const currentRevision = podTimeline.getLatestRevisionOfTime(time);
  if (!currentRevision) {
    return null;
  }

  return {
    id: podTimeline.resourcePath,
    type: DiagramElementType.Pod,
    name: podTimeline.getNameOfLayer(TimelineLayer.Name),
    namespace: podTimeline.getNameOfLayer(TimelineLayer.Namespace),
    containers: generateContainerDiagramElementsForFrame(
      podTimeline.resourcePath,
      currentRevision.parsedManifest,
    ),
    phase: getPodPhase(currentRevision.parsedManifest) || 'Unknown',
  };
}

/* eslint-disable @typescript-eslint/no-explicit-any */

function getNodeNameFromPodManifest(manifest: any): string | undefined {
  return manifest?.spec?.nodeName;
}

function getNodeNameFromPodBindingManifest(manifest: any): string | undefined {
  return manifest?.target?.name;
}

function getPodPhase(manifest: any): string | undefined {
  return manifest?.status?.phase;
}

function recursivelyTrackOwnerResources(
  inspectionData: InspectionData,
  from: string,
  time: number,
  resourceManifest: any,
  toLower: boolean,
  layerIndex: number,
  output: DiagramFrame,
) {
  const ownerTimelines = getOwnerTimelines(resourceManifest, inspectionData);
  ownerTimelines.forEach((timeline) => {
    const currentRevision = timeline.getLatestRevisionOfTime(time);
    if (currentRevision && !isDeletedResource(currentRevision)) {
      if (toLower) {
        if (output.lower.length === layerIndex) {
          output.lower.push([]);
          output.lowerPaths.push([]);
        }
      } else {
        if (output.upper.length === layerIndex) {
          output.upper.push([]);
          output.upperPaths.push([]);
        }
      }

      const targetLayer = toLower
        ? output.lower[layerIndex]
        : output.upper[layerIndex];
      const targetPathLayer = toLower
        ? output.lowerPaths[layerIndex]
        : output.upperPaths[layerIndex];

      const namespace = timeline.getNameOfLayer(TimelineLayer.Namespace);
      let isAlreadyContained =
        targetLayer.filter((e) => e.id === timeline.resourcePath).length > 0;
      if (!isAlreadyContained) {
        targetLayer.push({
          id: timeline.resourcePath,
          type: DiagramElementType.ReplicaSet,
          name: timeline.getNameOfLayer(TimelineLayer.Name),
          namespace,
        } as BasicDiagramNamespacedElement);
      }
      let isAlreadyBound =
        targetPathLayer.filter(
          (e) => e.outerID === timeline.resourcePath && e.innerID === from,
        ).length > 0;
      if (!isAlreadyBound) {
        targetPathLayer.push({
          type: DiagramPathType.Owner,
          innerID: from,
          outerID: timeline.resourcePath,
        });
      }
      recursivelyTrackOwnerResources(
        inspectionData,
        timeline.resourcePath,
        time,
        currentRevision.parsedManifest,
        toLower,
        layerIndex + 1,
        output,
      );
    }
  });
}

function getOwnerTimelines(
  manifest: any,
  inspectionData: InspectionData,
): ResourceTimeline[] {
  const ownerReferences = manifest?.metadata?.ownerReferences;
  if (!ownerReferences) {
    return [];
  }
  const result: ResourceTimeline[] = [];
  const ownNamespace = manifest?.metadata?.namespace ?? 'cluster-scope';
  for (const ownerReference of ownerReferences) {
    let apiVersion = ownerReference.apiVersion;
    if (apiVersion === 'v1') {
      apiVersion = 'core/v1';
    }
    const kind = ownerReference.kind.toLowerCase();
    const name = ownerReference.name;
    const ownerResourcePath = `${apiVersion}#${kind}#${ownNamespace}#${name}`;
    const ownerTimeline =
      inspectionData.getTimelineByResourcePath(ownerResourcePath);
    if (ownerTimeline) {
      result.push(ownerTimeline);
    }
  }
  return result;
}

function isDeletedPod(revision: ResourceRevision): boolean {
  const deletionGracePeriodSeconds = (revision.parsedManifest as any)?.metadata
    ?.deletionGracePeriodSeconds;
  if (deletionGracePeriodSeconds === undefined) {
    return (
      revision.lastMutationVerb === RevisionVerb.RevisionVerbDelete ||
      revision.lastMutationVerb === RevisionVerb.RevisionVerbDeleteCollection
    );
  } else {
    return deletionGracePeriodSeconds === 0;
  }
}

function isDeletedResource(revision: ResourceRevision): boolean {
  return (
    revision.lastMutationVerb === RevisionVerb.RevisionVerbDelete ||
    revision.lastMutationVerb === RevisionVerb.RevisionVerbDeleteCollection
  );
}

function generateContainerDiagramElementsForFrame(
  podResourcePath: string,
  podManifest: any,
): ContainerDiagramElement[] {
  const containers: ContainerDiagramElement[] = [];

  const initContainerStatuses = podManifest?.status?.initContainerStatuses;
  if (initContainerStatuses) {
    for (const initContainerStatus of initContainerStatuses) {
      containers.push(
        generateContainerDiagramElementFromContainerStatus(
          podResourcePath + '-init-' + initContainerStatus.name,
          ContainerType.Init,
          initContainerStatus,
        ),
      );
    }
  }

  const containerStatuses = podManifest?.status?.containerStatuses;
  if (containerStatuses) {
    for (const containerStatus of containerStatuses) {
      containers.push(
        generateContainerDiagramElementFromContainerStatus(
          podResourcePath + '-standard-' + containerStatus.name,
          ContainerType.Standard,
          containerStatus,
        ),
      );
    }
  }

  return containers;
}

function generateContainerDiagramElementFromContainerStatus(
  id: string,
  containerType: ContainerType,
  containerStatusManifest: any,
): ContainerDiagramElement {
  let containerStatus = ContainerStatus.Unknown;
  let exitCode = -1;
  let reason = '';
  if (containerStatusManifest?.state?.running) {
    if (containerStatusManifest?.ready) {
      containerStatus = ContainerStatus.Ready;
    } else {
      containerStatus = ContainerStatus.NotReady;
    }
  } else if (containerStatusManifest?.state?.terminated) {
    exitCode = containerStatusManifest.state.terminated.exitCode;
    reason = containerStatusManifest.state.terminated.reason;
    if (exitCode === 0) {
      containerStatus = ContainerStatus.TerminatedSuccess;
    } else {
      containerStatus = ContainerStatus.TerminatedFailure;
    }
  } else if (containerStatusManifest?.state?.waiting) {
    containerStatus = ContainerStatus.Waiting;
    reason = containerStatusManifest.state.waiting.reason;
  }
  return {
    id: id,
    type: DiagramElementType.Container,
    name: containerStatusManifest.name,
    containerType: containerType,
    status: containerStatus,
    exitCode: exitCode,
    reason: reason,
    restartCount: containerStatusManifest.restartCount ?? 0,
  };
}

/* eslint-enable @typescript-eslint/no-explicit-any */
