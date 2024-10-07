import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MainMenuComponent } from './main-menu.component';
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
import { NgxEnvModule } from '@ngx-env/core';
import { provideHttpClient } from '@angular/common/http';

describe('MainMenuComponent', () => {
  let component: MainMenuComponent;
  let fixture: ComponentFixture<MainMenuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [MainMenuComponent],
      imports: [MatDialogModule, MatIconModule, MatMenuModule, NgxEnvModule],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
        InspectionDataLoaderService,
        InspectionDataStoreService,
        provideHttpClient(),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(MainMenuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
