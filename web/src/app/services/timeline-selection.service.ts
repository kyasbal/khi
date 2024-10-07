import { Injectable } from '@angular/core';
import { BehaviorSubject, filter, map } from 'rxjs';
import { InspectionDataStoreService } from './inspection-data-store.service';
import { SelectionManagerService } from './selection-manager.service';
import { LogEntry } from '../store/log';

@Injectable({ providedIn: 'root' })
export class TimelineSelectionService {
  private $logs: BehaviorSubject<LogEntry[]> =
    this._inspectionDataStore.$allLogs;

  private $currentTime: BehaviorSubject<number> = new BehaviorSubject(0);

  constructor(
    private _inspectionDataStore: InspectionDataStoreService,
    private _logSelectionManager: SelectionManagerService,
  ) {
    this._logSelectionManager.selectedLog
      .pipe(
        filter((log) => log !== null),
        map((log) => log!.time),
      )
      .subscribe(this.$currentTime);
  }

  public seek(diff: number) {
    if (diff > 0) {
      this.seekToAfter(this.$currentTime.value + diff);
    } else {
      this.seekToBefore(this.$currentTime.value + diff);
    }
  }

  public seekToBefore(time: number) {
    const logs = this.$logs.value;
    if (logs.length == 0) return;
    if (logs[logs.length - 1].time < time) {
      this._logSelectionManager.changeSelectionByLog(logs.length - 1);
    } else {
      let left = 0;
      let right = logs.length - 1;
      while (right - left > 1) {
        const mid = Math.floor((left + right) / 2);
        if (logs[mid].time <= time) {
          left = mid;
        } else {
          right = mid;
        }
      }
      this._logSelectionManager.changeSelectionByLog(left);
    }
  }

  public seekToAfter(time: number) {
    const logs = this.$logs.value;
    if (logs.length == 0) return;
    if (logs[0].time > time) {
      this._logSelectionManager.changeSelectionByLog(0);
    } else {
      let left = 0;
      let right = logs.length - 1;
      while (right - left > 1) {
        const mid = Math.floor((left + right) / 2);
        if (logs[mid].time < time) {
          left = mid;
        } else {
          right = mid;
        }
      }
      this._logSelectionManager.changeSelectionByLog(right);
    }
  }
}
