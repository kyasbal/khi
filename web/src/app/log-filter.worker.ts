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

/// <reference lib="webworker" />

import * as LogFilterWorker from './worker-types';

addEventListener('message', (data) => {
  const query = data.data as LogFilterWorker.FilterQuery;
  const regex = new RegExp(query.regexInStr);
  const result = [];
  for (let logIndex = 0; logIndex < query.logs.length; logIndex++) {
    if (!regex.test(query.logs[logIndex])) {
      // Retrieve only not match
      result.push(logIndex);
    }
  }
  postMessage({
    notMatch: result,
    taskId: query.taskId,
    isKHIWorkerPacket: true,
  } as LogFilterWorker.FilterResult);
});
