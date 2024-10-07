import { ViewStateService } from '../services/view-state.service';
import { TimelinenCoordinateCalculator } from './timeline-coordinate-calculator';
import { asBehaviorSubject } from '../utils/observable-util';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';
import { Observable, combineLatestWith, debounceTime, delay } from 'rxjs';
import { CanvasSize } from './canvas/types';
import * as generated from '../generated';
import { SelectionManagerService } from '../services/selection-manager.service';
import { LogEntry } from '../store/log';
import { TimelineEntry } from '../store/timeline';

const STRIDE_SIZE_CANDIDATES = [
  1000,
  5000,
  10000,
  30000,
  60000,
  300000,
  600000,
  1800000,
  3600000,
  3 * 60 * 60 * 1000,
  6 * 60 * 60 * 1000,
  12 * 60 * 60 * 1000,
  24 * 60 * 60 * 1000,
  3 * 24 * 60 * 60 * 1000,
  7 * 24 * 60 * 60 * 1000,
];

interface SeverityLogGroupingResult {
  maxLogCountInSingleBucket: number;
  logsBySeverity: { [severity: string]: number }[];
}

export class BackgroundCanvas {
  private canvas: HTMLCanvasElement;

  private ctx: CanvasRenderingContext2D;

  private _invalidate = true;

  private _timezoneShift$ = asBehaviorSubject(
    this._viewStateService.timezoneShift,
    0,
  );

  public MINIMUM_DISTANT = 5;

  public MAX_RULER_THICKNESS = 5;

  selectedLogTimestamp = 0;

  currentAllFilteredLogs: LogEntry[] = [];

  currentHighlightedLogs: LogEntry[] = [];

  constructor(
    canvas: HTMLCanvasElement,
    private canvasSizeObservable: Observable<CanvasSize>,
    private timelineCoordinateCalculator: TimelinenCoordinateCalculator,
    private _inspectionData: InspectionDataStoreService,
    private _viewStateService: ViewStateService,
    private _selectionManagerService: SelectionManagerService,
  ) {
    this.canvas = canvas;
    const context = this.canvas.getContext('2d');
    if (context != null) {
      this.ctx = context;
    } else {
      throw new Error("Couldn't get canvas context");
    }

    this._inspectionData.$filteredLogs
      .pipe(
        combineLatestWith(
          this._inspectionData.$filteredTimelines,
          this._selectionManagerService.selectedTimelinesWithChildren,
        ),
      )
      .subscribe(([logs, filteredTimelines, selectedTimelines]) => {
        this.currentAllFilteredLogs = this.filterLogsWithTimelines(
          logs,
          filteredTimelines,
        );
        this.currentHighlightedLogs = this.filterLogsWithTimelines(
          logs,
          selectedTimelines,
        );
        this.invalidate();
      });
    this.canvasSizeObservable.pipe(debounceTime(5)).subscribe((size) => {
      this._fitCanvasSize(size.width, size.height);
      this.invalidate();
    });
    // To redraw canvas background in response to scale up/down the page with the browser feature
    this._viewStateService.devicePixelRatio.pipe(delay(1)).subscribe(() => {
      this.invalidate();
    });
    this.invalidate();
    this._timezoneShift$.subscribe(() => this.invalidate());
  }

  private _render() {
    if (this._invalidate) {
      this._invalidate = false;

      this._clearBackground();
      this.drawRulers();
      this.drawDateLabels();
      this._drawHistgram();
      this._drawSelectedLogLine();
      this._drawDataRegion();
    }
  }

  public invalidate() {
    this._invalidate = true;
    this._render();
  }

  private _fitCanvasSize(width: number, height: number) {
    this.canvas.width =
      width * this.timelineCoordinateCalculator.canvasPixelRatio;
    this.canvas.height =
      height * this.timelineCoordinateCalculator.canvasPixelRatio;
    this.canvas.style.scale =
      '' + 1 / this.timelineCoordinateCalculator.canvasPixelRatio;
    this.canvas.style.transformOrigin = '0% 0%';
    this._invalidate = true;
  }

  private _clearBackground() {
    this.ctx.fillStyle = 'white';
    this.ctx.strokeStyle = 'transparent';
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
  }

  private drawDateLabels() {
    const timeColor = '#A0A0A0';
    const dayLineColor = '#CCAAFF';
    const dayDurationInMs = 1000 * 60 * 60 * 24;
    const leftSafePaddingInPx = 430; // width of the scaling tips element.
    const time =
      this.timelineCoordinateCalculator.leftOffsetToTime(0) - dayDurationInMs;
    const date = this.timeToDate(time);
    const leftMostDate = this.toNearestMidnightBefore(date.getTime());
    const headerSize = this.timelineCoordinateCalculator.adjustPixelScale(60);

    const fontSize = this.timelineCoordinateCalculator.adjustPixelScale(10);
    this.ctx.font = `${fontSize}px serif`;
    let currentDateStart = leftMostDate.getTime();
    this.ctx.strokeStyle = timeColor;
    this.ctx.fillStyle = timeColor;
    let currentLabel = this.toDateLabelString(currentDateStart);
    const labelSize = this.ctx.measureText(currentLabel);
    let currentDateOffset =
      this.timelineCoordinateCalculator.timeToLeftOffset(currentDateStart);
    while (currentDateOffset < 0) {
      const nextDayOffset = this.timelineCoordinateCalculator.timeToLeftOffset(
        currentDateStart + dayDurationInMs,
      );
      if (nextDayOffset < 0) {
        currentDateStart += dayDurationInMs;
        currentLabel = this.toDateLabelString(currentDateStart);
        currentDateOffset =
          this.timelineCoordinateCalculator.timeToLeftOffset(currentDateStart);
        continue;
      }
      if (
        nextDayOffset - currentDateOffset > labelSize.width + 30 &&
        nextDayOffset > leftSafePaddingInPx + labelSize.width + 30
      ) {
        this.ctx.fillText(
          currentLabel,
          leftSafePaddingInPx - labelSize.width / 2,
          25,
        );
        currentDateStart += dayDurationInMs;
      }
      break;
    }
    while (
      currentDateStart <
      this.timelineCoordinateCalculator.rightMostTime() + dayDurationInMs
    ) {
      const currentLabel = this.toDateLabelString(currentDateStart);
      this.ctx.fillText(
        currentLabel,
        this.timelineCoordinateCalculator.timeToLeftOffset(currentDateStart),
        25,
      );

      const extraHeight =
        this.timelineCoordinateCalculator.adjustPixelScale(10);
      this.ctx.strokeStyle = dayLineColor;
      this.ctx.lineWidth = this.MAX_RULER_THICKNESS;
      this.ctx.beginPath();
      this.ctx.moveTo(
        this.timelineCoordinateCalculator.timeToLeftOffset(currentDateStart),
        headerSize - extraHeight,
      );
      this.ctx.lineTo(
        this.timelineCoordinateCalculator.timeToLeftOffset(currentDateStart),
        this.canvas.height,
      );
      this.ctx.stroke();
      currentDateStart += dayDurationInMs;
    }
  }

  private drawRulers() {
    const lineColor = '#E0E0E0';
    const timeColor = '#A0A0A0';
    this.ctx.strokeStyle = lineColor;
    let thickness = this.MAX_RULER_THICKNESS;
    const largestRulerStride = this._calcStrideSizeAt(0);
    const headerSize = this.timelineCoordinateCalculator.adjustPixelScale(60);
    let extraHeight = this.timelineCoordinateCalculator.adjustPixelScale(40);
    for (
      let rulerTimeStride = largestRulerStride;
      this.timelineCoordinateCalculator.durationToWidth(rulerTimeStride) >=
      this.MINIMUM_DISTANT;
      rulerTimeStride /= 5
    ) {
      this.ctx.lineWidth = thickness;
      this.ctx.beginPath();
      const isPrimaryLine = rulerTimeStride === largestRulerStride;
      let rulerTime =
        Math.ceil(
          this.timelineCoordinateCalculator.leftMostTime() / rulerTimeStride,
        ) * rulerTimeStride;
      while (
        this.timelineCoordinateCalculator.timeToLeftOffset(rulerTime) <
        this.canvas.width
      ) {
        if (isPrimaryLine) {
          this.ctx.strokeStyle = timeColor;
          this.ctx.fillStyle = timeColor;
          const fontSize =
            this.timelineCoordinateCalculator.adjustPixelScale(14);
          this.ctx.font = `${fontSize}px serif`;
          this.ctx.lineWidth = 1;
          const label = this.toTimeLabelString(rulerTime);
          const measured = this.ctx.measureText(label);
          this.ctx.fillText(
            label,
            this.timelineCoordinateCalculator.timeToLeftOffset(rulerTime) -
              measured.width / 2,
            headerSize - this.timelineCoordinateCalculator.adjustPixelScale(35),
          );
          this.ctx.lineWidth = thickness;
          this.ctx.stroke();
          this.ctx.strokeStyle = lineColor;
        }
        this.ctx.moveTo(
          this.timelineCoordinateCalculator.timeToLeftOffset(rulerTime),
          headerSize - extraHeight,
        );
        this.ctx.lineTo(
          this.timelineCoordinateCalculator.timeToLeftOffset(rulerTime),
          this.canvas.height,
        );
        rulerTime += rulerTimeStride;
      }
      this.ctx.stroke();
      thickness /= 2;
      extraHeight /= 2;
    }
  }

  private _drawDataRegion() {
    const range = this._inspectionData.$timeRange.value;
    const minimumFromLeftTime =
      range.begin - this.timelineCoordinateCalculator.leftMostTime();
    if (minimumFromLeftTime > 0) {
      this.ctx.fillStyle = 'rgba(0,0,0,0.1)';
      this.ctx.fillRect(
        0,
        0,
        this.timelineCoordinateCalculator.durationToWidth(minimumFromLeftTime),
        this.canvas.height,
      );
    }
    const maximumFromLeftTime =
      range.end - this.timelineCoordinateCalculator.leftMostTime();
    if (maximumFromLeftTime > 0) {
      this.ctx.fillStyle = 'rgba(0,0,0,0.1)';
      this.ctx.fillRect(
        this.timelineCoordinateCalculator.durationToWidth(maximumFromLeftTime),
        0,
        this.canvas.width,
        this.canvas.height,
      );
    }
  }

  private _drawHistgram() {
    if (this.currentAllFilteredLogs.length === 0) {
      return;
    }
    const range = this._inspectionData.$timeRange.value;
    const histogramMaxSizeInPx = 50;
    const secondStride = this._calcStrideSizeAt(0) / 50;
    const boxWidth =
      this.timelineCoordinateCalculator.durationToWidth(secondStride);
    const histgramLeftTime = Math.max(
      Math.ceil(
        this.timelineCoordinateCalculator.leftMostTime() / secondStride,
      ) *
        secondStride -
        secondStride,
      range.begin,
    );
    this.ctx.lineWidth = 1;

    // Count logs by severity with spliting timerange by secondStride
    const logCountsForAll = this.countLogsBySeverity(
      this.currentAllFilteredLogs,
      histgramLeftTime,
      Math.min(
        this.canvas.width,
        this.timelineCoordinateCalculator.timeToLeftOffset(range.end),
      ),
      secondStride,
    );
    const logCountsForHighlights = this.countLogsBySeverity(
      this.currentHighlightedLogs,
      histgramLeftTime,
      Math.min(
        this.canvas.width,
        this.timelineCoordinateCalculator.timeToLeftOffset(range.end),
      ),
      secondStride,
    );
    let currentBucketLeftTime = histgramLeftTime;
    let bucketIndex = 0;
    while (
      this.timelineCoordinateCalculator.timeToLeftOffset(
        currentBucketLeftTime,
      ) <
      Math.min(
        this.canvas.width,
        this.timelineCoordinateCalculator.timeToLeftOffset(range.end),
      )
    ) {
      let bottomHeight = 0;
      for (const severity of generated.severities) {
        // 33 is alpha part of RGB color code example: #0000FF33
        this.ctx.strokeStyle = generated.severityBorderColors[severity] + '33';
        this.ctx.fillStyle = generated.severityColors[severity] + '33';
        const histRatio =
          logCountsForAll.logsBySeverity[bucketIndex][severity] /
          logCountsForAll.maxLogCountInSingleBucket;
        const histHeight = histogramMaxSizeInPx * histRatio;
        this.ctx.fillRect(
          this.timelineCoordinateCalculator.timeToLeftOffset(
            currentBucketLeftTime,
          ),
          this.timelineCoordinateCalculator.adjustPixelScale(
            60 - histHeight - bottomHeight,
          ),
          boxWidth,
          this.timelineCoordinateCalculator.adjustPixelScale(histHeight),
        );
        this.ctx.strokeRect(
          this.timelineCoordinateCalculator.timeToLeftOffset(
            currentBucketLeftTime,
          ),
          this.timelineCoordinateCalculator.adjustPixelScale(
            60 - histHeight - bottomHeight,
          ),
          boxWidth,
          this.timelineCoordinateCalculator.adjustPixelScale(histHeight),
        );
        bottomHeight += histHeight;
      }
      bottomHeight = 0;
      for (const severity of generated.severities) {
        // CC is alpha part of RGB color code example: #0000FFCC
        this.ctx.strokeStyle = generated.severityBorderColors[severity] + 'CC';
        this.ctx.fillStyle = generated.severityColors[severity] + 'CC';
        const histRatio =
          logCountsForHighlights.logsBySeverity[bucketIndex][severity] /
          logCountsForAll.maxLogCountInSingleBucket;
        const histHeight = histogramMaxSizeInPx * histRatio;
        this.ctx.fillRect(
          this.timelineCoordinateCalculator.timeToLeftOffset(
            currentBucketLeftTime,
          ),
          this.timelineCoordinateCalculator.adjustPixelScale(
            60 - histHeight - bottomHeight,
          ),
          boxWidth,
          this.timelineCoordinateCalculator.adjustPixelScale(histHeight),
        );
        this.ctx.strokeRect(
          this.timelineCoordinateCalculator.timeToLeftOffset(
            currentBucketLeftTime,
          ),
          this.timelineCoordinateCalculator.adjustPixelScale(
            60 - histHeight - bottomHeight,
          ),
          boxWidth,
          this.timelineCoordinateCalculator.adjustPixelScale(histHeight),
        );
        bottomHeight += histHeight;
      }
      currentBucketLeftTime += secondStride;
      bucketIndex += 1;
    }
  }

  private countLogsBySeverity(
    logs: LogEntry[],
    leftMostTime: number,
    endTime: number,
    windowTimeWidth: number,
  ): SeverityLogGroupingResult {
    // Count logs by severity with spliting timerange by secondStride
    const logCounts: { [severity: string]: number }[] = [];
    let maxLogCountInSingleBucket = 0;
    let currentLogIndex = 0;
    let currentLeftTime = leftMostTime;
    for (
      ;
      currentLogIndex < logs.length &&
      logs[currentLogIndex].time < currentLeftTime;
      currentLogIndex++
    ) {
      // Skip non contained element for left most bucket
      // TODO: Use bisect to optimize
    }
    while (
      this.timelineCoordinateCalculator.timeToLeftOffset(currentLeftTime) <
      endTime
    ) {
      const logCountsInBucket: { [severity: string]: number } = {};
      for (const severity of generated.severities) {
        logCountsInBucket[severity] = 0;
      }
      let logsInBucket = 0;

      for (
        ;
        currentLogIndex < logs.length &&
        logs[currentLogIndex].time < currentLeftTime + windowTimeWidth;
        currentLogIndex++
      ) {
        const log = logs[currentLogIndex];
        logCountsInBucket[log.severity] += 1;
        logsInBucket++;
      }
      logCounts.push(logCountsInBucket);
      currentLeftTime += windowTimeWidth;
      if (logsInBucket > maxLogCountInSingleBucket) {
        maxLogCountInSingleBucket = logsInBucket;
      }
    }
    return {
      logsBySeverity: logCounts,
      maxLogCountInSingleBucket,
    };
  }

  private _calcStrideSizeAt(strideOrder: number) {
    const width = this.canvas.width;
    const maximumAllowedLargest = 10;
    for (let i = 0; i + strideOrder < STRIDE_SIZE_CANDIDATES.length; i++) {
      if (
        this.timelineCoordinateCalculator.widthToDuration(width) /
          STRIDE_SIZE_CANDIDATES[i] <=
        maximumAllowedLargest
      ) {
        return STRIDE_SIZE_CANDIDATES[i + strideOrder];
      }
    }
    return STRIDE_SIZE_CANDIDATES[STRIDE_SIZE_CANDIDATES.length - 1];
  }

  private _drawSelectedLogLine() {
    this.ctx.strokeStyle = '#33CC33';
    this.ctx.lineWidth = 4;
    this.ctx.beginPath();
    this.ctx.moveTo(
      this.timelineCoordinateCalculator.timeToLeftOffset(
        this.selectedLogTimestamp,
      ),
      this.timelineCoordinateCalculator.adjustPixelScale(20),
    );
    this.ctx.lineTo(
      this.timelineCoordinateCalculator.timeToLeftOffset(
        this.selectedLogTimestamp,
      ),
      this.timelineCoordinateCalculator.adjustPixelScale(60),
    );
    this.ctx.stroke();
  }

  private toNearestMidnightBefore(time: number) {
    const hourInMs = 1000 * 60 * 60;
    const date = this.timeToDate(time);
    const nearDate = Date.UTC(
      date.getUTCFullYear(),
      date.getUTCMonth(),
      date.getUTCDate(),
    );
    const nearDateInCurrentTimezone = new Date(
      nearDate - this._timezoneShift$.value * hourInMs,
    );
    return nearDateInCurrentTimezone;
  }

  private toTimeLabelString(time: number) {
    const date = this.timeToDate(time);
    const hour = ('' + date.getUTCHours()).padStart(2, '0');
    const minute = ('' + date.getUTCMinutes()).padStart(2, '0');
    const second = ('' + date.getUTCSeconds()).padStart(2, '0');
    return `${hour}:${minute}:${second}`;
  }

  private toDateLabelString(time: number) {
    const date = this.timeToDate(time);
    return `${date.getUTCFullYear()}/${
      date.getUTCMonth() + 1
    }/${date.getUTCDate()}`;
  }

  private timeToDate(time: number) {
    return new Date(time + this._timezoneShift$.value * 60 * 60 * 1000);
  }

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
}
