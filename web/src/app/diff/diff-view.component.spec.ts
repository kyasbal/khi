import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';

import { DiffViewComponent } from './diff-view.component';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from '../services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from '../services/frame-connection/window-connection-provider.service';
import { DiffModule } from './diff.module';

describe('DiffViewComponent', () => {
  let component: DiffViewComponent;
  let fixture: ComponentFixture<DiffViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [DiffViewComponent],
      imports: [MatButtonModule, MatToolbarModule, MatIconModule, DiffModule],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(DiffViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
