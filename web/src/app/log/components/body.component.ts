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

import { Component, computed, inject, input, resource } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LogHeaderComponent } from './header.component';
import { HighlightModule } from 'ngx-highlightjs';
import { ResolveTextPipe } from '../../common/resolve-text.pipe';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltip } from '@angular/material/tooltip';
import { LogEntry } from 'src/app/store/log';
import { KHIIconRegistrationModule } from 'src/app/shared/module/icon-registration.module';
import { MatButtonModule } from '@angular/material/button';
import { InspectionDataStoreService } from 'src/app/services/inspection-data-store.service';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { Clipboard, ClipboardModule } from '@angular/cdk/clipboard';
import { filter, firstValueFrom, of } from 'rxjs';
import jsyaml from 'js-yaml';
import { toSignal } from '@angular/core/rxjs-interop';

/**
 * Component responsible for displaying the detailed body of a log entry.
 * Provides actions such as copying the raw log content and copying a Cloud Logging query
 * for the specific log entry.
 */
@Component({
  selector: 'khi-log-body',
  templateUrl: './body.component.html',
  styleUrls: ['./body.component.scss'],
  imports: [
    CommonModule,
    LogHeaderComponent,
    HighlightModule,
    ResolveTextPipe,
    MatIconModule,
    MatTooltip,
    KHIIconRegistrationModule,
    MatButtonModule,
    MatSnackBarModule,
    ClipboardModule,
  ],
})
export class LogBodyComponent {
  private readonly clipboard = inject(Clipboard);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dataStore = inject(InspectionDataStoreService, {
    optional: true,
  });

  /**
   * The log entry model to display.
   */
  public log = input<LogEntry | null>(null);

  /**
   * Signal containing the current text reference resolver from the data store.
   */
  private readonly referenceResolver = toSignal(
    this.dataStore?.referenceResolver.pipe(filter((tb) => !!tb)) ?? of(null),
  );

  /**
   * Asynchronously loads the full log body text using the reference resolver.
   */
  private readonly logBody = resource({
    params: () => ({ resolver: this.referenceResolver(), log: this.log() }),
    loader: ({ params }) => {
      if (!params.log || !params.resolver) {
        return Promise.resolve('');
      }
      return firstValueFrom(params.resolver.getText(params.log.body));
    },
  });

  /**
   * Determines if the "Copy Query" button should be visible.
   * True only if a valid timestamp can be extracted from the loaded log body.
   */
  protected readonly showCopyQueryButton = computed(() => {
    return this.getTimestampStringFromBody() !== null;
  });

  /**
   * Copies the loaded log body text to the clipboard and displays a notification.
   */
  copyLog() {
    const hasValue = this.logBody.hasValue();
    if (!hasValue) {
      return;
    }
    this.showCopySnackbarMessage(this.clipboard.copy(this.logBody.value()));
  }

  /**
   * Copies a Cloud Logging query string uniquely identifying this log to the clipboard.
   * Extracts the insertId and timestamp from the log body to build the query.
   */
  copyLogQuery() {
    const log = this.log();
    if (!log) {
      return;
    }
    const timestampString = this.getTimestampStringFromBody();
    if (!timestampString) {
      return;
    }
    this.showCopySnackbarMessage(
      this.clipboard.copy(`(
-- Log query for "${log.summary}"
insertId="${log.insertId}" 
timestamp="${timestampString}"
)`),
    );
  }

  /**
   * Extracts the timestamp string from the loaded YAML/JSON log body.
   * @returns The extracted timestamp string, or null if it cannot be parsed or found.
   */
  private getTimestampStringFromBody(): string | null {
    const hasValue = this.logBody.hasValue();
    const log = this.log();
    if (!hasValue || !log) {
      return null;
    }
    try {
      const parsedLog = jsyaml.load(this.logBody.value()) as {
        [key: string]: string;
      };
      return parsedLog['timestamp'] ?? null;
    } catch {
      return null;
    }
  }

  /**
   * Displays a snackbar notification indicating the result of a copy action.
   * @param success Whether the copy to clipboard operation was successful.
   */
  private showCopySnackbarMessage(success: boolean) {
    this.snackBar.open(success ? 'Copied!' : 'Copy failed', undefined, {
      duration: 1000,
    });
  }
}
