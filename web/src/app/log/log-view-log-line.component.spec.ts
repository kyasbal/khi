import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LogViewLogLineComponent } from './log-view-log-line.component';
import { KHICommonModule } from '../common/common.module';
import { LogType, Severity } from '../generated';
import { KHIFileTextReference } from '../common/schema/khi-file-types';
import { MatTooltipModule } from '@angular/material/tooltip';
import { LogEntry } from '../store/log';

describe('LogViewLogLineComponent', () => {
  let component: LogViewLogLineComponent;
  let fixture: ComponentFixture<LogViewLogLineComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [LogViewLogLineComponent],
      imports: [KHICommonModule, MatTooltipModule],
      providers: [],
    }).compileComponents();
    fixture = TestBed.createComponent(LogViewLogLineComponent);
    component = fixture.componentInstance;
    component.log = new LogEntry(
      0,
      'foo',
      LogType.LogTypeAudit,
      Severity.SeverityInfo,
      0,
      'foo',
      {} as unknown as KHIFileTextReference,
      [],
    );
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
