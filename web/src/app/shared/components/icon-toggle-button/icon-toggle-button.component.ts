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

import { CommonModule } from '@angular/common';
import { Component, input, model } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { KHIIconRegistrationModule } from '../../module/icon-registration.module';

/**
 * A reusable toggle button component that displays a Material icon.
 */
@Component({
  selector: 'khi-icon-toggle-button',
  templateUrl: './icon-toggle-button.component.html',
  styleUrls: ['./icon-toggle-button.component.scss'],
  imports: [
    CommonModule,
    MatTooltipModule,
    MatIconModule,
    KHIIconRegistrationModule,
  ],
})
export class IconToggleButtonComponent {
  /**
   * The name of the Material symbol/icon to display.
   */
  icon = input<string>('');

  /**
   * The text to display in the tooltip when the user hovers over the button.
   */
  tooltip = input<string>('');

  /**
   * The two-way binding for the current selection state of the button.
   * Emits the new boolean value when toggled.
   */
  selected = model<boolean>(false);

  /**
   * Whether the button is disabled. If true, interactions are ignored and the disabled visual style is applied.
   */
  disabled = input<boolean>(false);

  /**
   * Handles click events on the button, toggling the current `selected` state.
   */
  onClick() {
    if (this.disabled()) return;
    this.selected.set(!this.selected());
  }
}
