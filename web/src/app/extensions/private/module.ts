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

import { NgModule } from '@angular/core';
import { KHIExtensionBundle } from 'src/app/extensions/extension-common/extension';
import { AnalyticsLifecycleExtension } from 'src/app/extensions/private/analytics/analytics-lifecycle-extension';
import { PrivateAnalyticsService } from 'src/app/extensions/private/analytics/private-analytics.service';
import { FRONTEND_ANALYTICS } from 'src/app/extensions/private/analytics/types';

@NgModule({
  imports: [],
  providers: [
    PrivateAnalyticsService,
    { provide: FRONTEND_ANALYTICS, useExisting: PrivateAnalyticsService },
    KHIExtensionBundle.forExtension(initExtension),
  ],
})
export class PrivateKHIExtension {}

function initExtension(extension: KHIExtensionBundle) {
  extension.addLifecycleHookExtension(AnalyticsLifecycleExtension);
}
