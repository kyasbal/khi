import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { NEVER, Observable, of } from 'rxjs';
import { AnnotationDecider, DECISION_HIDDEN } from './annotator';
import { ResourceRevisionChangePair } from '../store/timeline';
@Component({
  standalone: true,
  imports: [CommonModule, MatIconModule],
  templateUrl: './common-warning-message.component.html',
  styleUrl: './common-warning-message.component.sass',
})
export class CommonWarningMessageComponent {
  @Input()
  icon = 'warm';

  @Input()
  message: Observable<string> = NEVER;

  public static inputMapperForRevisionPair(
    icon: string,
    messageMapper: (source: ResourceRevisionChangePair) => string,
    predicate: (source: ResourceRevisionChangePair) => boolean,
  ): AnnotationDecider<ResourceRevisionChangePair> {
    return (source) => {
      if (!source || !predicate(source)) return DECISION_HIDDEN;
      return {
        inputs: {
          icon,
          message: of(messageMapper(source!)),
        },
      };
    };
  }
}
