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

import { Inject, Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LoginDialogComponent } from '../../extensions/data-loader/login-dialog/login.component';
import { InspectionDataLoaderService } from '../../services/data-loader.service';
import { conditionalModule } from '../util';
import { GoogleDriveAPI } from './google-drive-api';
import {
  PROGRESS_DIALOG_STATUS_UPDATOR,
  ProgressDialogStatusUpdator,
} from 'src/app/services/progress/progress-interface';
export function RegisterGoogleDriveExtensionProvidersIfEnabled() {
  return conditionalModule(
    !!process.env['NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER'],
    GoogleDriveDataLoaderService,
  );
}

@Injectable()
export class GoogleDriveDataLoaderService {
  constructor(
    @Inject(PROGRESS_DIALOG_STATUS_UPDATOR)
    private progress: ProgressDialogStatusUpdator,
    private loaderService: InspectionDataLoaderService,
    private _dialog: MatDialog,
    private _snackBar: MatSnackBar,
    private _driveAPI: GoogleDriveAPI,
  ) {}

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
