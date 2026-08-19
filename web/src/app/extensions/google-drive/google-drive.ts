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
import { firstValueFrom } from 'rxjs';
import { GoogleDriveAPI } from './google-drive-api';
import {
  PROGRESS_DIALOG_STATUS_UPDATOR,
  ProgressDialogStatusUpdator,
} from 'src/app/services/progress/progress-interface';
import { InspectionDataLoaderService } from 'src/app/services/data-loader.service';
import { LoginDialogComponent } from './login-dialog/login.component';

/**
 * Service for loading inspection data from Google Drive and importing it via backend.
 */
@Injectable()
export class GoogleDriveDataLoaderService {
  private readonly progress = inject<ProgressDialogStatusUpdator>(
    PROGRESS_DIALOG_STATUS_UPDATOR,
  );
  private readonly loaderService = inject(InspectionDataLoaderService);
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);
  private readonly driveAPI = inject(GoogleDriveAPI);

  /**
   * Loads inspection data from Google Drive by file ID and imports it into the backend.
   *
   * @param fileId Google Drive file ID.
   */
  public async load(fileId: string): Promise<void> {
    const dialogRef = this.dialog.open(LoginDialogComponent, {});
    await firstValueFrom(dialogRef.afterClosed());

    this.progress.show();
    this.progress.updateProgress({
      message: 'Loading inspection data from Google Drive...',
      mode: 'indeterminate',
      percent: 0,
    });
    try {
      const fileData = await this.driveAPI.getFileAsText(fileId);
      const file = new File([fileData], `${fileId}.khi`);
      await this.loaderService.importInspectionFile(file);
    } catch (e) {
      console.error(e);
      this.snackBar.open('Specified inspection data not found', 'Close', {
        duration: 10000,
      });
    } finally {
      this.progress.dismiss();
    }
  }
}
