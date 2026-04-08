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
import { appConfig } from './main';
import { ApplicationInitStatus } from '@angular/core';
import { ConsoleReporter, Reporter } from './app/common/reporter/reporter';

describe('main.ts appConfig', () => {
  it('should complete application initialization successfully', async () => {
    TestBed.configureTestingModule({
      // Provides ConsoleReporter for testing, as it is typically provided by plugins in production.
      providers: [
        { provide: Reporter, useClass: ConsoleReporter },
        ...appConfig.providers,
      ],
    });

    const status = TestBed.inject(ApplicationInitStatus);
    // Waits for all APP_INITIALIZERs to complete.
    await status.donePromise;
    expect(status.done).toBe(true);
  });
});
