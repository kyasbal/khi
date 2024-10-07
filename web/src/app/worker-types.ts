/* eslint-disable-next-line @typescript-eslint/no-explicit-any */
export function isKHIWorkerPacket(packet: any): packet is KHIWorkerPacket {
  return 'isKHIWorkerPacket' in packet && packet['isKHIWorkerPacket'];
}

export interface KHIWorkerPacket {
  isKHIWorkerPacket: boolean;
}

export interface FilterQuery extends KHIWorkerPacket {
  taskId: string;
  regexInStr: string;
  logs: string[];
}

export interface FilterResult extends KHIWorkerPacket {
  taskId: string;
  notMatch: number[];
}
