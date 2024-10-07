import { Annotator, AnnotationDecider } from '../annotator';
import { CommonToolbarButtonComponent } from '../common-toolbar-button.component';
import { LogAnnotatorResolver } from '../log/resolver';
import { inject } from '@angular/core';
import { InspectionDataStoreService } from 'src/app/services/inspection-data-store.service';
import { filter, map, of, withLatestFrom } from 'rxjs';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Clipboard } from '@angular/cdk/clipboard';
import * as jsyaml from 'js-yaml';
import { LogEntry } from 'src/app/store/log';

function copyLogEntryContentMapper(
  toolTip: string,
): AnnotationDecider<LogEntry> {
  return (l) => {
    if (!l) {
      return {
        inputs: {
          icon: 'content_paste',
          tooltip: toolTip,
          disabled: true,
          onClick: () => ({}),
        },
      };
    }
    const snackBar = inject(MatSnackBar);
    const clipboard = inject(Clipboard);
    const dataStore = inject(InspectionDataStoreService);
    return {
      inputs: {
        icon: 'content_paste',
        tooltip: toolTip,
        disabled: l && l.logIndex < 0,
        onClick: () => {
          of(l.body)
            .pipe(
              withLatestFrom(
                dataStore.textBufferSource.pipe(filter((tb) => !!tb)),
              ),
              map(([lr, tbs]) => {
                const copyText = tbs!.getText(lr);
                if (clipboard.copy(copyText)) {
                  return 'Copied!';
                } else {
                  return 'Copy failed';
                }
              }),
            )
            .subscribe((copiedMessage) => {
              snackBar.open(copiedMessage, undefined, { duration: 1000 });
            });
        },
      },
    };
  };
}

function copyLogQueryContentMapper(
  toolTip: string,
): AnnotationDecider<LogEntry> {
  return (l) => {
    if (!l) {
      return {
        // TODO: think better icon later
        inputs: {
          icon: 'markdown_paste',
          tooltip: toolTip,
          disabled: true,
          onClick: () => ({}),
        },
      };
    }
    const snackBar = inject(MatSnackBar);
    const clipboard = inject(Clipboard);
    const dataStore = inject(InspectionDataStoreService);
    return {
      inputs: {
        icon: 'markdown_paste',
        tooltip: toolTip,
        disabled: l && l.logIndex < 0,
        onClick: () => {
          of(l)
            .pipe(
              withLatestFrom(
                dataStore.textBufferSource.pipe(filter((tb) => !!tb)),
              ),
              map(([l, textSource]) => {
                const logBody = textSource!.getText(l.body);
                const parsedLog = jsyaml.load(logBody) as {
                  [key: string]: string;
                };
                const timestamp = parsedLog['timestamp'];
                return `(
-- Log query for "${l.message}"
insertId="${l.insertId}"
timestamp="${timestamp}"
)`;
              }),
            )
            .subscribe((query) => {
              let snackbarMessage = 'Copy failed';
              if (clipboard.copy(query)) {
                snackbarMessage = 'Copied!';
              }
              snackBar.open(snackbarMessage, undefined, { duration: 1000 });
            });
        },
      },
    };
  };
}
export function getDefaultLogToolAnnotatorResolver(): LogAnnotatorResolver {
  return new LogAnnotatorResolver([
    new Annotator(
      CommonToolbarButtonComponent,
      copyLogEntryContentMapper('Copy the log to clipboard'),
    ),
    new Annotator(
      CommonToolbarButtonComponent,
      copyLogQueryContentMapper('Copy query to clipboard'),
    ),
  ]);
}
