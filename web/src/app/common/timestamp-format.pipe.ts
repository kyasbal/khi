import { Pipe, PipeTransform } from '@angular/core';
import { map, Observable } from 'rxjs';
import { ViewStateService } from '../services/view-state.service';

@Pipe({
  name: 'tsf',
})
export class TimestampFormatPipe implements PipeTransform {
  constructor(private readonly _viewStateService: ViewStateService) {}

  transform(unix_time: number): Observable<string> {
    return this._viewStateService.timezoneShift.pipe(
      map((t) => TimestampFormatPipe.toShortDisplayTimestamp(unix_time, t)),
    );
  }
  public static toShortDisplayTimestamp(time: number, timezoneShift: number) {
    const date = new Date(time + timezoneShift * 60 * 60 * 1000);
    const hour = ('' + date.getUTCHours()).padStart(2, '0');
    const minute = ('' + date.getUTCMinutes()).padStart(2, '0');
    const second = ('' + date.getUTCSeconds()).padStart(2, '0');
    return `${hour}:${minute}:${second}`;
  }
}

@Pipe({
  name: 'tsflong',
})
export class LongTimestampFormatPipe implements PipeTransform {
  constructor(private readonly _viewStateService: ViewStateService) {}

  transform(unix_time: number): Observable<string> {
    return this._viewStateService.timezoneShift.pipe(
      map((t) => LongTimestampFormatPipe.toLongDisplayTimestamp(unix_time, t)),
    );
  }
  public static toLongDisplayTimestamp(time: number, timezoneShift: number) {
    const shiftHour = Math.floor(timezoneShift);
    const shiftMinute = Math.floor((timezoneShift % 1) * 60);
    const shiftHourStr = ('' + Math.abs(shiftHour)).padStart(2, '0');
    const shiftMinuteStr = ('' + Math.abs(shiftMinute)).padStart(2, '0');
    const sign = timezoneShift >= 0 ? '+' : '-';
    const shiftLabel = `${sign}${shiftHourStr}:${shiftMinuteStr}`;
    const date = new Date(time + timezoneShift * 60 * 60 * 1000);
    const year = '' + date.getUTCFullYear();
    const month = ('' + (date.getUTCMonth() + 1)).padStart(2, '0');
    const day = ('' + date.getUTCDate()).padStart(2, '0');
    const hour = ('' + date.getUTCHours()).padStart(2, '0');
    const minute = ('' + date.getUTCMinutes()).padStart(2, '0');
    const second = ('' + date.getUTCSeconds()).padStart(2, '0');
    const milliSec = ('' + date.getUTCMilliseconds()).padStart(3, '0');
    return `${year}-${month}-${day}T${hour}:${minute}:${second}.${milliSec}${shiftLabel}`;
  }
}
