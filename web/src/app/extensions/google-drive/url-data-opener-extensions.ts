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

import { inject } from '@angular/core';
import { URLDataOpenerExtension } from '../extension-common/extension-types/url-data-opener';
import { GoogleDriveDataLoaderService } from './google-drive';

export class GoogleDriveURLDataOpenerExtension implements URLDataOpenerExtension {
  tryOpen(): boolean {
    if (window.location.hash.length > 1) {
      const hash = window.location.hash.substring(1);
      const dataLoader = inject(GoogleDriveDataLoaderService);
      dataLoader.load(hash);
      return true;
    }
    return false;
  }
}
