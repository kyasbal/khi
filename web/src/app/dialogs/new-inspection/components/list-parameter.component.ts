/**
 * Copyright 2026 Google LLC
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

import { Component, computed, forwardRef, inject, input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import {
  ListParameterFormField,
  ListParameterFormFieldItem,
  ParameterHintType,
  ParameterInputType,
} from 'src/app/common/schema/form-types';
import { KHIIconRegistrationModule } from 'src/app/shared/module/icon-registration.module';
import { PARAMETER_STORE } from 'src/app/dialogs/new-inspection/components/service/parameter-store';
import { ParameterHeaderComponent } from 'src/app/dialogs/new-inspection/components/parameter-header.component';
import { ParameterHintComponent } from 'src/app/dialogs/new-inspection/components/parameter-hint.component';
import { TextParameterComponent } from 'src/app/dialogs/new-inspection/components/text-parameter.component';
import { FileParameterComponent } from 'src/app/dialogs/new-inspection/components/file-parameter.component';
import { SetParameterComponent } from 'src/app/dialogs/new-inspection/components/set-parameter.component';
import { CheckboxParameterComponent } from 'src/app/dialogs/new-inspection/components/checkbox-parameter.component';
import { GroupParameterComponent } from 'src/app/dialogs/new-inspection/components/group-parameter.component';

/**
 * A repeatable list of parameter form fields.
 */
@Component({
  selector: 'khi-new-inspection-list-parameter',
  templateUrl: './list-parameter.component.html',
  styleUrls: ['./list-parameter.component.scss'],
  imports: [
    CommonModule,
    MatButtonModule,
    MatIconModule,
    KHIIconRegistrationModule,
    MatTooltipModule,
    ParameterHeaderComponent,
    ParameterHintComponent,
    TextParameterComponent,
    FileParameterComponent,
    SetParameterComponent,
    CheckboxParameterComponent,
    forwardRef(() => GroupParameterComponent),
  ],
})
export class ListParameterComponent {
  /**
   * Exposes ParameterInputType enum to the template.
   */
  protected readonly ParameterInputType = ParameterInputType;

  /**
   * Exposes ParameterHintType enum to the template.
   */
  protected readonly ParameterHintType = ParameterHintType;

  /**
   * List parameter form field definition.
   */
  readonly parameter = input.required<ListParameterFormField>();

  private readonly store = inject(PARAMETER_STORE);

  private nextKeyCounter = 0;

  /**
   * Computed signal of currently active item keys.
   */
  readonly activeKeys = computed<string[]>(() => {
    const stored = this.store.get<string[]>(this.parameter().id)();
    if (Array.isArray(stored)) {
      return stored;
    }
    return this.parameter().items.map((item) => item.key);
  });

  /**
   * Computed signal of items that are currently active in the list.
   */
  readonly activeItems = computed<ListParameterFormFieldItem[]>(() => {
    const active = new Set(this.activeKeys());
    return this.parameter().items.filter((item) => active.has(item.key));
  });

  /**
   * Computed signal indicating whether a new item can be added.
   */
  readonly canAdd = computed<boolean>(() => {
    const max = this.parameter().maxCount;
    return max <= 0 || this.activeKeys().length < max;
  });

  /**
   * Computed signal indicating whether an item can be removed.
   */
  readonly canRemove = computed<boolean>(() => {
    const min = this.parameter().minCount;
    return this.activeKeys().length > min;
  });

  /**
   * Adds a new item key to the list parameter.
   */
  addItem(): void {
    if (!this.canAdd()) {
      return;
    }

    const allKeys = [
      ...this.activeKeys(),
      ...this.parameter().default,
      ...this.parameter().items.map((item) => item.key),
    ];

    let maxExisting = -1;
    for (const k of allKeys) {
      if (/^\d+$/.test(k)) {
        const val = parseInt(k, 10);
        if (val > maxExisting) {
          maxExisting = val;
        }
      }
    }

    const nextNum = Math.max(maxExisting + 1, this.nextKeyCounter);
    this.nextKeyCounter = nextNum + 1;
    this.store.set(this.parameter().id, [
      ...this.activeKeys(),
      String(nextNum),
    ]);
  }

  /**
   * Removes an item key from the list parameter.
   */
  removeItem(key: string): void {
    if (!this.canRemove()) {
      return;
    }
    this.store.set(
      this.parameter().id,
      this.activeKeys().filter((k) => k !== key),
    );
  }
}
