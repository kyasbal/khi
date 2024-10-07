export interface K8sResource {
  apiVersion: string;
  kind: string;
  metadata?: K8sMetadata;
}

export interface K8sControlledResource<
  Spec extends K8sSpec = K8sSpec,
  Status extends K8sStatus = K8sStatus,
> extends K8sResource {
  spec?: Spec;
  status?: Status;
}

export interface K8sMetadata {
  name: string;
  namespace?: string;
  uid?: string;
  ownerReferences?: K8sOwnerReference[];
  labels?: { [key: string]: string };
  annotations?: { [key: string]: string };
}

export interface K8sOwnerReference {
  apiVersion: string;
  kind: string;
  name: string;
  uid: string;
}

// TODO: The fields are not common in k8s resources
export interface K8sStatus {
  phase: string;
  conditions?: K8sCondition[];
  [label: string]: unknown;
}

export interface K8sSpec {
  [label: string]: unknown;
}

export interface K8sCondition {
  type: string;
  status: string;
  message: string;
  reason: string;
  lastTransitionTime: string;
  lastHeartbeatTime: string;
}

/*
 * Pod resource
 */

export type K8sPodResource = K8sControlledResource<PodSpec, PodStatus>;

export interface PodSpec extends K8sSpec {
  containers?: ContainerSpec[];
  initContainers?: ContainerSpec[];
  nodeName?: string;
}

export interface PodStatus extends K8sStatus {
  containerStatuses?: ContainerStatus[];
  initContainerStatuses?: ContainerStatus[];
  podIP?: string;
}

export interface ContainerSpec {
  name: string;
}

export interface ContainerStatus {
  name: string;
  ready: boolean;
  restartCount: number;
  started: boolean;
  state: ContainerStatusState;
}

export interface ContainerStatusState {
  running?: {
    startedAt: string;
  };
  terminated?: {
    exitCode: number;
    reason: string;
    finishedAt: string;
    containerID: string;
  };
}

/**
 * K8s node resource
 */

export type K8sNodeResource = K8sControlledResource<NodeSpec, NodeStatus>;

export interface NodeSpec extends K8sSpec {
  podCIDR?: string;
  taints?: {
    effect: string;
    key: string;
    timeAdded: string;
  }[];
}

export interface NodeStatus extends K8sStatus {
  addresses?: {
    address: string;
    type: string;
  }[];
  nodeInfo?: {
    architecture: string;
    bootID: string;
    containerRuntimeVersion: string;
    kernelVersion: string;
    kubeProxyVersion: string;
    kubeletVersion: string;
    machineID: string;
    operatingSystem: string;
    osImage: string;
    systemUUID: string;
  };
}

/**
 * K8s service resource
 */

export type K8sServiceResource = K8sControlledResource<
  ServiceSpec,
  ServiceStatus
>;

export interface ServiceSpec extends K8sSpec {
  type: 'NodePort' | 'ClusterIP' | 'Loadbalancer' | 'ExternalName';
  selector?: { [key: string]: string };
}

export interface ServiceStatus extends K8sStatus {
  clusterIp: string;
}

export interface K8sPodBindingResource extends K8sResource {
  target?: {
    kind: string;
    name: string;
  };
}
