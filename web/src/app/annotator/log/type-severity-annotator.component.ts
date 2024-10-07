import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { AnnotationDecider } from '../annotator';
import { LogTypeMetadata } from 'src/app/generated';
import { LogEntry } from 'src/app/store/log';

@Component({
  standalone: true,
  templateUrl: './type-severity-annotator.component.html',
  styleUrls: ['./type-severity-annotator.component.sass'],
  imports: [CommonModule, MatIconModule],
})
export class TypeSeverityAnnotatorComponent {
  @Input()
  logType = 'N/A';

  @Input()
  severity = 'N/A';

  public static inputMapper: AnnotationDecider<LogEntry> = (
    l?: LogEntry | null,
  ) => {
    let logType = 'N/A';
    if (l !== null && l !== undefined) {
      logType = LogTypeMetadata[l.logType].label;
    }
    return {
      inputs: {
        logType: logType,
        severity: l?.logSeverityLabel ?? 'N/A',
      },
    };
  };
}
