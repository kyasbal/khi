import { Injectable } from '@angular/core';
import {
  BehaviorSubject,
  Observable,
  combineLatestWith,
  filter,
  map,
  shareReplay,
  startWith,
} from 'rxjs';
import { InspectionData, TimelineRange } from '../models/inspection-data';
import { asBehaviorSubject } from '../utils/observable-util';
import { SelectOnlyDeeperOrEqual } from '../utils/timeline-collection-util';
import { FilterWorkerService } from './filter-worker.service';
import { TextBufferLoader } from './data-loader.service';
import { ParentRelationship } from '../generated';
import { TimelineEntry, TimelineLayer } from '../store/timeline';
import { LogEntry } from '../store/log';

/**
 * Manage the inspection data(Timelines, Logs,...) after receiving it from somewhere
 * Provides filter feature on this layer
 */
@Injectable({ providedIn: 'root' })
export class InspectionDataStoreService {
  /**
   * Timeline filter subjects
   */
  public $resourceNameTimelineFilterRegex: BehaviorSubject<string> =
    new BehaviorSubject('');
  public $kindTimelineFilter: BehaviorSubject<Set<string>> =
    new BehaviorSubject(new Set());
  public $namespaceTimelineFilter: BehaviorSubject<Set<string>> =
    new BehaviorSubject(new Set());

  public $filteredOutLogIndices: BehaviorSubject<Set<number>> =
    new BehaviorSubject(new Set());

  public subresourceParentRelationshipFilter: BehaviorSubject<
    Set<ParentRelationship>
  > = new BehaviorSubject(new Set());

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

  public textBufferSource = new BehaviorSubject<TextBufferLoader | null>(null);

  /**
   * Timeline related inspection sub data
   */
  public $allTimelines: BehaviorSubject<TimelineEntry[]> = asBehaviorSubject(
    this.currentValidInspectionData.pipe(
      map((t) => t.timelines),
      startWith([]),
      shareReplay(1),
    ),
    [],
  );
  public $filteredTimelines: BehaviorSubject<TimelineEntry[]> =
    asBehaviorSubject(
      this.$allTimelines.pipe(
        combineLatestWith(
          this.$kindTimelineFilter,
          this.$namespaceTimelineFilter,
          this.subresourceParentRelationshipFilter,
          this.$resourceNameTimelineFilterRegex,
        ),
        map(
          ([
            timelines,
            kindFilter,
            namespaceFilter,
            relationshipFilter,
            resourceRegex,
          ]) =>
            this._timelineFilter(
              timelines,
              kindFilter,
              namespaceFilter,
              relationshipFilter,
              resourceRegex,
            ),
        ),
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
  public $resourceKinds = this.currentValidInspectionData.pipe(
    map((t) => t.kinds),
    startWith(new Set<string>()),
  );
  public $resourceNamespaces = this.currentValidInspectionData.pipe(
    map((t) => t.namespaces),
    startWith(new Set<string>()),
  );
  public subresourceRelationships = this.currentValidInspectionData.pipe(
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

  /**
   * Log related inspection sub data
   */
  public $allLogs: BehaviorSubject<LogEntry[]> = asBehaviorSubject(
    this.currentValidInspectionData.pipe(
      map((t) => t.logs),
      startWith([]),
    ),
    [],
  );
  public $filteredLogs = new BehaviorSubject<LogEntry[]>([]);

  private filterWorker: FilterWorkerService = new FilterWorkerService(this);

  constructor() {
    // Replace filter subject fields with the default values when the available filter options are updated(When new inspection data was loaded.)
    this.$resourceKinds.subscribe((k) => this.$kindTimelineFilter.next(k));
    this.$resourceNamespaces.subscribe((n) =>
      this.$namespaceTimelineFilter.next(n),
    );
    this.subresourceRelationships.subscribe((rels) => {
      this.subresourceParentRelationshipFilter.next(rels);
    });
    this.$allLogs.subscribe((allLogs) => {
      this.$filteredLogs.next(allLogs);
    });

    this.$filteredOutLogIndices.subscribe(() => this._onLogFilterUpdate());
  }

  /**
   * Get the list of timelines related to the specified log entry.
   */
  public findTimelineFromLog(
    log: LogEntry,
    includingFilteredOut = false,
  ): TimelineEntry[] {
    const source = includingFilteredOut
      ? this.$allTimelines.value
      : this.$filteredTimelines.value;
    const sourceSet = new Set(source);
    const result: TimelineEntry[] = [];
    for (const relatedTimeline of log.relatedTimelines) {
      if (sourceSet.has(relatedTimeline)) result.push(relatedTimeline);
    }
    return result;
  }

  public setNewInspectionData(
    data: InspectionData,
    textBufferSource: TextBufferLoader,
  ) {
    this.$inspectionData.next(data);
    this.textBufferSource.next(textBufferSource);
  }

  public setResourceNameRegexes(regexes: string) {
    this.$resourceNameTimelineFilterRegex.next(regexes == '' ? '.*' : regexes);
  }

  public setKindFilter(kinds: Set<string>) {
    this.$kindTimelineFilter.next(kinds);
  }

  public setNamespaceFilter(namespaces: Set<string>) {
    this.$namespaceTimelineFilter.next(namespaces);
  }

  public setRelationshipFilter(relationships: Set<ParentRelationship>) {
    this.subresourceParentRelationshipFilter.next(relationships);
  }

  public async setLogRegexFilter(filter: string) {
    this.$filteredOutLogIndices.next(
      await this.filterWorker.filterLogs(this.$allLogs.value, filter),
    );
  }

  private _timelineFilter(
    allTimelines: TimelineEntry[],
    currentKindFilter: Set<string>,
    currentNamespaceFilter: Set<string>,
    currentSubresourceRelationshipFilter: Set<ParentRelationship>,
    regexFilter: string,
  ): TimelineEntry[] {
    const filteredTimelines: TimelineEntry[] = [];
    let nameFilterRegexs: RegExp[] = [];
    const resourceNameRegexes = regexFilter;
    if (resourceNameRegexes === '') {
      nameFilterRegexs.push(/.*/);
    } else {
      nameFilterRegexs = resourceNameRegexes
        .split(' ')
        .filter((f) => f != '')
        .map((filter) => new RegExp(filter));
    }

    let lastNamespaceLayer: TimelineEntry | null = null;
    for (const timeline of allTimelines) {
      if (!currentKindFilter.has(timeline.getNameOfLayer(TimelineLayer.Kind)))
        continue;
      if (
        timeline.layer >= TimelineLayer.Namespace &&
        !currentNamespaceFilter.has(
          timeline.getNameOfLayer(TimelineLayer.Namespace),
        )
      )
        continue;
      if (
        timeline.layer >= TimelineLayer.Name &&
        !this._processRegexNameFilters(timeline, nameFilterRegexs)
      )
        continue;
      if (timeline.layer < TimelineLayer.Namespace) {
        filteredTimelines.push(timeline);
      } else if (timeline.layer === TimelineLayer.Namespace) {
        lastNamespaceLayer = timeline;
      } else if (timeline.layer === TimelineLayer.Name) {
        if (lastNamespaceLayer) {
          filteredTimelines.push(lastNamespaceLayer);
          lastNamespaceLayer = null;
        }
        filteredTimelines.push(timeline);
      } else {
        if (
          currentSubresourceRelationshipFilter.has(timeline.parentRelationship)
        ) {
          filteredTimelines.push(timeline);
        }
      }
    }
    // Remove unused kind or namespace
    const result = SelectOnlyDeeperOrEqual(
      filteredTimelines,
      TimelineLayer.Name,
    );
    return result;
  }

  private _logFilter(allLogs: LogEntry[]): LogEntry[] {
    const logIndices = this.$filteredOutLogIndices.value;
    return allLogs.filter((l) => !logIndices.has(l.logIndex));
  }

  private _onLogFilterUpdate() {
    this.$filteredLogs.next(this._logFilter(this.$allLogs.value));
  }

  private _processRegexNameFilters(
    elem: TimelineEntry,
    filterRegexs: RegExp[],
  ): boolean {
    for (const regex of filterRegexs) {
      if (regex.test(elem.getNameOfLayer(TimelineLayer.Name))) return true;
    }
    return false;
  }
}
