import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl } from '@angular/forms';
import { RegexValidator } from './regex-validator';

@Component({
  selector: 'khi-header-regex-input',
  templateUrl: './regex-input.component.html',
  styleUrls: ['./regex-input.component.sass'],
})
export class RegexInputComponent {
  @Input()
  label = '';

  @Output()
  regexFilterChange: EventEmitter<string> = new EventEmitter();

  regexInput: FormControl = new FormControl('', [RegexValidator()]);

  regexFormErrorMessage(): string {
    return this.regexInput.errors!['regex'] as string;
  }

  onFilterChange() {
    if (!this.regexInput.valid) return;
    this.regexFilterChange.emit(this.regexInput.value);
  }
}
