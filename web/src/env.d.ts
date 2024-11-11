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

/* eslint-disable */
declare var process: {
  env: {
    NG_APP_BACKEND_URL_PREFIX: string;
    NG_APP_ENV: string;
    NG_APP_VIEWER_MODE: string;
    NG_APP_VERSION: string;
    NG_APP_GOOGLE_DRIVE_CLIENT_ID: string;
    NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER: string;
    NG_APP_REPORT_BUG_URL: string;
    NG_APP_GRAPH_PAGE: string;
    NG_APP_GTAG_ID: string;
    NG_APP_DOCUMENT_URL: string;
    // Replace the line below with your environment variable for better type checking
    [key: string]: any;
  };
};
