import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Observable, ReplaySubject } from 'rxjs';
import { ProgressDialogComponent } from '../../dialogs/progress/progress.component';
import {
  CurrentProgress,
  PROGRESS_DIALOG_STATUS_OBSERVER,
  PROGRESS_DIALOG_STATUS_UPDATOR,
  ProgressDialogStatusObserver,
  ProgressDialogStatusUpdator,
} from './progress-interface';
import { Injectable, Provider } from '@angular/core';

@Injectable()
export class ProgressDialogService
  implements ProgressDialogStatusObserver, ProgressDialogStatusUpdator
{
  /**
   * Returns Angular providers for interfaces implemented on this class.
   */
  public static providers(): Provider[] {
    return [
      {
        provide: PROGRESS_DIALOG_STATUS_OBSERVER,
        useClass: ProgressDialogService,
      },
      {
        provide: PROGRESS_DIALOG_STATUS_UPDATOR,
        useExisting: PROGRESS_DIALOG_STATUS_OBSERVER,
      },
    ];
  }

  private currentStatusSubject = new ReplaySubject<CurrentProgress>(1);

  /**
   * The dialog refernece to currently shown dialog.
   * It's null when no progress dialog is showing up now.
   */
  private dialogRef: MatDialogRef<ProgressDialogComponent> | null = null;

  /**
   * The counter of show() call.
   * The dismiss() function will close the dialog only when this count become 0.
   */
  private showCount = 0;

  constructor(private _dialog: MatDialog) {}

  status(): Observable<CurrentProgress> {
    return this.currentStatusSubject;
  }

  show(): void {
    this.currentStatusSubject.next({
      message: '',
      percent: 0,
      mode: 'indeterminate',
    });
    if (this.showCount == 0) {
      this.dialogRef = this._dialog.open(ProgressDialogComponent, {
        disableClose: true,
      });
    }
    this.showCount += 1;
  }

  updateProgress(progress: CurrentProgress): void {
    this.currentStatusSubject.next(progress);
  }

  dismiss(): void {
    this.showCount -= 1;
    // Close dialog only when user calls dismiss same times of calling show
    if (this.showCount == 0) {
      if (this.dialogRef) {
        this.dialogRef.close();
        this.dialogRef = null;
      }
      this.currentStatusSubject.next({
        message: '',
        percent: 0,
        mode: 'determinate',
      });
    }
  }
}
