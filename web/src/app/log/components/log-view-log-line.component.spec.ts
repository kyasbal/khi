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

import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LogViewLogLineComponent } from './log-view-log-line.component';
import { LogEntry } from 'src/app/store/log';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TimestampFormatPipe } from 'src/app/common/timestamp-format.pipe';
import { By } from '@angular/platform-browser';
import { LogType, Severity } from 'src/app/zzz-generated';
import { ToTextReferenceFromKHIFileBinary } from 'src/app/common/loader/reference-type';

describe('LogViewLogLineComponent', () => {
  let component: LogViewLogLineComponent;
  let fixture: ComponentFixture<LogViewLogLineComponent>;

  const mockLog = new LogEntry(
    1,
    'mock-insert-id',
    LogType.LogTypeUnknown,
    Severity.SeverityUnknown,
    1700000000000,
    'Test summary',
    ToTextReferenceFromKHIFileBinary({ offset: 0, len: 0, buffer: 0 }),
    [],
  );

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LogViewLogLineComponent, MatTooltipModule, TimestampFormatPipe],
    }).compileComponents();

    fixture = TestBed.createComponent(LogViewLogLineComponent);
    component = fixture.componentInstance;

    // Set required input
    fixture.componentRef.setInput('log', mockLog);

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display log message', () => {
    const messageEl = fixture.debugElement.query(By.css('.message'));
    expect(messageEl.nativeElement.textContent).toContain('Test summary'); // The template uses logEntry.summary in the message div based on the diff we saw earlier
  });

  it('should emit lineClick event on click', () => {
    spyOn(component.lineClick, 'emit');
    const rowEl = fixture.debugElement.query(By.css('.log-row'));
    rowEl.nativeElement.click();
    expect(component.lineClick.emit).toHaveBeenCalledWith(mockLog);
  });

  it('should emit lineClick event on Enter key', () => {
    spyOn(component.lineClick, 'emit');
    const rowEl = fixture.debugElement.query(By.css('.log-row'));
    const event = new KeyboardEvent('keydown', { key: 'Enter' });
    rowEl.nativeElement.dispatchEvent(event);
    expect(component.lineClick.emit).toHaveBeenCalledWith(mockLog);
  });

  it('should emit lineClick event on Space key', () => {
    spyOn(component.lineClick, 'emit');
    const rowEl = fixture.debugElement.query(By.css('.log-row'));
    const event = new KeyboardEvent('keydown', { key: ' ' });
    rowEl.nativeElement.dispatchEvent(event);
    expect(component.lineClick.emit).toHaveBeenCalledWith(mockLog);
  });

  it('should emit lineHover event on mouseover', () => {
    spyOn(component.lineHover, 'emit');
    const rowEl = fixture.debugElement.query(By.css('.log-row'));
    rowEl.nativeElement.dispatchEvent(new Event('mouseover'));
    expect(component.lineHover.emit).toHaveBeenCalledWith(mockLog);
  });

  it('should apply highlight class when highlighted input is true', () => {
    fixture.componentRef.setInput('highlighted', true);
    fixture.detectChanges();
    const rowEl = fixture.debugElement.query(By.css('.log-row'));
    expect(rowEl.nativeElement.classList.contains('highlight')).toBeTrue();
  });

  it('should apply selected class when selected input is true', () => {
    fixture.componentRef.setInput('selected', true);
    fixture.detectChanges();
    const rowEl = fixture.debugElement.query(By.css('.log-row'));
    expect(rowEl.nativeElement.classList.contains('selected')).toBeTrue();
  });
});
