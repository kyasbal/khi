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

import { InjectionToken } from '@angular/core';
import { PageType } from '../../extension-common/extension-types/lifecycle-hook';

export const FRONTEND_ANALYTICS = new InjectionToken<FrontendAnalytics>(
  'FRONTEND_ANAYTICS',
);

/**
 * Type of activity recorded in Analytics
 */
export enum KHIAnalyticsActivityType {
  Init = 'INIT',
  Inspect = 'INSPECT',
  OpenInspectionData = 'OPEN_INSPECTION_DATA',
}

export type KHIAnalyticsActivityMetadata = {
  [metadataKey: string]: number | string;
};

/**
 * FrontendAnalytics is an interface to report events.
 */
export interface FrontendAnalytics {
  /**
   * initialize analytics provider. This must be called when a page loaded.
   */
  init(pageType: PageType): void;

  /**
   * Report the event.
   * @param event type of event
   * @param metadata supplimental information attached to this event.
   */
  report(
    event: KHIAnalyticsActivityType,
    metadata: KHIAnalyticsActivityMetadata,
  ): void;

  /**
   * Returns the reference to global metadata in addition to the metadata argument reported with `report` method.
   */
  globalMetadata(): KHIAnalyticsActivityMetadata;
}
