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

import { Injectable } from '@angular/core';
import {
  BehaviorSubject,
  Observable,
  ReplaySubject,
  combineLatest,
  debounceTime,
  filter,
  map,
  shareReplay,
  startWith,
  switchMap,
} from 'rxjs';
import { InspectionData, TimelineRange } from '../models/inspection-data';
import { asBehaviorSubject } from '../utils/observable-util';
import { FilterWorkerService } from './filter-worker.service';
import { ParentRelationship } from '../generated';
import { TimelineEntry } from '../store/timeline';
import { LogEntry } from '../store/log';
import { ReferenceResolverStore } from '../common/loader/reference-resolver';

/**
 * InspectionDataStore provides observable to the inspection data loaded.
 * This store won't provide any filterings performed in response to user's interactions.
 * Implementation of this class must compute values and emit them only when another inspection data was loaded.
 */
export interface InspectionDataStore {
  /**
   * allTimelines emits the array of timeline entry without any filter.
   */
  allTimelines: Observable<TimelineEntry[]>;

  /**
   * availableKinds emits the set of all kind names found in the inspection data.
   */
  availableKinds: Observable<Set<string>>;

  /**
   * availableNamespaces emits the set of all namespaces found in the inspection data.
   */
  availableNamespaces: Observable<Set<string>>;

  /**
   * availableSubresourceParentRelationships emits the set of all parent relationships of subresources in the inspection data.
   */
  availableSubresourceParentRelationships: Observable<Set<ParentRelationship>>;
}

/**
 * Manage the inspection data(Timelines, Logs,...) after receiving it from somewhere
 * Provides filter feature on this layer
 */
@Injectable({ providedIn: 'root' })
export class InspectionDataStoreService implements InspectionDataStore {
  /**
   * Source of the inspection data
   */
  public $inspectionData: BehaviorSubject<InspectionData | null> =
    new BehaviorSubject<InspectionData | null>(null);

  /**
   * Inspectiondata with null check.
   */
  public currentValidInspectionData = this.$inspectionData.pipe(
    filter((d) => !!d),
  ) as Observable<InspectionData>;

  public textBufferSource = new BehaviorSubject<ReferenceResolverStore | null>(
    null,
  );

  /**
   * Timeline related inspection sub data
   */
  public allTimelines: BehaviorSubject<TimelineEntry[]> = asBehaviorSubject(
    this.currentValidInspectionData.pipe(
      map((t) => t.timelines),
      startWith([]),
      shareReplay(1),
    ),
    [],
  );

  public $timeRange: BehaviorSubject<TimelineRange> = asBehaviorSubject(
    this.currentValidInspectionData.pipe(
      map((t) => t.range),
      startWith(new TimelineRange(0, 0)),
    ),
    new TimelineRange(0, 0),
  );
  public availableKinds = this.currentValidInspectionData.pipe(
    map((t) => t.kinds),
    startWith(new Set<string>()),
  );
  public availableNamespaces = this.currentValidInspectionData.pipe(
    map((t) => t.namespaces),
    startWith(new Set<string>()),
  );
  public availableSubresourceParentRelationships =
    this.currentValidInspectionData.pipe(
      filter((data) => data !== null),
      map((data) => data.relationships),
      shareReplay(1),
      startWith(new Set<ParentRelationship>()),
      map((relationshipSet) => {
        // Unknown type is not used in subresources. Ignore this type to be included in the applicable filter.
        relationshipSet.delete(ParentRelationship.Unknown);
        return relationshipSet;
      }),
    );

  public allLogs: BehaviorSubject<LogEntry[]> = asBehaviorSubject(
    this.currentValidInspectionData.pipe(
      map((t) => t.logs),
      startWith([]),
    ),
    [],
  );
  private logFilter = new ReplaySubject<string>(1);

  /**
   * An observable emits the Set of log indices to be filtered out.
   */
  public filteredOutLogIndicesSet = combineLatest([
    this.allLogs,
    this.logFilter.pipe(startWith('')),
  ]).pipe(
    debounceTime(0),
    switchMap(([allLogs, filter]) =>
      this.filterWorker.filterLogs(allLogs, filter),
    ),
    shareReplay(1),
    startWith(new Set<number>()),
  );

  /**
   * An observable emits list of logs filtered.
   */
  public filteredLogs = combineLatest([
    this.allLogs,
    this.filteredOutLogIndicesSet,
  ]).pipe(
    debounceTime(0),
    map(([allLogs, filteredOutLogs]) =>
      allLogs.filter((_, index) => !filteredOutLogs.has(index)),
    ),
    shareReplay(1),
  );

  private filterWorker: FilterWorkerService = new FilterWorkerService(this);

  public setNewInspectionData(
    data: InspectionData,
    textBufferSource: ReferenceResolverStore,
  ) {
    this.$inspectionData.next(data);
    this.textBufferSource.next(textBufferSource);
  }

  public async setLogRegexFilter(filter: string) {
    this.logFilter.next(filter);
  }
}
