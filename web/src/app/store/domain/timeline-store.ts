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

/**
 * Base data structure representing a revision in a timeline.
 */
export interface RevisionData {
  readonly time: bigint;
  // ...other fields will be defined in the concrete implementation
}

/**
 * Base data structure representing a timeline with its revisions.
 */
export interface TimelineData<TRev extends RevisionData> {
  readonly id: number;
  readonly revisions: TRev[];
}

/**
 * Store for managing timelines and querying their revisions efficiently.
 */
export class TimelineStore<
  TRev extends RevisionData,
  TTimeline extends TimelineData<TRev>,
> {
  private readonly timelines = new Map<number, TTimeline>();

  /**
   * Adds a new timeline to the store.
   */
  add(timeline: TTimeline): void {
    this.timelines.set(timeline.id, timeline);
  }

  /**
   * Retrieves a timeline by its ID.
   */
  getTimeline(id: number): TTimeline | undefined {
    return this.timelines.get(id);
  }

  /**
   * Sorts revisions in all timelines by timestamp.
   * Should be called after all timelines and revisions are ingested.
   */
  sort(): void {
    for (const tl of this.timelines.values()) {
      tl.revisions.sort((a, b) =>
        a.time < b.time ? -1 : a.time > b.time ? 1 : 0,
      );
    }
  }

  /**
   * High performance binary search for fetching the latest revision at or before a given timestamp.
   */
  getLatestRevisionAt(timelineId: number, time: bigint): TRev | null {
    const tl = this.timelines.get(timelineId);
    if (!tl || tl.revisions.length === 0) return null;

    let low = 0;
    let high = tl.revisions.length - 1;
    let ans: TRev | null = null;

    while (low <= high) {
      const mid = Math.floor((low + high) / 2);
      if (tl.revisions[mid].time <= time) {
        ans = tl.revisions[mid];
        low = mid + 1; // Check if there is a later one that is still <= time
      } else {
        high = mid - 1;
      }
    }

    return ans;
  }
}
