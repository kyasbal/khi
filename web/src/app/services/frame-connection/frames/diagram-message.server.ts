/**
 * Copyright 2025 Google LLC
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

import { inject, Injectable } from '@angular/core';
import { WindowConnectorService } from '../window-connector.service';
import { Subject } from 'rxjs';
import { QUERY_DIAGRAM_DATA } from 'src/app/common/schema/inter-window-messages';

/**
 * A message server running on the main page.
 */
@Injectable({ providedIn: 'root' })
export class DiagramMessageServer {
  private destroyed = new Subject();
  private readonly connector = inject(WindowConnectorService);

  public start() {
    this.connector.receiver(QUERY_DIAGRAM_DATA).subscribe();
  }

  public stop() {
    this.destroyed.next(undefined);
  }
}
