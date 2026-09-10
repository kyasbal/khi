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

import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import {
  RegisteredTaskGroupInfo,
  RegisteredTaskInfo,
  TaskDependencyCardinality,
  TaskDependencyScope,
} from 'src/app/generated/api/v1/inspection_task_graph_pb';
import { Step1RegistryTableComponent } from './step1-registry-table.component';

describe('Step1RegistryTableComponent', () => {
  let component: Step1RegistryTableComponent;
  let fixture: ComponentFixture<Step1RegistryTableComponent>;

  const mockTasks: RegisteredTaskInfo[] = [
    {
      taskImplementationId: 'parser.k8s.audit-log',
      taskReferenceId: 'parser.k8s.audit',
      priority: 100,
      isFeature: false,
      isDefaultFeature: false,
      featureLabel: '',
      featureDescription: '',
      dependencies: [],
      selectorRequirements: [],
      compatibleInspectionTypes: [],
      labels: {
        log_type: 'audit',
      },
    } as unknown as RegisteredTaskInfo,
    {
      taskImplementationId: 'parser.k8s.audit-log-specialized',
      taskReferenceId: 'parser.k8s.audit',
      priority: 200,
      isFeature: false,
      isDefaultFeature: false,
      featureLabel: '',
      featureDescription: '',
      dependencies: [],
      selectorRequirements: [],
      compatibleInspectionTypes: [],
      labels: {
        log_type: 'audit',
        env: 'prod',
      },
    } as unknown as RegisteredTaskInfo,
  ];

  const mockGroups: RegisteredTaskGroupInfo[] = [
    {
      taskReferenceId: 'parser.k8s.audit',
      tasks: mockTasks,
    } as unknown as RegisteredTaskGroupInfo,
    {
      taskReferenceId: 'timeline.builder',
      tasks: [
        {
          taskImplementationId: 'timeline.builder.default',
          taskReferenceId: 'timeline.builder',
          priority: 50,
          isFeature: false,
          isDefaultFeature: false,
          featureLabel: '',
          featureDescription: '',
          isInitialTask: false,
          dependencies: [
            {
              cardinality: TaskDependencyCardinality.FAN_IN,
              scope: TaskDependencyScope.ALL,
              targetReferenceId: '',
              targetTag: 'timeline-producer',
            },
          ],
          primaryFeature: '',
          subfeatures: [],
          isGlobal: false,
          labels: {},
        } as unknown as RegisteredTaskInfo,
      ],
    } as unknown as RegisteredTaskGroupInfo,
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Step1RegistryTableComponent, NoopAnimationsModule],
    }).compileComponents();

    fixture = TestBed.createComponent(Step1RegistryTableComponent);
    component = fixture.componentInstance;
  });

  it('renders all task groups initially', () => {
    fixture.componentRef.setInput('taskGroups', mockGroups);
    fixture.detectChanges();

    expect(component.filteredTaskGroups().length).toBe(2);
    expect(component.totalImplementationsCount()).toBe(3);
  });

  it('filters groups by search query matching reference ID', () => {
    fixture.componentRef.setInput('taskGroups', mockGroups);
    fixture.detectChanges();

    component.searchQuery.set('timeline');
    expect(component.filteredTaskGroups().length).toBe(1);
    expect(component.filteredTaskGroups()[0].taskReferenceId).toBe(
      'timeline.builder',
    );
  });

  it('filters groups by search query matching implementation ID', () => {
    fixture.componentRef.setInput('taskGroups', mockGroups);
    fixture.detectChanges();

    component.searchQuery.set('specialized');
    expect(component.filteredTaskGroups().length).toBe(1);
    expect(component.filteredTaskGroups()[0].taskReferenceId).toBe(
      'parser.k8s.audit',
    );
  });

  it('clears search query', () => {
    component.searchQuery.set('test');
    component.clearSearch();
    expect(component.searchQuery()).toBe('');
  });

  it('emits selectImplementation when implementation ID clicked', () => {
    const spy = jasmine.createSpy('selectImplementation');
    component.selectImplementation.subscribe(spy);

    fixture.componentRef.setInput('taskGroups', mockGroups);
    fixture.detectChanges();

    const element: HTMLElement = fixture.nativeElement;
    const implLink = element.querySelector('.impl-id') as HTMLElement;
    implLink.click();

    expect(spy).toHaveBeenCalledWith('parser.k8s.audit-log');
  });
});
