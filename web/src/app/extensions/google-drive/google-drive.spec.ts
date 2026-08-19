/**
 * Copyright 2026 Google LLC
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

import { TestBed } from '@angular/core/testing';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { of } from 'rxjs';
import { GoogleDriveDataLoaderService } from './google-drive';
import { GoogleDriveAPI } from './google-drive-api';
import { InspectionDataLoaderService } from 'src/app/services/data-loader.service';
import {
  PROGRESS_DIALOG_STATUS_UPDATOR,
  ProgressDialogStatusUpdator,
} from 'src/app/services/progress/progress-interface';

describe('GoogleDriveDataLoaderService', () => {
  let service: GoogleDriveDataLoaderService;
  let dialogSpy: jasmine.SpyObj<MatDialog>;
  let snackBarSpy: jasmine.SpyObj<MatSnackBar>;
  let driveAPISpy: jasmine.SpyObj<GoogleDriveAPI>;
  let loaderServiceSpy: jasmine.SpyObj<InspectionDataLoaderService>;
  let progressSpy: jasmine.SpyObj<ProgressDialogStatusUpdator>;
  let dialogRefSpy: jasmine.SpyObj<MatDialogRef<unknown>>;

  beforeEach(() => {
    dialogRefSpy = jasmine.createSpyObj<MatDialogRef<unknown>>('MatDialogRef', [
      'afterClosed',
    ]);
    dialogRefSpy.afterClosed.and.returnValue(of(undefined));

    dialogSpy = jasmine.createSpyObj<MatDialog>('MatDialog', ['open']);
    dialogSpy.open.and.returnValue(dialogRefSpy);

    snackBarSpy = jasmine.createSpyObj<MatSnackBar>('MatSnackBar', ['open']);
    driveAPISpy = jasmine.createSpyObj<GoogleDriveAPI>('GoogleDriveAPI', [
      'getFileAsText',
    ]);
    loaderServiceSpy = jasmine.createSpyObj<InspectionDataLoaderService>(
      'InspectionDataLoaderService',
      ['importInspectionFile', 'loadInspectionDataFromBackend'],
    );
    progressSpy = jasmine.createSpyObj<ProgressDialogStatusUpdator>(
      'ProgressDialogStatusUpdator',
      ['show', 'updateProgress', 'dismiss'],
    );

    TestBed.configureTestingModule({
      providers: [
        GoogleDriveDataLoaderService,
        { provide: MatDialog, useValue: dialogSpy },
        { provide: MatSnackBar, useValue: snackBarSpy },
        { provide: GoogleDriveAPI, useValue: driveAPISpy },
        {
          provide: InspectionDataLoaderService,
          useValue: loaderServiceSpy,
        },
        { provide: PROGRESS_DIALOG_STATUS_UPDATOR, useValue: progressSpy },
      ],
    });

    service = TestBed.inject(GoogleDriveDataLoaderService);
  });

  it('downloads file from drive and imports to backend', async () => {
    const dummyBuffer = new ArrayBuffer(8);

    driveAPISpy.getFileAsText.and.returnValue(Promise.resolve(dummyBuffer));
    loaderServiceSpy.importInspectionFile.and.returnValue(Promise.resolve());

    await service.load('test-file-id');

    expect(dialogSpy.open).toHaveBeenCalled();
    expect(progressSpy.show).toHaveBeenCalled();
    expect(driveAPISpy.getFileAsText).toHaveBeenCalledWith('test-file-id');
    expect(loaderServiceSpy.importInspectionFile).toHaveBeenCalledWith(
      jasmine.any(File),
    );
    expect(progressSpy.dismiss).toHaveBeenCalled();
  });

  it('shows snackbar when drive download fails', async () => {
    driveAPISpy.getFileAsText.and.returnValue(
      Promise.reject(new Error('Network error')),
    );

    await service.load('invalid-file-id');

    expect(snackBarSpy.open).toHaveBeenCalledWith(
      'Specified inspection data not found',
      'Close',
      jasmine.any(Object),
    );
    expect(progressSpy.dismiss).toHaveBeenCalled();
  });
});
