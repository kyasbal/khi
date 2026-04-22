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
 * Store for managing interned strings by their ID.
 */
export class StringPoolStore {
  private readonly pool = new Map<number, string>();

  /**
   * Registers a string with its ID.
   */
  add(id: number, value: string): void {
    this.pool.set(id, value);
  }

  /**
   * Resolves an interned string by its ID.
   */
  get(id: number): string {
    return this.pool.get(id) || '';
  }
}
