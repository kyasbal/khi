import {
  CdkVirtualScrollViewport,
  FixedSizeVirtualScrollStrategy,
  VIRTUAL_SCROLL_STRATEGY,
} from '@angular/cdk/scrolling';
import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import {
  BehaviorSubject,
  Observable,
  Subject,
  combineLatest,
  delay,
  filter,
  map,
  shareReplay,
  takeUntil,
  withLatestFrom,
} from 'rxjs';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';
import { SelectionManagerService } from '../services/selection-manager.service';
import { ObservableCSSClassBinder } from '../utils/observable-css-class-binder';
import { LogEntry } from '../store/log';
import { TimelineEntry } from '../store/timeline';

class LogViewScrollingStrategy extends FixedSizeVirtualScrollStrategy {
  constructor() {
    super(23.33, 500, 1000);
  }
}

interface LogViewSelectionMoveCommand {
  direction: 'next' | 'prev';
}

@Component({
  selector: 'khi-log-view',
  templateUrl: './log-view.component.html',
  styleUrls: ['./log-view.component.sass'],
  providers: [
    { provide: VIRTUAL_SCROLL_STRATEGY, useClass: LogViewScrollingStrategy },
  ],
})
export class LogViewComponent implements OnInit, OnDestroy {
  private destoroyed = new Subject<void>();

  @ViewChild(CdkVirtualScrollViewport)
  viewPort!: CdkVirtualScrollViewport;

  filterByTimeline = new BehaviorSubject(true);

  logViewSelectionMoveCommand = new Subject<LogViewSelectionMoveCommand>();

  includeTimelineChildren =
    this.selectionManager.timelineSelectionShouldIncludeChildren;

  selectedLog = this.selectionManager.selectedLog;
  shownLogs: Observable<LogEntry[]> = combineLatest([
    this.inspectionDataStore.$filteredLogs,
    this.filterByTimeline,
    this.selectionManager.selectedTimelinesWithChildren,
  ]).pipe(
    map(([filteredLogs, filterByTimeline, timelines]) => {
      if (!filterByTimeline || timelines.length === 0) return filteredLogs; // show all of the filtered logs when `filterByTimeline` option is disabled.
      return this.filterLogsWithTimelines(filteredLogs, timelines);
    }),
    shareReplay(1),
  );
  allLogsCount = this.inspectionDataStore.$allLogs.pipe(
    map((logs) => logs.length),
  );
  shownLogsCount = this.shownLogs.pipe(map((logs) => logs.length));

  logViewHeight: BehaviorSubject<number> = new BehaviorSubject(
    Math.max(document.body.clientHeight - 400, 300),
  );

  highlightLogBinder = new ObservableCSSClassBinder(
    'highlight',
    'log-view-log-index',
    this.selectionManager.highlightLogIndices,
    new Set(),
  );
  selectedLogBinder = new ObservableCSSClassBinder(
    'selected',
    'log-view-log-index',
    this.selectionManager.selectedLogIndex.pipe(
      map((index) => new Set(index === -1 ? [] : [index])),
    ),
    new Set(),
  );

  disableScrollForNext = false;

  constructor(
    private inspectionDataStore: InspectionDataStoreService,
    private selectionManager: SelectionManagerService,
  ) {
    this.logViewSelectionMoveCommand
      .pipe(
        takeUntil(this.destoroyed),
        withLatestFrom(this.selectionManager.selectedLogIndex, this.shownLogs),
      )
      .subscribe(([command, index, logs]) => {
        if (index < 0) return;
        const currentSelected = this.searchArrayIndexOfLog(logs, index);
        if (index < 0) return;
        const direction = command.direction === 'prev' ? -1 : 1;
        const nextSelected = Math.max(
          0,
          Math.min(logs.length - 1, currentSelected + direction),
        );
        this.selectionManager.changeSelectionByLog(logs[nextSelected]);
      });
  }
  ngOnDestroy(): void {
    this.destoroyed.next();
  }

  ngOnInit(): void {
    this.initScrollEventOnScroll();
  }

  _selectLog(logEntry: LogEntry) {
    this.disableScrollForNext = true;
    this.selectionManager.changeSelectionByLog(logEntry);
  }

  _onLogHover(logEntry: LogEntry) {
    this.selectionManager.onHighlightLog(logEntry);
  }

  _resizeStart() {
    window.addEventListener('mouseup', () => {
      window.removeEventListener('mousemove', this._resizeMove);
    });
    window.addEventListener('mousemove', this._resizeMove);
  }

  _resizeMove = (e: MouseEvent) => {
    const current = this.logViewHeight.value;
    this.logViewHeight.next(
      Math.min(
        Math.max(100, current - e.movementY),
        document.body.clientHeight - 400,
      ),
    );
    this.viewPort.checkViewportSize();
  };

  _onScroll() {
    this.selectedLogBinder.invalidate();
    this.highlightLogBinder.invalidate();
  }

  /**
   * filterLogsWithTimelines returns a list of LogEntries being related to any of given list of timelines.
   */
  private filterLogsWithTimelines(
    logs: LogEntry[],
    timelines: TimelineEntry[],
  ): LogEntry[] {
    const logIndices = new Set<number>();
    for (const timeline of timelines) {
      for (const revision of timeline.revisions) {
        logIndices.add(revision.logIndex);
      }
      for (const event of timeline.events) {
        logIndices.add(event.logIndex);
      }
    }
    const result: LogEntry[] = [];
    for (const log of logs) {
      if (logIndices.has(log.logIndex)) {
        result.push(log);
      }
    }
    return result;
  }

  private initScrollEventOnScroll() {
    combineLatest([this.shownLogs, this.selectedLog])
      .pipe(
        filter(([, selected]) => selected !== null),
        delay(1),
      ) // Wait virtual scroll box to update the list before scroll
      .subscribe(([logs, selected]) => {
        if (!this.disableScrollForNext) {
          const arrayIndex = this.searchArrayIndexOfLog(
            logs,
            selected!.logIndex,
          );
          if (arrayIndex >= 0) {
            this.viewPort.scrollToIndex(arrayIndex, 'smooth');
          }
        }
        this.disableScrollForNext = false;
      });
  }

  public setIncludeTimelineChildren(include: boolean) {
    this.includeTimelineChildren.next(include);
  }

  /**
   * Search the array index of logEntry in given log array.
   * This method assume the log array is already sorted in logIndex order.
   * This will return -1 when the specified logIndex is not found in the list of the logs.
   */
  private searchArrayIndexOfLog(logs: LogEntry[], logIndex: number): number {
    if (logs.length === 0) return -1;
    if (logs[0].logIndex === logIndex) return 0;

    // binary search
    let less = 0;
    let moreOrEqual = logs.length - 1;
    while (Math.abs(less - moreOrEqual) > 1) {
      const middle = Math.floor((less + moreOrEqual) / 2);
      if (logs[middle].logIndex < logIndex) {
        less = middle;
      } else {
        moreOrEqual = middle;
      }
    }

    return logs[moreOrEqual].logIndex == logIndex ? moreOrEqual : -1;
  }

  public onKeyDown(keyEvent: KeyboardEvent) {
    if (keyEvent.key === 'ArrowDown') {
      this.logViewSelectionMoveCommand.next({
        direction: 'next',
      });
      keyEvent.preventDefault();
    }
    if (keyEvent.key === 'ArrowUp') {
      this.logViewSelectionMoveCommand.next({
        direction: 'prev',
      });
      keyEvent.preventDefault();
    }
  }
}
