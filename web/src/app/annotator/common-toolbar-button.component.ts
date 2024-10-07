import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { AnnotationDecision } from './annotator';

@Component({
  standalone: true,
  templateUrl: './common-toolbar-button.component.html',
  styleUrl: './common-toolbar-button.component.sass',
  imports: [CommonModule, MatIconModule, MatTooltipModule],
})
export class CommonToolbarButtonComponent {
  @Input()
  public icon = 'content_paste';

  @Input()
  public tooltip = '';

  @Input()
  public onClick = () => {};

  @Input()
  public disabled = false;

  public triggerOnClock() {
    this.onClick();
  }

  public static disabledAnnotationDecision(
    icon: string,
    tooltip: string,
  ): CommonToolbarButtonInput {
    return {
      hidden: true,
      inputs: {
        icon,
        tooltip,
        disabled: true,
        onClick: () => {},
      },
    };
  }
}

export interface CommonToolbarButtonInput extends AnnotationDecision {
  inputs: {
    icon: string;
    tooltip: string;
    disabled: boolean;
    onClick: () => void;
  };
}
