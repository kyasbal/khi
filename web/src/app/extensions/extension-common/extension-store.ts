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

import { TimelineEntry } from 'src/app/store/timeline';
import {
  DisplayableTimelineNavigatorExtension,
  TimelineNavigatorExtension,
} from './extension-types/timeline-navigator';

/**
 * ExtensionStore is the type to hold the reference to plugin instances of each extensible points in KHI.
 */
export interface ExtensionStore {
  timelineNavigatorExtensions: TimelineNavigatorExtension[];
}

/**
 * GlobalExtensionStore is the singleton instance of ExtensionStore and it's the default ExtensionStore in KHI frontend.
 */
export const GlobalExtensionStore: ExtensionStore = {
  timelineNavigatorExtensions: [],
};

/**
 * ExtensionStoreUtil provide utilitiy function to interact the actual extension references.
 */
export class ExtensionStoreUtil {
  /**
   * Returns the visible extensions for the given timeline.
   */
  public static getVisibleTimelineNavigatorExtensions(
    store: ExtensionStore,
    timeline: TimelineEntry,
  ): DisplayableTimelineNavigatorExtension[] {
    return store.timelineNavigatorExtensions
      .filter((extension) => extension.show(timeline))
      .map((extension) => extension.getDisplayable(timeline));
  }
}
