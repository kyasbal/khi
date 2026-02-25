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

import { Meta, StoryObj } from '@storybook/angular';
import { TypeSeverityComponent } from './type-severity.component';

export default {
  title: 'log/TypeSeverityComponent',
  component: TypeSeverityComponent,
} as Meta<TypeSeverityComponent>;

type Story = StoryObj<TypeSeverityComponent>;

export const Info: Story = {
  args: {
    logType: 'k8s_audit',
    severity: 'INFO',
  },
};

export const Warning: Story = {
  args: {
    logType: 'k8s_audit',
    severity: 'WARNING',
  },
};

export const ErrorSeverity: Story = {
  args: {
    logType: 'k8s_node',
    severity: 'ERROR',
  },
};

export const Unknown: Story = {
  args: {
    logType: 'k8s_container',
    severity: 'UNKNOWN',
  },
};
