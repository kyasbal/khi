/**
 * Copyright 2024 Google LLC
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
