import { Component, Inject } from '@angular/core';
import {
  PROGRESS_DIALOG_STATUS_OBSERVER,
  ProgressDialogStatusObserver,
} from 'src/app/services/progress/progress-interface';

@Component({
  templateUrl: './progress.component.html',
  styleUrls: ['./progress.component.sass'],
})
export class ProgressDialogComponent {
  public currentStatus = this.progressObserver.status();

  constructor(
    @Inject(PROGRESS_DIALOG_STATUS_OBSERVER)
    private progressObserver: ProgressDialogStatusObserver,
  ) {}
}
