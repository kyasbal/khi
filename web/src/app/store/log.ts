import {
  KHIFileTextReference,
  KHILogAnnotation,
} from '../common/schema/khi-file-types';
import {
  LogType,
  LogTypeMetadata,
  Severity,
  SeverityMetadata,
} from '../generated';
import { TimelineEntry } from './timeline';

export class LogEntry {
  /**
   * Set of timelines relate to this log.
   */
  public relatedTimelines: Set<TimelineEntry> = new Set();

  public logTypeLabel = LogTypeMetadata[this.logType].label;

  public logSeverityLabel = SeverityMetadata[this.severity].label;

  constructor(
    public readonly logIndex: number,
    public readonly insertId: string,
    public readonly logType: LogType,
    public readonly severity: Severity,
    public readonly time: number,
    public readonly message: string,
    public readonly body: KHIFileTextReference,
    public readonly annotations: KHILogAnnotation[],
  ) {}

  public static clone(entry: LogEntry): LogEntry {
    return new LogEntry(
      entry.logIndex,
      entry.insertId,
      entry.logType,
      entry.severity,
      entry.time,
      entry.message,
      entry.body,
      entry.annotations,
    );
  }
}

/**
 * NullLog is just a placeholder for log reference when the resource status is not inferred from any logs.
 */
export const NullLog = new LogEntry(
  -1,
  '',
  LogType.LogTypeUnknown,
  Severity.SeverityUnknown,
  0,
  '',
  { offset: 0, len: 0, buffer: 0 },
  [],
);
