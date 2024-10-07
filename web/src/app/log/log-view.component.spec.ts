import { OverlayModule } from '@angular/cdk/overlay';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FormsModule } from '@angular/forms';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatToolbarModule } from '@angular/material/toolbar';
import { KHICommonModule } from '../common/common.module';
import { LogBodyComponent } from './body.component';

import { LogViewComponent } from './log-view.component';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from '../services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from '../services/frame-connection/window-connection-provider.service';
import { IconToggleButtonComponent } from './icon-toggle-button.component';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { LOG_TOOL_ANNOTATOR_RESOLVER } from '../annotator/log-tool/resolver';
import { getDefaultLogToolAnnotatorResolver } from '../annotator/log-tool/default';

describe('LogViewComponent', () => {
  let component: LogViewComponent;
  let fixture: ComponentFixture<LogViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [
        LogViewComponent,
        LogBodyComponent,
        IconToggleButtonComponent,
      ],
      imports: [
        OverlayModule,
        MatToolbarModule,
        MatSlideToggleModule,
        MatIconModule,
        MatTooltipModule,
        FormsModule,
        KHICommonModule,
      ],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
        {
          provide: LOG_TOOL_ANNOTATOR_RESOLVER,
          useValue: getDefaultLogToolAnnotatorResolver(),
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(LogViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
