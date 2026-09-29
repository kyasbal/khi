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

import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import {
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
import { flattenDefaultValues } from 'src/app/dialogs/new-inspection/utils/new-inspection.utils';

const mockFileListParameter: ListParameterFormField = {
  id: 'log-file-list',
  label: 'Log Files',
  description: 'Upload multiple audit or node log files for inspection.',
  type: ParameterInputType.List,
  hint: '',
  hintType: ParameterHintType.None,
  minCount: 1,
  maxCount: 5,
  addButtonLabel: 'Add Log File',
  default: ['0'],
  items: [
    {
      key: '0',
      field: {
        id: 'file-entry-0',
        label: 'Log File',
        description: 'Choose a log bundle.',
        type: ParameterInputType.File,
        hint: '',
        hintType: ParameterHintType.None,
        token: { id: 'sample-token' },
        status: UploadStatus.Waiting,
      },
    },
  ],
};

const mockNodeGroupListParameter: ListParameterFormField = {
  id: 'node-groups',
  label: 'Node Groups',
  description: 'Configure repeatable node pool inspection criteria.',
  type: ParameterInputType.List,
  hint: '',
  hintType: ParameterHintType.None,
  minCount: 0,
  maxCount: 0,
  addButtonLabel: 'Add Node Group',
  default: ['pool-0'],
  items: [
    {
      key: 'pool-0',
      field: {
        id: 'group-entry-0',
        label: 'Pool Configuration',
        description: 'Pool parameters',
        type: ParameterInputType.Group,
        hint: '',
        hintType: ParameterHintType.None,
        collapsible: false,
        collapsedByDefault: false,
        children: [
          {
            id: 'pool-name-0',
            label: 'Pool Name',
            description: 'Name of the node pool.',
            type: ParameterInputType.Text,
            hint: '',
            hintType: ParameterHintType.None,
            default: 'default-pool',
            readonly: false,
            suggestions: ['default-pool', 'highmem-pool'],
            validationTiming: ParameterFormValidationTiming.Change,
          },
          {
            id: 'include-system-pods-0',
            label: 'Include System Pods',
            description: 'Inspect kube-system daemonsets on this pool.',
            type: ParameterInputType.Checkbox,
            hint: '',
            hintType: ParameterHintType.None,
            default: true,
            readonly: false,
          },
        ],
      },
    },
  ],
};

const createInitializedParameterStore = (param: ListParameterFormField) => {
  const store = new DefaultParameterStore();
  const defaults = flattenDefaultValues([param]);
  store.setDefaultValues(defaults);
  store.setValidatedParameters(defaults);
  return store;
};

const meta: Meta<ListParameterComponent> = {
  title: 'Dialogs/NewInspection/ListParameter',
  component: ListParameterComponent,
  tags: ['autodocs'],
  decorators: [
    moduleMetadata({
      imports: [BrowserAnimationsModule],
      providers: [
        {
          provide: PARAMETER_STORE,
          useFactory: () =>
            createInitializedParameterStore(mockFileListParameter),
        },
        {
          provide: FILE_UPLOADER,
          useClass: MockFileUploader,
        },
      ],
    }),
  ],
};

export default meta;
type Story = StoryObj<ListParameterComponent>;

export const RepeatableFileList: Story = {
  args: {
    parameter: mockFileListParameter,
  },
};

export const RepeatableNodeGroupList: Story = {
  args: {
    parameter: mockNodeGroupListParameter,
  },
};

export const MinMaxConstrainedList: Story = {
  args: {
    parameter: {
      ...mockFileListParameter,
      id: 'constrained-file-list',
      minCount: 1,
      maxCount: 2,
    },
  },
};
