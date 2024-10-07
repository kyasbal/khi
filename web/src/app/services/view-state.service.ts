import { Injectable } from '@angular/core';
import {
  BehaviorSubject,
  Observable,
  Subject,
  animationFrames,
  distinctUntilChanged,
  map,
  shareReplay,
} from 'rxjs';

/**
 * A service to manage statuses used for view in application wide.
 */
@Injectable({ providedIn: 'root' })
export class ViewStateService {
  /**
   * Rendering small sharp shapes with WebGL in the default pixel ratio can be blurry result in high resolution display like Mac's retina display.
   * Scale the value by 1.5 by default. TODO: change this value with regarding the performance.
   */
  public static DEVICE_PIXEL_RATIO_SCALE = 1.5;

  private timezoneShiftSubject: BehaviorSubject<number> = new BehaviorSubject(
    -new Date().getTimezoneOffset() / 60,
  );

  /**
   * Number of the hours differences from UTC
   */
  public timezoneShift: Observable<number> = this.timezoneShiftSubject;

  private timelineStateResetCommandSubject: Subject<null> = new Subject();

  /**
   * Emit value when timeline scale and offset reset was requesed.
   */
  public timelineStateResetCommand: Observable<null> =
    this.timelineStateResetCommandSubject;

  private timeOffsetSubject = new BehaviorSubject(0);

  public timeOffset: Observable<number> = this.timeOffsetSubject;

  private pixelPerTimeSubject = new BehaviorSubject<number>(0.01);

  public pixelPerTime: Observable<number> = this.pixelPerTimeSubject;

  private visibleWidthSubject = new BehaviorSubject(100);

  public visibleWidth: Observable<number> = this.visibleWidthSubject.pipe(
    shareReplay(1),
  );

  public devicePixelRatio = animationFrames().pipe(
    map(
      () => window.devicePixelRatio * ViewStateService.DEVICE_PIXEL_RATIO_SCALE,
    ),
    distinctUntilChanged(),
    shareReplay(1),
  );

  public setTimezoneShift(timezoneShift: number): void {
    this.timezoneShiftSubject.next(timezoneShift);
  }

  /**
   * Set the offset of timeline view
   * @param offset offset time of left most edge of the timeline view
   */
  public setTimeOffset(offset: number): void {
    this.timeOffsetSubject.next(offset);
  }

  public getTimeOffset(): number {
    return this.timeOffsetSubject.value;
  }

  public setPixelPerTime(pixelPerTime: number): void {
    return this.pixelPerTimeSubject.next(pixelPerTime);
  }

  public getPixelPerTime(): number {
    return this.pixelPerTimeSubject.value;
  }

  public setVisibleWidth(widthInPx: number): void {
    this.visibleWidthSubject.next(widthInPx);
  }

  public getVisibleWidth(): number {
    return this.visibleWidthSubject.value;
  }

  /**
   * Reset scale and offset of timeline status.
   */
  public resetTimelineStatus(): void {
    this.timelineStateResetCommandSubject.next(null);
  }
}
