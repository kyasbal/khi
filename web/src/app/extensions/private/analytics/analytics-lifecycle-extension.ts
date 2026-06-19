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

import { inject } from '@angular/core';
import {
  LifecycleHookExtension,
  PageType,
} from 'src/app/extensions/extension-common/extension-types/lifecycle-hook';
import { FRONTEND_ANALYTICS, KHIAnalyticsActivityType } from './types';
import { randomString } from 'src/app/utils/random';
import { InspectionDataV2 } from 'src/app/store/domain/inspection-data';

/**
 * AnalyticsLifecycleExtension reports event on lifecycle events with the injected FRONTEND_ANALYTICS service.
 *
 */
export const AnalyticsLifecycleExtension: LifecycleHookExtension = {
  onPageLoaded: (page: PageType) => {
    const analytics = inject(FRONTEND_ANALYTICS);
    analytics.init(page);
  },
  onInspectionDataOpen: (
    inspectionData: InspectionDataV2,
    rawData: ArrayBuffer,
  ) => {
    const analytics = inject(FRONTEND_ANALYTICS);
    const reportInAsync = async () => {
      // analytics code to inspection data metrics
      const hash = await sha512FromArrayBuffer(rawData);
      const openId = randomString();
      analytics.globalMetadata()['inspectionDataHash'] = hash;
      analytics.globalMetadata()['openId'] = openId;

      analytics.report(KHIAnalyticsActivityType.OpenInspectionData, {
        logLength: inspectionData.logStore.logs.length,
        decompressedTextBufferLength:
          inspectionData.metadata?.header?.fileSize ?? 0,
        revisionCount: inspectionData.timelineStore.timelines.reduce(
          (prev, next) => next.revisions.length + prev,
          0,
        ),
        eventCount: inspectionData.timelineStore.timelines.reduce(
          (prev, next) => next.events.length + prev,
          0,
        ),
      });
    };
    reportInAsync();
  },

  onInspectionStart: () => {
    const analytics = inject(FRONTEND_ANALYTICS);
    analytics.report(KHIAnalyticsActivityType.Inspect, {});
  },
};

/**
 * Generate SHA-512 hash string from given ArrayBuffer
 * @param source source of the hash
 * @returns SHA-512 hash in hex-string
 */
async function sha512FromArrayBuffer(source: ArrayBuffer): Promise<string> {
  const hashBuffer = await crypto.subtle.digest('SHA-512', source);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map((b) => b.toString(16).padStart(2, '0')).join('');
}
