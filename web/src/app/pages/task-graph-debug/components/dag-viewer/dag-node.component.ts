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

import {
  ChangeDetectionStrategy,
  Component,
  computed,
  input,
  output,
} from '@angular/core';
import { DagPositionedNode } from 'src/app/pages/task-graph-debug/components/dag-viewer/dag-viewer.model';

const CHAR_WIDTH_ESTIMATE_PX = 7.5;
const NODE_HORIZONTAL_PADDING_PX = 28;
const FEATURE_BADGE_OFFSET_X = 62;
const INIT_BADGE_OFFSET_X = 42;

/**
 * Renders an individual task node card inside the SVG DAG canvas.
 */
@Component({
  // eslint-disable-next-line @angular-eslint/component-selector
  selector: 'g[khi-dag-node]',
  templateUrl: './dag-node.component.html',
  styleUrls: ['./dag-node.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: {
    class: 'dag-node-group',
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
   * Domain prefix extracted from the reference ID, such as "khi.google.com/".
   * Returns empty string if the reference ID does not contain a slash.
   */
  readonly displayDomain = computed(() => {
    const ref = this.node().referenceId;
    const slashIndex = ref.indexOf('/');
    if (slashIndex === -1) {
      return '';
    }
    return ref.slice(0, slashIndex + 1);
  });

  /**
   * Main task identifier after the domain prefix.
   */
  readonly displayTaskName = computed(() => {
    const ref = this.node().referenceId;
    const slashIndex = ref.indexOf('/');
    const name = slashIndex === -1 ? ref : ref.slice(slashIndex + 1);
    const maxChars = Math.max(
      10,
      Math.floor(
        (this.node().width - NODE_HORIZONTAL_PADDING_PX) /
          CHAR_WIDTH_ESTIMATE_PX,
      ),
    );
    if (name.length > maxChars) {
      return `${name.slice(0, maxChars - 3)}...`;
    }
    return name;
  });

  /**
   * Vertical coordinate (Y) for the task name SVG text element.
   */
  readonly taskTextY = computed(() => (this.displayDomain() ? 35 : 26));

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
    const offsetX = this.node().isFeature ? FEATURE_BADGE_OFFSET_X : 0;
    return `translate(${offsetX}, 0)`;
  });

  /**
   * SVG transform for the Priority badge.
   */
  readonly prioBadgeTransform = computed(() => {
    let offsetX = 0;
    if (this.node().isFeature) {
      offsetX += FEATURE_BADGE_OFFSET_X;
    }
    if (this.node().isInitialTask) {
      offsetX += INIT_BADGE_OFFSET_X;
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
