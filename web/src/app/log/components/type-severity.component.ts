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

import { Component, input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';

/**
 * `TypeSeverityComponent` displays a visual badge representing the severity and type of a log entry.
 * It uses predefined CSS rules to apply semantic colors based on the severity level
 * (e.g., info, warning, error, fatal).
 */
@Component({
  selector: 'khi-type-severity',
  standalone: true,
  templateUrl: './type-severity.component.html',
  styleUrls: ['./type-severity.component.scss'],
  imports: [CommonModule, MatIconModule, MatTooltipModule],
})
export class TypeSeverityComponent {
  /**
   * The type of the log entry (e.g., 'k8s_audit', 'k8s_container').
   * Displayed distinctly alongside the severity.
   */
  logType = input('N/A');

  /**
   * The severity level of the log entry (e.g., 'INFO', 'WARNING', 'ERROR').
   * Determines the semantic color of the displayed badge.
   */
  severity = input('N/A');
}
