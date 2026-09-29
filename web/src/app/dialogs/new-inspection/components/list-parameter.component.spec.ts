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

import { NgModule, provideZoneChangeDetection } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatIconRegistry } from '@angular/material/icon';
import { By } from '@angular/platform-browser';
import {
  BrowserTestingModule,
  platformBrowserTesting,
} from '@angular/platform-browser/testing';
import {
  GroupParameterFormField,
  ListParameterFormField,
  ParameterFormValidationTiming,
  ParameterHintType,
  ParameterInputType,
  UploadStatus,
} from 'src/app/common/schema/form-types';
import {
  FILE_UPLOADER,
  MockFileUploader,
} from 'src/app/dialogs/new-inspection/components/service/file-uploader';
import {
  DefaultParameterStore,
  PARAMETER_STORE,
} from 'src/app/dialogs/new-inspection/components/service/parameter-store';
import { ListParameterComponent } from 'src/app/dialogs/new-inspection/components/list-parameter.component';

@NgModule({ providers: [provideZoneChangeDetection()] })
export class ZoneChangeDetectionModule {}

describe('ListParameterComponent', () => {
  let fixture: ComponentFixture<ListParameterComponent>;
  let parameterStore: DefaultParameterStore;

  beforeAll(() => {
    TestBed.resetTestEnvironment();
    TestBed.initTestEnvironment(
      [ZoneChangeDetectionModule, BrowserTestingModule],
      platformBrowserTesting(),
      { teardown: { destroyAfterEach: false } },
    );
  });

  beforeEach(async () => {
    parameterStore = new DefaultParameterStore();
    await TestBed.configureTestingModule({
      imports: [BrowserAnimationsModule],
      providers: [
        {
          provide: FILE_UPLOADER,
          useValue: new MockFileUploader(),
        },
        {
          provide: PARAMETER_STORE,
          useValue: parameterStore,
        },
      ],
    }).compileComponents();
    const matIconRegistry = TestBed.inject(MatIconRegistry);
    matIconRegistry.setDefaultFontSetClass('material-symbols-outlined');
    fixture = TestBed.createComponent(ListParameterComponent);
  });

  afterEach(() => {
    parameterStore.destroy();
  });

  it('renders list items with single file items', () => {
    const mockList: ListParameterFormField = {
      id: 'file-list',
      type: ParameterInputType.List,
      label: 'File List',
      description: 'Upload files',
      hint: '',
      hintType: ParameterHintType.None,
      default: ['0'],
      minCount: 0,
      maxCount: 5,
      addButtonLabel: 'Add File',
      items: [
        {
          key: '0',
          field: {
            id: 'file-0',
            type: ParameterInputType.File,
            label: 'File 0',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            token: { id: 'tok-0' },
            status: UploadStatus.Waiting,
          },
        },
      ],
    };

    fixture.componentRef.setInput('parameter', mockList);
    fixture.detectChanges();

    const cards = fixture.debugElement.queryAll(By.css('.list-item-card'));
    expect(cards.length).toBe(1);
    const fileComp = fixture.debugElement.query(
      By.css('khi-new-inspection-file-parameter'),
    );
    expect(fileComp).toBeTruthy();
  });

  it('renders list items with nested group items', () => {
    const mockGroup: GroupParameterFormField = {
      id: 'grp-0',
      type: ParameterInputType.Group,
      label: 'Group 0',
      description: '',
      hint: '',
      hintType: ParameterHintType.None,
      collapsible: false,
      collapsedByDefault: false,
      children: [
        {
          id: 'text-in-grp',
          type: ParameterInputType.Text,
          label: 'Text in Group',
          description: '',
          hint: '',
          hintType: ParameterHintType.None,
          default: '',
          readonly: false,
          suggestions: [],
          validationTiming: ParameterFormValidationTiming.Change,
        },
      ],
    };

    const mockList: ListParameterFormField = {
      id: 'group-list',
      type: ParameterInputType.List,
      label: 'Group List',
      description: '',
      hint: '',
      hintType: ParameterHintType.None,
      default: ['item-0'],
      minCount: 0,
      maxCount: 0,
      addButtonLabel: 'Add Group',
      items: [
        {
          key: 'item-0',
          field: mockGroup,
        },
      ],
    };

    fixture.componentRef.setInput('parameter', mockList);
    fixture.detectChanges();

    const groupComp = fixture.debugElement.query(
      By.css('khi-new-inspection-group-parameter'),
    );
    expect(groupComp).toBeTruthy();
  });

  it('clicking add button appends a new monotonic key to PARAMETER_STORE', () => {
    const mockList: ListParameterFormField = {
      id: 'list-add',
      type: ParameterInputType.List,
      label: 'List',
      description: '',
      hint: '',
      hintType: ParameterHintType.None,
      default: ['0'],
      minCount: 0,
      maxCount: 5,
      addButtonLabel: 'Add Item',
      items: [
        {
          key: '0',
          field: {
            id: 'item-0',
            type: ParameterInputType.Text,
            label: 'Item 0',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
      ],
    };

    fixture.componentRef.setInput('parameter', mockList);
    fixture.detectChanges();

    const addBtn = fixture.debugElement.query(By.css('.add-item-btn'));
    addBtn.nativeElement.click();
    fixture.detectChanges();

    expect(parameterStore.get<string[]>('list-add')()).toEqual(['0', '1']);
  });

  it('clicking remove button removes that specific item key while preserving others', () => {
    const mockList: ListParameterFormField = {
      id: 'list-rem',
      type: ParameterInputType.List,
      label: 'List',
      description: '',
      hint: '',
      hintType: ParameterHintType.None,
      default: ['0', '1', '2'],
      minCount: 1,
      maxCount: 5,
      addButtonLabel: 'Add Item',
      items: [
        {
          key: '0',
          field: {
            id: 'item-0',
            type: ParameterInputType.Text,
            label: 'Item 0',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
        {
          key: '1',
          field: {
            id: 'item-1',
            type: ParameterInputType.Text,
            label: 'Item 1',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
        {
          key: '2',
          field: {
            id: 'item-2',
            type: ParameterInputType.Text,
            label: 'Item 2',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
      ],
    };

    parameterStore.set('list-rem', ['0', '1', '2']);
    fixture.componentRef.setInput('parameter', mockList);
    fixture.detectChanges();

    const removeButtons = fixture.debugElement.queryAll(
      By.css('.remove-item-btn'),
    );
    expect(removeButtons.length).toBe(3);
    removeButtons[1].nativeElement.click();
    fixture.detectChanges();

    const remainingCards = fixture.debugElement.queryAll(
      By.css('.list-item-card'),
    );
    expect(remainingCards.length).toBe(2);

    expect(parameterStore.get<string[]>('list-rem')()).toEqual(['0', '2']);
  });

  it('deleting a middle item and then adding a new item never reuses a previously used key', () => {
    const mockList: ListParameterFormField = {
      id: 'list-reuse',
      type: ParameterInputType.List,
      label: 'List',
      description: '',
      hint: '',
      hintType: ParameterHintType.None,
      default: ['0', '1'],
      minCount: 0,
      maxCount: 5,
      addButtonLabel: 'Add Item',
      items: [
        {
          key: '0',
          field: {
            id: 'item-0',
            type: ParameterInputType.Text,
            label: 'Item 0',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
        {
          key: '1',
          field: {
            id: 'item-1',
            type: ParameterInputType.Text,
            label: 'Item 1',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
      ],
    };

    parameterStore.set('list-reuse', ['0', '1']);
    fixture.componentRef.setInput('parameter', mockList);
    fixture.detectChanges();

    // Remove item '1'
    fixture.componentInstance.removeItem('1');
    fixture.detectChanges();
    expect(parameterStore.get<string[]>('list-reuse')()).toEqual(['0']);

    // Add item: max existing was 1, next key should be at least '2'
    fixture.componentInstance.addItem();
    fixture.detectChanges();
    expect(parameterStore.get<string[]>('list-reuse')()).toEqual(['0', '2']);
  });

  it('disables add button when maxCount > 0 and activeKeys length >= maxCount', () => {
    const mockList: ListParameterFormField = {
      id: 'list-max',
      type: ParameterInputType.List,
      label: 'List',
      description: '',
      hint: '',
      hintType: ParameterHintType.None,
      default: ['0', '1'],
      minCount: 0,
      maxCount: 2,
      addButtonLabel: 'Add',
      items: [
        {
          key: '0',
          field: {
            id: 'item-0',
            type: ParameterInputType.Text,
            label: '0',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
        {
          key: '1',
          field: {
            id: 'item-1',
            type: ParameterInputType.Text,
            label: '1',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
      ],
    };

    parameterStore.set('list-max', ['0', '1']);
    fixture.componentRef.setInput('parameter', mockList);
    fixture.detectChanges();

    const addBtn = fixture.debugElement.query(By.css('.add-item-btn'));
    expect(addBtn.nativeElement.disabled).toBeTrue();
  });

  it('disables remove button when activeKeys length <= minCount', () => {
    const mockList: ListParameterFormField = {
      id: 'list-min',
      type: ParameterInputType.List,
      label: 'List',
      description: '',
      hint: '',
      hintType: ParameterHintType.None,
      default: ['0'],
      minCount: 1,
      maxCount: 5,
      addButtonLabel: 'Add',
      items: [
        {
          key: '0',
          field: {
            id: 'item-0',
            type: ParameterInputType.Text,
            label: '0',
            description: '',
            hint: '',
            hintType: ParameterHintType.None,
            default: '',
            readonly: false,
            suggestions: [],
            validationTiming: ParameterFormValidationTiming.Change,
          },
        },
      ],
    };

    parameterStore.set('list-min', ['0']);
    fixture.componentRef.setInput('parameter', mockList);
    fixture.detectChanges();

    const removeBtn = fixture.debugElement.query(By.css('.remove-item-btn'));
    expect(removeBtn.nativeElement.disabled).toBeTrue();
  });
});
