import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { NEVER, Observable } from 'rxjs';
import { AnnotationDecision } from './annotator';

export interface ReferenceViewModel {
  url: string;
  displayText: string;
}

@Component({
  standalone: true,
  imports: [CommonModule, MatIconModule, MatTooltipModule],
  templateUrl: './common-reference-list.component.html',
  styleUrl: './common-reference-list.component.sass',
})
export class CommonReferenceListComponent {
  @Input()
  references: Observable<ReferenceViewModel[]> = NEVER;

  @Input()
  isGray: boolean = false;

  @Input()
  header: string = 'References';

  onLinkClick(value: ReferenceViewModel) {
    window.open(value.url, '_blank');
  }
}

export interface CommonReferenceListAnnotatorDecision
  extends AnnotationDecision {
  inputs: {
    header: string;
    isGray: boolean;
    references: Observable<ReferenceViewModel[]>;
  };
}
