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
import { PrivateKHIExtension } from 'src/app/extensions/private/module';
import { PublicKHIExtension } from 'src/app/extensions/public/module';
import { links } from './private-links';
import {
  DownloadEnvironmentConfig,
  UploadEnvironmentConfig,
} from './environment-types';

export const environment = {
  production: true,
  useBinaryFormat: true,
  bugReportUrl:
    'https://b.corp.google.com/issues/new?component=1265687&template=1747079',
  documentUrl: 'http://go/khi',
  upload: {
    chunkSizeBytes: 512 * 1024,
    maxConcurrency: 128,
  } as UploadEnvironmentConfig,
  download: {
    chunkSizeBytes: 512 * 1024,
    maxConcurrency: 16,
  } as DownloadEnvironmentConfig,
  pluginModules: [PrivateKHIExtension, PublicKHIExtension],
  options: {
    GTAG_ID: 'G-JJ6G0C6V06',
    VIEWER_MODE: true,
  } as Record<string, unknown>,
  links,
  /*
   *  KHI Hosted is behind of GAE standard to put it behind of Uber-Proxy.
   *  GAE standard doesn't support server side streaming feature.
   *  So we need to use polling legacy for KHI Hosted.
   */
  usePollingLegacy: true,
};
