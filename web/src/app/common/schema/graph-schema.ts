export function emptyGraphData(): GraphData {
  return {
    graphTime: '-',
    nodes: [],
    services: [],
    podOwners: {
      daemonset: [],
      job: [],
      replicaset: [],
    },
    podOwnerOwners: {
      cronjob: [],
      deployment: [],
    },
  };
}

export interface GraphResourceData {
  deletedAt?: string;
  updatedAt?: string;
}

export type PodOwnerKinds = 'daemonset' | 'replicaset' | 'job';
export type PodOwnerOwnerKinds = 'deployment' | 'cronjob';

export interface GraphData {
  graphTime: string;
  nodes: GraphNode[];
  services: ServiceGraphData[];
  podOwners: { [kind in PodOwnerKinds]: GraphPodOwner[] };
  podOwnerOwners: { [kind in PodOwnerOwnerKinds]: GraphPodOwnerOwner[] };
}

export interface LabeledGraphElement {
  labels: { [label_key: string]: string };
}

export interface ArchGraphCondition {
  type: string;
  message: string;
  status: string;
  is_positive_status: boolean;
}

export interface NamespacedArchGraphResource {
  uid?: string;
  name: string;
  namespace: string;
}

export interface GraphNode extends LabeledGraphElement, GraphResourceData {
  name: string;
  pods: PodGraphData[];
  podCIDR: string;
  externalIP: string;
  internalIP: string;
  taints: string[];
  conditions: ArchGraphCondition[];
}

export interface PodGraphData
  extends LabeledGraphElement,
    NamespacedArchGraphResource,
    GraphResourceData {
  containers: ContainerGraphData[];
  podIP: string;
  conditions: ArchGraphCondition[];
  phase: string;
  isPhaseHealthy: boolean;
  ownerUids: Set<string>;
}

export interface ContainerGraphData {
  name: string;
  isInitContainer: boolean;
  isStatusHealthy: boolean;
  status: string;
  reason: string;
  code: number;
  ready: boolean;
  statusReadFromManifest: boolean;
}

export interface ServiceGraphData
  extends NamespacedArchGraphResource,
    LabeledGraphElement,
    GraphResourceData {
  type: string;
  clusterIp: string;
  connectedPods: PodConnectionGraphData[];
}

export interface PodConnectionGraphData {
  node: GraphNode;
  pod: PodGraphData;
}

export interface PodOwnerConnectionGraphData {
  podOwner: GraphPodOwner;
}

export interface GraphPodOwnerBase
  extends NamespacedArchGraphResource,
    LabeledGraphElement,
    GraphResourceData {
  status: { [key: string]: unknown };
}

export interface GraphPodOwner extends GraphPodOwnerBase {
  ownerUids: Set<string>;
  connectedPods: PodConnectionGraphData[];
}

export interface GraphPodOwnerOwner extends GraphPodOwnerBase {
  connectedPodOwners: PodOwnerConnectionGraphData[];
}
