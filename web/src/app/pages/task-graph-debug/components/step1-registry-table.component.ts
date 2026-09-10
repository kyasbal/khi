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
  signal,
} from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import {
  RegisteredTaskGroupInfo,
  RegisteredTaskInfo,
  TaskDependencyCardinality,
  TaskDependencyScope,
} from 'src/app/generated/api/v1/inspection_task_graph_pb';
import { TaskLabelEntry } from 'src/app/pages/task-graph-debug/types/task-graph-debug.model';
import { KHIIconRegistrationModule } from 'src/app/shared/module/icon-registration.module';

/**
 * Renders Step 1: All registered tasks grouped by TaskReferenceID.
 */
@Component({
  selector: 'khi-step1-registry-table',
  templateUrl: './step1-registry-table.component.html',
  styleUrls: ['./step1-registry-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatIconModule, MatButtonModule, KHIIconRegistrationModule],
})
export class Step1RegistryTableComponent {
  /**
   * Export TaskDependencyCardinality enum to the template.
   */
  protected readonly TaskDependencyCardinality = TaskDependencyCardinality;

  /**
   * All registered task groups by reference identifier.
   */
  readonly taskGroups = input<readonly RegisteredTaskGroupInfo[]>([]);

  /**
   * Emitted when user clicks on a task implementation to inspect.
   */
  readonly selectImplementation = output<string>();

  /**
   * Text search query filter.
   */
  readonly searchQuery = signal<string>('');

  /**
   * Task groups filtered by search query matching reference ID or implementation ID.
   */
  readonly filteredTaskGroups = computed<readonly RegisteredTaskGroupInfo[]>(
    () => {
      const q = this.searchQuery().trim().toLowerCase();
      const groups = this.taskGroups();
      if (!q) {
        return groups;
      }
      return groups.filter((g) => {
        if (g.taskReferenceId.toLowerCase().includes(q)) {
          return true;
        }
        return g.tasks.some((task) =>
          task.taskImplementationId.toLowerCase().includes(q),
        );
      });
    },
  );

  /**
   * Total number of distinct task implementations registered.
   */
  readonly totalImplementationsCount = computed<number>(() =>
    this.taskGroups().reduce((acc, g) => acc + g.tasks.length, 0),
  );

  /**
   * Updates the search filter query.
   */
  onSearchInput(event: Event): void {
    const inputElement = event.target as HTMLInputElement;
    this.searchQuery.set(inputElement.value);
  }

  /**
   * Clears the search filter query.
   */
  clearSearch(): void {
    this.searchQuery.set('');
  }

  /**
   * Converts task metadata labels into an iterable array of key-value tuples.
   */
  getLabelEntries(task: RegisteredTaskInfo): readonly TaskLabelEntry[] {
    return Object.entries(task.labels).map(([key, value]) => ({ key, value }));
  }

  /**
   * Returns human-readable label for dependency scope.
   */
  getScopeLabel(scope: TaskDependencyScope): string {
    switch (scope) {
      case TaskDependencyScope.ACTIVE_FEATURES:
        return 'Active features';
      case TaskDependencyScope.ACTIVE_GRAPH:
        return 'Active graph';
      case TaskDependencyScope.ALL:
        return 'All candidates';
      default:
        return 'Unspecified';
    }
  }
}
