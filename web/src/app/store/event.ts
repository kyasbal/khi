import { LogType, Severity } from '../generated';

export class ResourceEvent {
  constructor(
    public logIndex: number,
    public ts: number,
    public logType: LogType,
    public logSeverity: Severity,
  ) {}

  public static clone(event: ResourceEvent): ResourceEvent {
    return new ResourceEvent(
      event.logIndex,
      event.ts,
      event.logType,
      event.logSeverity,
    );
  }
}
