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
import { Component, input, inject } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { map, shareReplay } from 'rxjs';
import { SelectionManagerService } from 'src/app/services/selection-manager.service';

/**
 * Represents a view model for a single resource reference link.
 */
export interface ResourceRefAnnotationViewModel {
  label: string;
  path: string;
}

/**
 * `ResourceReferenceListComponent` renders a list of related resources extracted from a loaded log.
 * It displays clickable chips that allow the user to highlight or select specific timelines
 * directly from the log details view.
 */
@Component({
  selector: 'khi-resource-reference-list',
  standalone: true,
  templateUrl: './resource-reference-list.component.html',
  styleUrl: './resource-reference-list.component.scss',
  imports: [CommonModule, MatIconModule],
})
export class ResourceReferenceListComponent {
  private readonly selectionManager = inject(SelectionManagerService);

  /**
   * A list of resolved resource references to display.
   */
  refs = input<ResourceRefAnnotationViewModel[]>([]);

  /**
   * Observable tracking the currently selected timeline path to visually indicate selection state.
   */
  currentSelectedTimelinePath = this.selectionManager.selectedTimeline.pipe(
    map((t) => t?.resourcePath ?? ''),
    shareReplay(1),
  );

  /**
   * Select the resource at the resource path.
   */
  public selectResource(resourcePath: string) {
    this.selectionManager.onSelectTimeline(resourcePath);
  }

  /**
   * Highlight the resource at the resource path.
   */
  public highlightResource(resourcePath: string) {
    this.selectionManager.onHighlightTimeline(resourcePath);
  }
}
