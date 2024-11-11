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

import { Injectable, Injector } from '@angular/core';
import { GoogleDriveDataLoaderService } from './google-drive';

@Injectable()
export class DataLoadSourceExtension {
  constructor(private _injector: Injector) {}

  public load(fileId: string): void {
    if (
      process.env['NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER'] &&
      process.env['NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER'] != 'false'
    ) {
      const loader = this._injector.get<GoogleDriveDataLoaderService>(
        GoogleDriveDataLoaderService,
      );
      loader.load(fileId);
    }
  }
}
