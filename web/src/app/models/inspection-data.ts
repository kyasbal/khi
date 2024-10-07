import { InspectionMetadataHeader } from '../common/schema/metadata-types';
import { ParentRelationship } from '../generated';
import { LogEntry } from '../store/log';
import { TimelineEntry, TimelineLayer } from '../store/timeline';

export class TimelineRange {
  constructor(
    public begin: number,
    public end: number,
  ) {}

  get duration(): number {
    return this.end - this.begin;
  }
}

/**
 * The root viewmodel of inspection data
 */
export class InspectionData {
  public readonly namespaces: Set<string>;

  public readonly kinds: Set<string>;

  constructor(
    public readonly header: InspectionMetadataHeader,
    public readonly rawInspectionData: ArrayBuffer,
    public readonly range: TimelineRange,
    public readonly timelines: TimelineEntry[],
    public readonly relationships: Set<ParentRelationship>,
    public readonly logs: LogEntry[],
  ) {
    this.namespaces = this._getUniqueValueOfTimelineLayer(
      TimelineLayer.Namespace,
    );
    this.kinds = this._getUniqueValueOfTimelineLayer(TimelineLayer.Kind);
  }

  private _getUniqueValueOfTimelineLayer(layer: TimelineLayer): Set<string> {
    const result = new Set<string>();
    for (const t of this.timelines) {
      if (t.layer === layer) {
        result.add(t.getNameOfLayer(layer));
      }
    }
    return result;
  }
}
