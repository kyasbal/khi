import { ViewStateService } from '../services/view-state.service';

/**
 * Calculator utility used for converting pixels - times
 */
export class TimelinenCoordinateCalculator {
  constructor(private viewState: ViewStateService) {
    this.viewState.devicePixelRatio.subscribe(
      (pixelRatio) => (this.canvasPixelRatio = pixelRatio),
    );
  }

  public canvasPixelRatio = 1;

  private _pixelScale(asCanvasCoordinate: boolean) {
    return asCanvasCoordinate ? this.canvasPixelRatio : 1;
  }

  public adjustPixelScale(size: number): number {
    return size * this.canvasPixelRatio;
  }

  public leftMostTime(): number {
    return this.viewState.getTimeOffset();
  }

  // TODO: buggy. The value somehow become smaller than the expected value.
  // Multiply 4 as a temporary fix.
  public rightMostTime(): number {
    return (
      this.leftMostTime() +
      this.widthToDuration(this.viewState.getVisibleWidth() * 4)
    );
  }

  public timeToLeftOffset(time: number, canvasCoordinate = true): number {
    const offsetFromLeftTime = time - this.leftMostTime();
    return this.durationToWidth(offsetFromLeftTime, canvasCoordinate);
  }

  public timeToRightOffset(time: number, canvasCoodinate = true): number {
    const timeToRight = this.rightMostTime() - time;
    return this.durationToWidth(timeToRight, canvasCoodinate);
  }

  public durationToWidth(duration: number, canvasCoodinate = true): number {
    return (
      duration *
      this.viewState.getPixelPerTime() *
      this._pixelScale(canvasCoodinate)
    );
  }

  public widthToDuration(width: number, canvasCoodinate = true): number {
    return (
      width /
      this._pixelScale(canvasCoodinate) /
      this.viewState.getPixelPerTime()
    );
  }

  public leftOffsetToTime(leftOffset: number, canvasCoordinate = true): number {
    const durationFromLeft = this.widthToDuration(leftOffset, canvasCoordinate);
    return this.leftMostTime() + durationFromLeft;
  }
}
