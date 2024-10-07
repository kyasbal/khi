import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MatDialogModule } from '@angular/material/dialog';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from '../services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from '../services/frame-connection/window-connection-provider.service';
import { InspectionDataLoaderService } from '../services/data-loader.service';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { GraphMenuComponent } from './graph-menu.component';

describe('GraphMenuComponent', () => {
  let component: GraphMenuComponent;
  let fixture: ComponentFixture<GraphMenuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GraphMenuComponent],
      imports: [MatDialogModule, MatIconModule, MatMenuModule],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
        InspectionDataLoaderService,
        InspectionDataStoreService,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(GraphMenuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
