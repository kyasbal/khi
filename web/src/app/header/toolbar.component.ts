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

import { Component, OnDestroy, OnInit, inject } from '@angular/core';
import { FormControl } from '@angular/forms';
import {
  BehaviorSubject,
  Subject,
  combineLatest,
  debounceTime,
  distinctUntilChanged,
  map,
  takeUntil,
} from 'rxjs';
import * as generated from '../generated';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';
import { ViewStateService } from '../services/view-state.service';
import { nonEmptyOrDefaultString } from '../utils/state-util';
import { SelectionManagerService } from '../services/selection-manager.service';
import { BreakpointObserver } from '@angular/cdk/layout';

@Component({
  selector: 'khi-header-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrls: ['./toolbar.component.sass'],
})
export class ToolbarComponent implements OnInit, OnDestroy {
  private destoroyed = new Subject<void>();

  private breakpointObserver = inject(BreakpointObserver);

  showButtonLabel = this.breakpointObserver
    .observe(['(min-width: 1200px)'])
    .pipe(map((result) => result.matches));

  LOG_FILTER_DEBOUNCE_TIME = 200;

  filterControl: FormControl = new FormControl('');

  timezoneShift$ = this.viewStateService.timezoneShift;

  kinds$ = this.inspectionDataStore.$resourceKinds;
  includedKinds$ = this.inspectionDataStore.$kindTimelineFilter;
  namespaces$ = this.inspectionDataStore.$resourceNamespaces;
  subresourceRelationships =
    this.inspectionDataStore.subresourceRelationships.pipe(
      map((rels) => {
        const relationshipLabels = new Set<string>();
        for (const relationship of rels) {
          relationshipLabels.add(
            generated.ParentRelationshipToLabel(relationship),
          );
        }
        return relationshipLabels;
      }),
    );
  includedSubresourceRelationships =
    this.inspectionDataStore.subresourceParentRelationshipFilter.pipe(
      map((rels) => {
        const relationshipLabels = new Set<string>();
        for (const relationship of rels) {
          relationshipLabels.add(
            generated.ParentRelationshipToLabel(relationship),
          );
        }
        return relationshipLabels;
      }),
    );
  includedNamespaces$ = this.inspectionDataStore.$namespaceTimelineFilter;
  logTypes = new Set(generated.logTypes);

  logTypeFilterOpen = false;

  logFilter$ = new BehaviorSubject<string>('');

  logOrTimelineNotSelected = combineLatest([
    this.selectionManager.selectedLog,
    this.selectionManager.selectedTimeline,
  ]).pipe(map(([l, t]) => l == null || t == null));

  kindFilterOpen = false;
  namespaceFilterOpen = false;
  subresourceFilterOpen = false;

  constructor(
    private selectionManager: SelectionManagerService,
    private viewStateService: ViewStateService,
    private inspectionDataStore: InspectionDataStoreService,
  ) {}
  ngOnDestroy(): void {
    this.destoroyed.next();
  }

  ngOnInit(): void {
    const $reductedFilterEvent = this.logFilter$.pipe(
      takeUntil(this.destoroyed),
      map((a) => nonEmptyOrDefaultString(a, '.*')),
      debounceTime(this.LOG_FILTER_DEBOUNCE_TIME),
      distinctUntilChanged(),
    );
    $reductedFilterEvent.subscribe((filter) => {
      this.inspectionDataStore.setLogRegexFilter(filter);
    });
  }

  onTimezoneshiftCommit(event: Event) {
    const value = +(event.target as HTMLInputElement).value;
    this.viewStateService.setTimezoneShift(value);
  }

  onOpenKindFilter() {
    this.kindFilterOpen = true;
    this.namespaceFilterOpen = false;
    this.logTypeFilterOpen = false;
    this.subresourceFilterOpen = false;
  }

  onOpenNamespaceFilter() {
    this.kindFilterOpen = false;
    this.namespaceFilterOpen = true;
    this.logTypeFilterOpen = false;
    this.subresourceFilterOpen = false;
  }

  onOpenSubresourceRelationshipFilter() {
    this.kindFilterOpen = false;
    this.namespaceFilterOpen = false;
    this.logTypeFilterOpen = false;
    this.subresourceFilterOpen = true;
  }

  onCloseFilters() {
    this.kindFilterOpen = false;
    this.namespaceFilterOpen = false;
    this.logTypeFilterOpen = false;
    this.subresourceFilterOpen = false;
  }

  onKindFilterCommit(kinds: Set<string>) {
    this.inspectionDataStore.setKindFilter(kinds);
  }

  onNamespaceFilterCommit(namespaces: Set<string>) {
    this.inspectionDataStore.setNamespaceFilter(namespaces);
  }

  onSubresourceRelationshipFilterCommit(
    subresourceRelationshipLabels: Set<string>,
  ) {
    const relationships = [];
    for (const relationshipLabel of subresourceRelationshipLabels) {
      relationships.push(
        generated.ParseParentRelationshipLabel(relationshipLabel),
      );
    }
    this.inspectionDataStore.setRelationshipFilter(new Set(relationships));
  }

  onNameFilterChange(filter: string) {
    this.inspectionDataStore.setResourceNameRegexes(filter);
  }

  onLogFilterChange(filter: string) {
    this.logFilter$.next(filter);
  }

  openLogTypeFilter() {
    this.logTypeFilterOpen = true;
    this.kindFilterOpen = false;
    this.namespaceFilterOpen = false;
  }

  onDrawDiagram() {
    window.open(window.location.pathname + '/graph', '_blank');
  }
}
