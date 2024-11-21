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
import { InspectionDataStore } from './inspection-data-store.service';
import {
  combineLatest,
  connectable,
  debounceTime,
  map,
  merge,
  ReplaySubject,
  shareReplay,
  Subject,
} from 'rxjs';
import { ParentRelationship } from '../generated';
import { TimelineEntry, TimelineLayer } from '../store/timeline';

/**
 * Injection token for the default TimelineFilter.
 */
export const DEFAULT_TIMELINE_FILTER = new InjectionToken(
  'DEFAULT_TIMELINE_FILTER',
);

/**
 * TimelineFilter provides Observable fields for timelines.
 * It listen changes on inspection data store and filters timelines with given conditions.
 */
export class TimelineFilter {
  constructor(public readonly dataStore: InspectionDataStore) {
    this.kindTimelineFilter.connect();
    this.namespaceTimelineFilter.connect();
    this.subresourceParentRelationshipFilter.connect();
    this.resourceNameTimelineRegexFilter.connect();
  }

  /**
   * Observable for currently selected kind name set.
   */
  private readonly kindTimelineFilterSubject = new Subject<Set<string>>();

  /**
   * Observable for currently selected kind name set.
   * This observable also emits the all kind names on the list of available kind names are changed.
   */
  public readonly kindTimelineFilter = connectable(
    merge(this.kindTimelineFilterSubject, this.dataStore.availableKinds),
    {
      connector: () => new ReplaySubject(1),
      resetOnDisconnect: false,
    },
  );

  /**
   * Observable for currently selected namespace set.
   */
  private readonly namespaceTimelineFilterSubject = new Subject<Set<string>>();

  /**
   * Observable for currently selected namespace set.
   * This observable also emits the all namespaces on the list of available namespaces are changed.
   */
  public readonly namespaceTimelineFilter = connectable(
    merge(
      this.namespaceTimelineFilterSubject,
      this.dataStore.availableNamespaces,
    ),
    {
      connector: () => new ReplaySubject(1),
      resetOnDisconnect: false,
    },
  );

  /**
   * Observable for currently selected parent relationship set for subresources.
   */
  private readonly subresourceParentRelationshipFilterSubject = new Subject<
    Set<ParentRelationship>
  >();

  /**
   * Observable for currently selected parent relationships of subresources.
   * This observable also emits the all parent relationships of subresource when the available parent relationships are changed.
   */
  public readonly subresourceParentRelationshipFilter = connectable(
    merge(
      this.subresourceParentRelationshipFilterSubject,
      this.dataStore.availableSubresourceParentRelationships,
    ),
    {
      connector: () => new ReplaySubject(1),
      resetOnDisconnect: false,
    },
  );

  /**
   * Observable for current active regex name filter for resource names.
   */
  private readonly resourceNameTimelineRegexFilterSubject =
    new Subject<string>();

  /**
   * Observable for currently used regex filter for resource name.
   * This observable resets the filter when a new data loaded on the data store.
   */
  public readonly resourceNameTimelineRegexFilter = connectable(
    merge(
      this.resourceNameTimelineRegexFilterSubject,
      this.dataStore.allTimelines.pipe(map(() => '')),
    ),
    {
      connector: () => new ReplaySubject(1),
      resetOnDisconnect: false,
    },
  );

  /**
   * Observable emitting RegExp list parsed from resourceNameTimelineRegexFilter. This allows user to use multiple filter with splitting regex filter with white space.
   */
  private readonly resourceNameTimelineRegexFilterInRegExpList =
    this.resourceNameTimelineRegexFilter.pipe(
      map((regex) =>
        regex
          .split(' ')
          .filter((f) => f != '')
          .map((filter) => new RegExp(filter)),
      ), // convert to the list of RegExp[]
      map((regexps) => (regexps.length == 0 ? [/.*/] : regexps)), // return a regex matching anything when no filter provided.
    );

  /**
   * The list of timelines filtered by this TimelineFilter.
   */
  public readonly filteredTimeline = combineLatest([
    this.dataStore.allTimelines,
    this.kindTimelineFilter,
    this.namespaceTimelineFilter,
    this.subresourceParentRelationshipFilter,
    this.resourceNameTimelineRegexFilterInRegExpList,
  ]).pipe(
    debounceTime(0),
    map(
      ([
        timelines,
        kindSet,
        namespaceSet,
        subresourceParentRelationshipSet,
        regexList,
      ]) =>
        this.filterTimelines(
          timelines,
          kindSet,
          namespaceSet,
          subresourceParentRelationshipSet,
          regexList,
        ),
    ),
    shareReplay(1),
  );

  /**
   * setKindFilter limits the filteredTimeline by kind names.
   */
  public setKindFilter(kinds: Set<string>) {
    this.kindTimelineFilterSubject.next(kinds);
  }

  /**
   * setNamespaceFilter limits the filteredTimeline by namespace names.
   */
  public setNamespaceFilter(namespaces: Set<string>) {
    this.namespaceTimelineFilterSubject.next(namespaces);
  }

  /**
   * setSubresourceParentRelationshipFilter limits the filteredTimeline of subresources by its relationship.
   */
  public setSubresourceParentRelationshipFilter(
    relationships: Set<ParentRelationship>,
  ) {
    this.subresourceParentRelationshipFilterSubject.next(relationships);
  }

  /**
   * setResourceNameRegexFilter limits the filteredTimeline of resources by its name.
   */
  public setResourceNameRegexFilter(regexFilter: string) {
    this.resourceNameTimelineRegexFilterSubject.next(regexFilter);
  }

  private filterTimelines(
    timelines: TimelineEntry[],
    kindFilter: Set<string>,
    namespaceFilter: Set<string>,
    parentRelationshipFilter: Set<ParentRelationship>,
    regexFilter: RegExp[],
  ): TimelineEntry[] {
    const filteredTimelines: TimelineEntry[] = [];
    let lastProcessedKind: TimelineEntry | null = null;
    let lastProcessedNamespace: TimelineEntry | null = null;
    for (const timeline of timelines) {
      if (!kindFilter.has(timeline.getNameOfLayer(TimelineLayer.Kind))) {
        // timeline is filtered by kind name.
        continue;
      }
      if (
        timeline.layer >= TimelineLayer.Namespace &&
        !namespaceFilter.has(timeline.getNameOfLayer(TimelineLayer.Namespace))
      ) {
        // timeline is filtered by namespace name and the timeline is not a kind.
        continue;
      }
      const resourceName = timeline.getNameOfLayer(TimelineLayer.Name);
      if (
        timeline.layer >= TimelineLayer.Name &&
        !regexFilter.some((regex) => regex.test(resourceName))
      ) {
        // timeline is filtered by name filter and the timeline is not a kind or a namespace.
        continue;
      }
      if (
        timeline.layer >= TimelineLayer.Subresource &&
        !parentRelationshipFilter.has(timeline.parentRelationship)
      ) {
        // timeline is filtered by subresource parent relationship filter, and the timeline is not a kind,a namespace or a resource.
        continue;
      }

      if (timeline.layer === TimelineLayer.Kind) {
        // if the timeline is at kind layer, the timeline will be ignored when there is no child timelines are included in the filter result.
        // Deferring including kind timeline until the next namespace layer timeline being included.
        lastProcessedKind = timeline;
        continue;
      } else if (timeline.layer === TimelineLayer.Namespace) {
        // if the timeline is at namespace layer, the timeline will be ignored when there is no child timelines are included in the filter result.
        // Deferring including namespace timeline until the next resource layer timeline being included.
        lastProcessedNamespace = timeline;
        continue;
      } else if (timeline.layer === TimelineLayer.Name) {
        if (lastProcessedNamespace) {
          if (lastProcessedKind) {
            filteredTimelines.push(lastProcessedKind);
            lastProcessedKind = null;
          }
          filteredTimelines.push(lastProcessedNamespace);
          lastProcessedNamespace = null;
        }
      }
      filteredTimelines.push(timeline);
    }
    return filteredTimelines;
  }
}
