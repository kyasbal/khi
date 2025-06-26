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

import { Injectable, inject } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { GoogleDriveAPI } from './google-drive-api';
import {
  PROGRESS_DIALOG_STATUS_UPDATOR,
  ProgressDialogStatusUpdator,
} from 'src/app/services/progress/progress-interface';
import { InspectionDataLoaderService } from 'src/app/services/data-loader.service';
import { LoginDialogComponent } from './login-dialog/login.component';

@Injectable()
export class GoogleDriveDataLoaderService {
  private readonly progress = inject<ProgressDialogStatusUpdator>(
    PROGRESS_DIALOG_STATUS_UPDATOR,
  );
  private loaderService = inject(InspectionDataLoaderService);
  private _dialog = inject(MatDialog);
  private _snackBar = inject(MatSnackBar);
  private _driveAPI = inject(GoogleDriveAPI);

  public async load(fileId: string): Promise<void> {
    const dialogRef = this._dialog.open(LoginDialogComponent, {});
    dialogRef.afterClosed().subscribe({
      complete: async () => {
        this.progress.show();
        this.progress.updateProgress({
          message: 'Loading inspection data from Google Drive...',
          mode: 'indeterminate',
          percent: 0,
        });
        try {
          const fileData = await this._driveAPI.getFileAsText(fileId);
          this.loaderService.loadInspectionDataDirect(fileData);
        } catch (e) {
          console.error(e);
          this._snackBar.open('Specified inspection data not found', 'Close', {
            duration: 10000,
          });
        }
        this.progress.dismiss();
      },
    });
  }
}
