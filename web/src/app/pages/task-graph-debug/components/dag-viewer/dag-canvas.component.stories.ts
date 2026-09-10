/**
 * Copyright 2026 Google LLC
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

import { Meta, StoryObj } from '@storybook/angular';
import { TaskDependencyCardinality } from 'src/app/generated/api/v1/inspection_task_graph_pb';
import { DagCanvasComponent } from 'src/app/pages/task-graph-debug/components/dag-viewer/dag-canvas.component';
import {
  DagViewerEdge,
  DagViewerNode,
} from 'src/app/pages/task-graph-debug/components/dag-viewer/dag-viewer.model';

const sampleNodes: DagViewerNode[] = [
  {
    id: 'k8s/cluster-info-reader#abc1',
    referenceId: 'k8s/cluster-info-reader',
    isFeature: false,
    isInitialTask: true,
    topologicalOrder: 0,
    priority: 100,
    labels: { cluster: 'gke' },
  },
  {
    id: 'k8s/audit-log-collector#abc2',
    referenceId: 'k8s/audit-log-collector',
    isFeature: true,
    isInitialTask: false,
    topologicalOrder: 1,
    priority: 150,
    labels: { logType: 'audit' },
  },
  {
    id: 'k8s/container-log-collector#abc3',
    referenceId: 'k8s/container-log-collector',
    isFeature: true,
    isInitialTask: false,
    topologicalOrder: 1,
    priority: 120,
    labels: { logType: 'container' },
  },
  {
    id: 'k8s/pod-timeline-mapper#abc4',
    referenceId: 'k8s/pod-timeline-mapper',
    isFeature: false,
    isInitialTask: false,
    topologicalOrder: 2,
    priority: 100,
    labels: {},
  },
  {
    id: 'k8s/node-timeline-mapper#abc5',
    referenceId: 'k8s/node-timeline-mapper',
    isFeature: false,
    isInitialTask: false,
    topologicalOrder: 2,
    priority: 100,
    labels: {},
  },
  {
    id: 'k8s/timeline-aggregator#abc6',
    referenceId: 'k8s/timeline-aggregator',
    isFeature: false,
    isInitialTask: false,
    topologicalOrder: 3,
    priority: 200,
    labels: {},
  },
];

const sampleEdges: DagViewerEdge[] = [
  {
    id: 'e1',
    sourceId: 'k8s/cluster-info-reader#abc1',
    destinationId: 'k8s/audit-log-collector#abc2',
    sourceReferenceId: 'k8s/cluster-info-reader',
    cardinality: TaskDependencyCardinality.POINT_TO_POINT,
    tag: '',
    priority: 100,
  },
  {
    id: 'e2',
    sourceId: 'k8s/cluster-info-reader#abc1',
    destinationId: 'k8s/container-log-collector#abc3',
    sourceReferenceId: 'k8s/cluster-info-reader',
    cardinality: TaskDependencyCardinality.POINT_TO_POINT,
    tag: '',
    priority: 100,
  },
  {
    id: 'e3',
    sourceId: 'k8s/audit-log-collector#abc2',
    destinationId: 'k8s/pod-timeline-mapper#abc4',
    sourceReferenceId: 'k8s/audit-log-collector',
    cardinality: TaskDependencyCardinality.POINT_TO_POINT,
    tag: '',
    priority: 100,
  },
  {
    id: 'e4',
    sourceId: 'k8s/container-log-collector#abc3',
    destinationId: 'k8s/node-timeline-mapper#abc5',
    sourceReferenceId: 'k8s/container-log-collector',
    cardinality: TaskDependencyCardinality.POINT_TO_POINT,
    tag: '',
    priority: 100,
  },
  {
    id: 'e5',
    sourceId: 'k8s/pod-timeline-mapper#abc4',
    destinationId: 'k8s/timeline-aggregator#abc6',
    sourceReferenceId: 'k8s/pod-timeline-mapper',
    cardinality: TaskDependencyCardinality.FAN_IN,
    tag: 'timeline-fragments',
    priority: 100,
  },
  {
    id: 'e6',
    sourceId: 'k8s/node-timeline-mapper#abc5',
    destinationId: 'k8s/timeline-aggregator#abc6',
    sourceReferenceId: 'k8s/node-timeline-mapper',
    cardinality: TaskDependencyCardinality.FAN_IN,
    tag: 'timeline-fragments',
    priority: 100,
  },
];

const meta: Meta<DagCanvasComponent> = {
  title: 'TaskGraphDebug/DagViewer/DagCanvas',
  component: DagCanvasComponent,
  tags: ['autodocs'],
  args: {
    nodes: sampleNodes,
    edges: sampleEdges,
  },
};

export default meta;
type Story = StoryObj<DagCanvasComponent>;

export const Default: Story = {
  args: {},
};

export const Empty: Story = {
  args: {
    nodes: [],
    edges: [],
  },
};
