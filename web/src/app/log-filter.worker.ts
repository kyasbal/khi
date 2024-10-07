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
