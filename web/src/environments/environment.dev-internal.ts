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

// This file can be replaced during build by using the `fileReplacements` array.
// `ng build` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

import { GoogleDriveDataLoaderKHIExtension } from 'src/app/extensions/google-drive/module';
import { PrivateKHIExtension } from 'src/app/extensions/private/module';
import { PublicKHIExtension } from 'src/app/extensions/public/module';
import { links } from './documents';

export const environment = {
  production: false,
  viewerMode: false,
  bugReportUrl:
    'https://b.corp.google.com/issues/new?component=1265687&template=1747079',
  documentUrl: 'http://go/khi',
  pluginModules: [
    PrivateKHIExtension,
    PublicKHIExtension,
    GoogleDriveDataLoaderKHIExtension,
  ],
  options: {
    GTAG_ID: 'G-JJ6G0C6V06',
  } as Record<string, unknown>,
  links,
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.
