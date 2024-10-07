import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TitleBarComponent } from './titlebar.component';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from '../services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from '../services/frame-connection/window-connection-provider.service';
import { NgxEnvModule } from '@ngx-env/core';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { KHICommonModule } from '../common/common.module';
import { MatMenuModule } from '@angular/material/menu';

describe('TitlebarComponent', () => {
  let component: TitleBarComponent;
  let fixture: ComponentFixture<TitleBarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TitleBarComponent],
      imports: [
        NgxEnvModule,
        MatIconModule,
        MatButtonModule,
        KHICommonModule,
        MatMenuModule,
      ],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(TitleBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
