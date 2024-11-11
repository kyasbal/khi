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

import * as LogFilterWorker from '../worker-types';
import { InspectionDataStoreService } from './inspection-data-store.service';
import { randomString } from '../utils/random';
import { LogEntry } from '../store/log';

export class FilterWorkerService {
  private logFilterWorker: Worker;
  private _taskCompletionHandler: { [taskId: string]: (a: unknown) => void } =
    {};
  constructor(private dataStore: InspectionDataStoreService) {
    this.logFilterWorker = new Worker(
      new URL('../log-filter.worker', import.meta.url),
    );
    this.logFilterWorker.onmessage = (d) => {
      const filterResult = d.data;
      if (!LogFilterWorker.isKHIWorkerPacket(filterResult)) return;
      const filterResultTyped = filterResult as LogFilterWorker.FilterResult;
      if (filterResultTyped.taskId in this._taskCompletionHandler) {
        this._taskCompletionHandler[filterResultTyped.taskId](
          new Set(filterResultTyped.notMatch),
        );
      } else {
        console.error(`Unknown task ID ${filterResultTyped.taskId}`);
      }
    };
  }

  public filterLogs(
    allLogs: LogEntry[],
    regexInStr: string,
  ): Promise<Set<number>> {
    const taskId = randomString();
    const bufferResolver = this.dataStore.textBufferSource.value;
    if (bufferResolver === null) {
      return Promise.resolve(new Set());
    }
    const allLogBodies = allLogs.map((l) => bufferResolver.getText(l.body));
    return new Promise((resolve) => {
      this._taskCompletionHandler[taskId] = resolve as (value: unknown) => void;
      this.logFilterWorker.postMessage({
        taskId,
        regexInStr,
        logs: allLogBodies,
      } as LogFilterWorker.FilterQuery);
    });
  }
}
