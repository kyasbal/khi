import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LogHeaderComponent } from './header.component';
import { WINDOW_CONNECTION_PROVIDER } from '../services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from '../services/frame-connection/window-connection-provider.service';
import { LOG_ANNOTATOR_RESOLVER } from '../annotator/log/resolver';
import { getDefaultLogAnnotatorResolver } from '../annotator/log/default';

describe('LogHeaderComponent', () => {
  let component: LogHeaderComponent;
  let fixture: ComponentFixture<LogHeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [LogHeaderComponent],
      imports: [],
      providers: [
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
        {
          provide: LOG_ANNOTATOR_RESOLVER,
          useValue: getDefaultLogAnnotatorResolver(),
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(LogHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
