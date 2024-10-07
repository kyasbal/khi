import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DiffComponent } from './diff.component';
import { InMemoryWindowConnectionProvider } from 'src/app/services/frame-connection/window-connection-provider.service';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from 'src/app/services/frame-connection/window-connector.service';
import { DiffModule } from 'src/app/diff/diff.module';
import { DiffPageDataSource } from 'src/app/services/frame-connection/frames/diff-page-datasource.service';

describe('DiffComponent', () => {
  let component: DiffComponent;
  let fixture: ComponentFixture<DiffComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DiffModule],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
        DiffPageDataSource,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(DiffComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
