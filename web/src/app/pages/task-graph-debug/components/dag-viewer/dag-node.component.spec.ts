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

import { Component } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DagNodeComponent } from 'src/app/pages/task-graph-debug/components/dag-viewer/dag-node.component';
import { DagPositionedNode } from 'src/app/pages/task-graph-debug/components/dag-viewer/dag-viewer.model';

@Component({
  standalone: true,
  imports: [DagNodeComponent],
  template: `
    <svg>
      <g
        khi-dag-node
        [node]="node"
        [isSelected]="isSelected"
        [isHighlighted]="isHighlighted"
        [isDimmed]="isDimmed"
        (selectNode)="onSelectNode($event)"
      ></g>
    </svg>
  `,
})
class TestHostComponent {
  node: DagPositionedNode = mockNode;
  isSelected = false;
  isHighlighted = false;
  isDimmed = false;
  emittedNode?: DagPositionedNode;

  onSelectNode(node: DagPositionedNode): void {
    this.emittedNode = node;
  }
}

const mockNode: DagPositionedNode = {
  id: 'khi.test.task#1a2b3c',
  referenceId: 'khi.test.task',
  isFeature: false,
  isInitialTask: false,
  topologicalOrder: 2,
  priority: 100,
  labels: {},
  x: 120,
  y: 80,
  width: 260,
  height: 84,
  layer: 1,
};

describe('DagNodeComponent', () => {
  let hostComponent: TestHostComponent;
  let fixture: ComponentFixture<TestHostComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TestHostComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TestHostComponent);
    hostComponent = fixture.componentInstance;
  });

  it('renders node texts and positions correctly', () => {
    hostComponent.node = mockNode;
    fixture.detectChanges();

    const group = fixture.nativeElement.querySelector('.dag-node-group');
    expect(group).toBeTruthy();
    expect(group.getAttribute('transform')).toBe('translate(120, 80)');

    const refText = fixture.nativeElement.querySelector('.node-ref-text');
    expect(refText.textContent.trim()).toBe('khi.test.task');

    const implText = fixture.nativeElement.querySelector('.node-impl-text');
    expect(implText.textContent.trim()).toBe('#1a2b3c');

    const prioText = fixture.nativeElement.querySelector(
      '.prio-badge .badge-text',
    );
    expect(prioText.textContent.trim()).toBe('P100');
  });

  it('renders feature and initial task badges when enabled', () => {
    hostComponent.node = {
      ...mockNode,
      isFeature: true,
      isInitialTask: true,
    };
    fixture.detectChanges();

    const featureBadge = fixture.nativeElement.querySelector('.feature-badge');
    expect(featureBadge).toBeTruthy();

    const initBadge = fixture.nativeElement.querySelector('.init-badge');
    expect(initBadge).toBeTruthy();
  });

  it('emits selectNode when clicked', () => {
    hostComponent.node = mockNode;
    fixture.detectChanges();

    const group = fixture.nativeElement.querySelector(
      '.dag-node-group',
    ) as SVGGraphicsElement;
    group.dispatchEvent(new MouseEvent('click', { bubbles: true }));

    expect(hostComponent.emittedNode).toEqual(mockNode);
  });

  it('applies selected and dimmed CSS classes based on inputs', () => {
    hostComponent.node = mockNode;
    hostComponent.isSelected = true;
    hostComponent.isDimmed = true;
    fixture.detectChanges();

    const group = fixture.nativeElement.querySelector('.dag-node-group');
    expect(group.classList.contains('selected')).toBeTrue();
    expect(group.classList.contains('dimmed')).toBeTrue();
  });
});
