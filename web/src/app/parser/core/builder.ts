import { StringPoolStore } from 'src/app/store/domain/string-pool-store';
import { LogStore, TimeStampedLog } from 'src/app/store/domain/log-store';
import {
  TimelineStore,
  RevisionData,
  TimelineData,
} from 'src/app/store/domain/timeline-store';

/**
 * ParsedKHIFile is the domain representation of the KHI file.
 * It holds highly optimized stores that Smart Components can use to quickly
 * query logs and timeline states (e.g. via binary search).
 */
export interface ParsedKHIFile {
  readonly stringPool: StringPoolStore;
  readonly logStore: LogStore<TimeStampedLog>;
  readonly timelineStore: TimelineStore<
    RevisionData,
    TimelineData<RevisionData>
  >;
  // Future: metadataStore, styleStore, etc.
}

/**
 * InspectionDataBuilder accumulates decoded data during the streaming phase.
 * It provides the context needed for cross-referencing (e.g. looking up a string by ID).
 */
export class InspectionDataBuilder {
  public readonly stringPool = new StringPoolStore();
  public readonly logStore = new LogStore<TimeStampedLog>();
  public readonly timelineStore = new TimelineStore<
    RevisionData,
    TimelineData<RevisionData>
  >();

  /**
   * Finalizes the stores (e.g., sort them by timestamp) and return the domain model.
   */
  build(): ParsedKHIFile {
    this.logStore.sort();
    this.timelineStore.sort();

    return {
      stringPool: this.stringPool,
      logStore: this.logStore,
      timelineStore: this.timelineStore,
    };
  }
}
