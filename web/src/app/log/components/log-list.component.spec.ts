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
import { LogListComponent } from './log-list.component';
import { LogEntry } from '../../store/log';
import { ResourceTimeline } from '../../store/timeline';
import { ScrollingModule } from '@angular/cdk/scrolling';
import { LogType, Severity } from '../../zzz-generated';
import { ToTextReferenceFromKHIFileBinary } from 'src/app/common/loader/reference-type';
import { ResourceRevision } from '../../store/revision';
import { ResourceEvent } from '../../store/event';

describe('LogListComponent', () => {
  let component: LogListComponent;
  let fixture: ComponentFixture<LogListComponent>;

  const mockLogs: LogEntry[] = [
    new LogEntry(
      0,
      '',
      LogType.LogTypeUnknown,
      Severity.SeverityUnknown,
      1000,
      'sum0',
      ToTextReferenceFromKHIFileBinary({ offset: 0, len: 0, buffer: 0 }),
      [],
    ),
    new LogEntry(
      1,
      '',
      LogType.LogTypeUnknown,
      Severity.SeverityUnknown,
      2000,
      'sum1',
      ToTextReferenceFromKHIFileBinary({ offset: 0, len: 0, buffer: 0 }),
      [],
    ),
    new LogEntry(
      2,
      '',
      LogType.LogTypeUnknown,
      Severity.SeverityUnknown,
      3000,
      'sum2',
      ToTextReferenceFromKHIFileBinary({ offset: 0, len: 0, buffer: 0 }),
      [],
    ),
  ];

  const mockTimelines: ResourceTimeline[] = [
    new ResourceTimeline(
      'mock-timeline',
      'mock-path',
      [{ logIndex: 0 } as unknown as ResourceRevision],
      [
        new ResourceEvent(
          2,
          0,
          LogType.LogTypeUnknown,
          Severity.SeverityUnknown,
        ),
      ],
      0, // Fallback for ParentRelationship enum
    ),
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LogListComponent, ScrollingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(LogListComponent);
    component = fixture.componentInstance;

    // Set required inputs
    fixture.componentRef.setInput('allLogsCount', 3);
    fixture.componentRef.setInput('filteredLogs', mockLogs);
    fixture.componentRef.setInput('selectedLogIndex', -1);
    fixture.componentRef.setInput('highlightLogIndices', new Set<number>());
    fixture.componentRef.setInput('selectedTimelinesWithChildren', []);

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show all logs when filterByTimeline is false', () => {
    fixture.componentRef.setInput('filterByTimeline', false);
    fixture.detectChanges();
    expect(component['shownLogs']()).toEqual(mockLogs);
  });

  it('should filter logs by timeline when filterByTimeline is true', () => {
    fixture.componentRef.setInput('filterByTimeline', true);
    fixture.componentRef.setInput(
      'selectedTimelinesWithChildren',
      mockTimelines,
    );
    fixture.detectChanges();

    const expectedLogs = [mockLogs[0], mockLogs[2]];
    expect(component['shownLogs']()).toEqual(expectedLogs);
  });

  it('should emit logSelected event', () => {
    spyOn(component.logSelected, 'emit');
    component['selectLog'](mockLogs[0]);
    expect(component.logSelected.emit).toHaveBeenCalledWith(mockLogs[0]);
  });

  it('should emit logHovered event', () => {
    spyOn(component.logHovered, 'emit');
    component['onLogHover'](mockLogs[0]);
    expect(component.logHovered.emit).toHaveBeenCalledWith(mockLogs[0]);
  });
});
