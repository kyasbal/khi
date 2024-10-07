import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ArchitectureGraphComponent } from './architecture-graph.component';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from '../../../services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from '../../../services/frame-connection/window-connection-provider.service';
import { FRONTEND_ANALYTICS } from 'src/app/services/analytics/types';
import { NopFrontendAnalytics } from 'src/app/services/analytics/nop';
import { GraphPageDataSource } from 'src/app/services/frame-connection/frames/graph-page-datasource.service';

describe('ArchitectureGraphComponent', () => {
  let component: ArchitectureGraphComponent;
  let fixture: ComponentFixture<ArchitectureGraphComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ArchitectureGraphComponent],
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

    fixture = TestBed.createComponent(ArchitectureGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
