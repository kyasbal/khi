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
} from '../../extension-common/extension-types/lifecycle-hook';
import { FRONTEND_ANALYTICS, KHIAnalyticsActivityType } from './types';
import { InspectionData } from 'src/app/store/inspection-data';
import { randomString } from 'src/app/utils/random';
import { sha512FromArrayBuffer } from 'src/app/utils/hash';
import { ReferenceType } from 'src/app/common/loader/interface';
import {
  KHIFileReferenceResolver,
  ReferenceResolverStore,
} from 'src/app/common/loader/reference-resolver';

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
    inspectionData: InspectionData,
    textBufferSource: ReferenceResolverStore,
    rawData: ArrayBuffer,
  ) => {
    const analytics = inject(FRONTEND_ANALYTICS);
    const reportInAsync = async () => {
      // analytics code to inspection data metrics
      const hash = await sha512FromArrayBuffer(rawData);
      const openId = randomString();
      analytics.globalMetadata()['inspectionDataHash'] = hash;
      analytics.globalMetadata()['openId'] = openId;

      // Calculate the total text data size
      const binaryPartReader = textBufferSource.resolvers.find((r) =>
        r.isSupportedReferenceType(ReferenceType.KHIFileBinary),
      );
      let binaryPartSize = 0;
      if (binaryPartReader) {
        binaryPartSize = (
          binaryPartReader as KHIFileReferenceResolver
        ).sourceBuffers.reduce((prev, next) => next.byteLength + prev, 0);
      }

      analytics.report(KHIAnalyticsActivityType.OpenInspectionData, {
        logLength: rawData.byteLength,
        decompressedTextBufferLength: binaryPartSize,
        revisionCount: inspectionData.timelines.reduce(
          (prev, next) => next.revisions.length + prev,
          0,
        ),
        eventCount: inspectionData.timelines.reduce(
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
