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

import { Component, computed, input, output } from '@angular/core';
import { DagPositionedNode } from 'src/app/pages/task-graph-debug/components/dag-viewer/dag-viewer.model';

/**
 * Renders an individual task node card inside the SVG DAG canvas.
 */
@Component({
  selector: 'g[khi-dag-node]',
  templateUrl: './dag-node.component.html',
  styleUrls: ['./dag-node.component.scss'],
  host: {
    'class': 'dag-node-group',
    '[class.selected]': 'isSelected()',
    '[class.highlighted]': 'isHighlighted()',
    '[class.dimmed]': 'isDimmed()',
    '[class.feature]': 'node().isFeature',
    '[class.initial-task]': 'node().isInitialTask',
    '[attr.transform]': 'transform()',
    '(click)': 'onClick($event)',
  },
})
export class DagNodeComponent {
  /**
   * Positioned node containing task metadata and coordinates.
   */
  readonly node = input.required<DagPositionedNode>();

  /**
   * Whether this node is currently selected by the user.
   */
  readonly isSelected = input<boolean>(false);

  /**
   * Whether this node is part of the highlighted upstream/downstream dependency path.
   */
  readonly isHighlighted = input<boolean>(false);

  /**
   * Whether this node is dimmed because another node is active.
   */
  readonly isDimmed = input<boolean>(false);

  /**
   * Emits when the user clicks on this node.
   */
  readonly selectNode = output<DagPositionedNode>();

  /**
   * SVG transform string positioning the node group at its top-left coordinates.
   */
  readonly transform = computed(
    () => `translate(${this.node().x}, ${this.node().y})`,
  );

  /**
   * Truncated or formatted reference ID for the card header.
   */
  readonly displayReferenceId = computed(() => {
    const ref = this.node().referenceId;
    if (ref.length > 28) {
      return `...${ref.slice(-25)}`;
    }
    return ref;
  });

  /**
   * Implementation hash identifier for the secondary label.
   */
  readonly displayImplementationId = computed(() => {
    const id = this.node().id;
    const parts = id.split('#');
    if (parts.length > 1) {
      return `#${parts[1]}`;
    }
    return id;
  });

  /**
   * SVG transform for the Init badge.
   */
  readonly initBadgeTransform = computed(() => {
    const offsetX = this.node().isFeature ? 62 : 0;
    return `translate(${offsetX}, 0)`;
  });

  /**
   * SVG transform for the Priority badge.
   */
  readonly prioBadgeTransform = computed(() => {
    let offsetX = 0;
    if (this.node().isFeature) {
      offsetX += 62;
    }
    if (this.node().isInitialTask) {
      offsetX += 42;
    }
    return `translate(${offsetX}, 0)`;
  });

  /**
   * Handles click events on the node group.
   */
  onClick(event: MouseEvent): void {
    event.stopPropagation();
    this.selectNode.emit(this.node());
  }
}
