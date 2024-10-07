import { TestBed } from '@angular/core/testing';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { NgxEnvModule } from '@ngx-env/core';
import { GraphComponent } from './graph.component';
import { ArchitectureGraphComponent } from './architecture-graph/architecture-graph.component';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from '../../services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from '../../services/frame-connection/window-connection-provider.service';
import { HeaderModule } from 'src/app/header/header.module';
import { FRONTEND_ANALYTICS } from 'src/app/services/analytics/types';
import { NopFrontendAnalytics } from 'src/app/services/analytics/nop';
import { GraphPageDataSource } from 'src/app/services/frame-connection/frames/graph-page-datasource.service';

describe('GraphComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GraphComponent, ArchitectureGraphComponent],
      imports: [NgxEnvModule, MatToolbarModule, MatIconModule, HeaderModule],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
        {
          provide: FRONTEND_ANALYTICS,
          useClass: NopFrontendAnalytics,
        },
        GraphPageDataSource,
      ],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(GraphComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });
});
