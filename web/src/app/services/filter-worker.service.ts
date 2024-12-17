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

import * as LogFilterWorker from '../worker/worker-types';
import { InspectionDataStoreService } from './inspection-data-store.service';
import { randomString } from '../utils/random';
import { LogEntry } from '../store/log';
import { forkJoin, map, mergeMap, Observable, of, Subject, take } from 'rxjs';
import { ReferenceResolverStore } from '../common/loader/reference-resolver';

/**
 * FilterWorkerService provides log filter feature with regex.
 * KHI use multiple WebWorker to filter them because filtering massive count of logs with regex in main thread cause browser page freezed.
 */
export class FilterWorkerService {
  /**
   * The count of worker pool.
   */
  private static LOG_FILTER_WORKER_POOL_COUNT: number = 8;

  /**
   * The maximum count of logs. Logs are divided into subsets not exceeding this count and distributed to multiple workers.
   */
  private static MAX_LOG_COUNT_PER_SINGLE_FILTER_SUBTASK: number = 5000;

  /**
   * Index of the next worker to be used.
   */
  private nextWorker: number = 0;

  /**
   * The list of workers.
   */
  private logFilterWorkers: Worker[] = [];

  /**
   * Map of subjects to receive the result from worker.
   *
   */
  private _taskCompletionHandler: { [taskId: string]: Subject<number[]> } = {};
  constructor(private dataStore: InspectionDataStoreService) {
    for (let i = 0; i < FilterWorkerService.LOG_FILTER_WORKER_POOL_COUNT; i++) {
      this.logFilterWorkers.push(
        new Worker(
          new URL('../worker/log-filter/log-filter.worker', import.meta.url),
        ),
      );
      this.logFilterWorkers[i].onmessage = (d) => this._onMessage(d);
    }
  }

  public filterLogs(
    allLogs: LogEntry[],
    regexInStr: string,
  ): Observable<Set<number>> {
    const bufferResolver = this.dataStore.textBufferSource.value;
    if (bufferResolver === null) {
      return of();
    }

    const taskResultSubject = new Subject<Set<number>>();
    const filteredIndexSet = new Set<number>();
    const filterSubRequests = new Subject<{ start: number; length: number }>();
    filterSubRequests
      .pipe(
        mergeMap(({ start, length }) =>
          this.requestFilterSubset(
            allLogs,
            bufferResolver,
            regexInStr,
            start,
            length,
          ),
        ),
      )
      .subscribe((subsetIndicesOfFilteredOutLogs) => {
        subsetIndicesOfFilteredOutLogs.forEach((index) =>
          filteredIndexSet.add(index),
        );
        taskResultSubject.next(filteredIndexSet);
      });

    for (
      let i = 0;
      i < allLogs.length;
      i += FilterWorkerService.MAX_LOG_COUNT_PER_SINGLE_FILTER_SUBTASK
    ) {
      filterSubRequests.next({
        start: i,
        length: Math.min(
          FilterWorkerService.MAX_LOG_COUNT_PER_SINGLE_FILTER_SUBTASK,
          allLogs.length - i,
        ),
      });
    }
    filterSubRequests.complete();
    return taskResultSubject;
  }

  private _onMessage(d: MessageEvent) {
    const filterResult = d.data;
    if (!LogFilterWorker.isKHIWorkerPacket(filterResult)) return;

    const filterResultTyped = filterResult as LogFilterWorker.FilterResult;
    if (filterResultTyped.taskId in this._taskCompletionHandler) {
      this._taskCompletionHandler[filterResultTyped.taskId].next(
        filterResultTyped.notMatch,
      );
      this._taskCompletionHandler[filterResultTyped.taskId].complete();
      delete this._taskCompletionHandler[filterResultTyped.taskId];
    } else {
      console.error(`Unknown task ID ${filterResultTyped.taskId}`);
    }
  }

  /**
   * Send subset of logs to the log filter WebWorker to get log ids not matching with given regex.
   */
  private requestFilterSubset(
    logs: LogEntry[],
    textLoader: ReferenceResolverStore,
    regexInStr: string,
    startIndex: number,
    length: number,
  ): Observable<number[]> {
    const taskId = randomString();
    const taskLoader = new Subject<number[]>();
    this._taskCompletionHandler[taskId] = taskLoader;

    // Get the body of the logs for subset of logs and send them to the worker.
    FilterWorkerServieUtil.logEntriesToFilterWorkerLogs(
      textLoader,
      logs.slice(startIndex, startIndex + length),
    ).subscribe((transferrableLog) => {
      this.logFilterWorkers[this.nextWorker].postMessage({
        taskId,
        regexInStr,
        logs: transferrableLog,
      });
    });

    this.nextWorker =
      (this.nextWorker + 1) % FilterWorkerService.LOG_FILTER_WORKER_POOL_COUNT;
    return taskLoader;
  }
}

/**
 * Provides utility functions used from FilterWorkerService.
 */
export class FilterWorkerServieUtil {
  /**
   * Convert an array of LogEntry to an array of LogFilterWorker.FilterWorkerLog with resolving its body text.
   */
  public static logEntriesToFilterWorkerLogs(
    referenceResolverStore: ReferenceResolverStore,
    logs: LogEntry[],
  ): Observable<LogFilterWorker.FilterWorkerLog[]> {
    return forkJoin(
      logs.map((l) =>
        referenceResolverStore.getText(l.body).pipe(
          map(
            (logBody) =>
              ({
                index: l.logIndex,
                logBody: logBody,
              }) as LogFilterWorker.FilterWorkerLog,
          ),
        ),
      ),
    ).pipe(take(1));
  }
}
