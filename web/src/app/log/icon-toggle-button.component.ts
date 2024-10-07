import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'khi-icon-toggle-button',
  templateUrl: './icon-toggle-button.component.html',
  styleUrls: ['./icon-toggle-button.component.sass'],
})
export class IconToggleButtonComponent {
  @Input()
  icon = '';

  @Input()
  tooltip = '';

  @Input()
  selected: boolean | null = false;

  @Output()
  selectedChange = new EventEmitter<boolean>();

  @Input()
  disabled: boolean | null = false;

  onClick() {
    this.selectedChange.emit(!this.selected);
  }
}
