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
 * Base abstraction for logs to support generic time-based queries.
 */
export interface TimeStampedLog {
  readonly time: bigint;
  // Other domain specific fields will be attached by the specific Log implementation
}

/**
 * Store for managing and querying logs efficiently.
 */
export class LogStore<TLog extends TimeStampedLog> {
  private readonly logs: TLog[] = [];

  /**
   * Adds a new log to the store.
   */
  add(log: TLog): void {
    this.logs.push(log);
  }

  /**
   * Sorts logs by timestamp. Should be called after all logs are added.
   */
  sort(): void {
    this.logs.sort((a, b) => (a.time < b.time ? -1 : a.time > b.time ? 1 : 0));
  }

  /**
   * Gets all logs in the given time range [start, end].
   * Optimized with binary search.
   */
  getLogsInRange(start: bigint, end: bigint): TLog[] {
    const startIndex = this.binarySearchFirstAfterOrEqual(start);
    if (startIndex === -1) return [];

    const endIndex = this.binarySearchFirstAfter(end);

    return this.logs.slice(
      startIndex,
      endIndex === -1 ? this.logs.length : endIndex,
    );
  }

  private binarySearchFirstAfterOrEqual(time: bigint): number {
    let low = 0;
    let high = this.logs.length - 1;
    let ans = -1;
    while (low <= high) {
      const mid = Math.floor((low + high) / 2);
      if (this.logs[mid].time >= time) {
        ans = mid;
        high = mid - 1;
      } else {
        low = mid + 1;
      }
    }
    return ans;
  }

  private binarySearchFirstAfter(time: bigint): number {
    let low = 0;
    let high = this.logs.length - 1;
    let ans = -1;
    while (low <= high) {
      const mid = Math.floor((low + high) / 2);
      if (this.logs[mid].time > time) {
        ans = mid;
        high = mid - 1;
      } else {
        low = mid + 1;
      }
    }
    return ans;
  }
}
