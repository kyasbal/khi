import { TestBed } from '@angular/core/testing';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { AppComponent } from './main.component';
import { DiffModule } from '../../diff/diff.module';
import { DataLoadSourceExtension } from '../../extensions/data-loader/extension';
import { HeaderModule } from '../../header/header.module';
import { LogModule } from '../../log/log.module';
import { InspectionDataLoaderService } from '../../services/data-loader.service';
import { TimelineModule } from '../../timeline/timeline.module';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from 'src/app/services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from 'src/app/services/frame-connection/window-connection-provider.service';
import { MatIconModule } from '@angular/material/icon';
import { provideHttpClient } from '@angular/common/http';
import { POPUP_MANAGER } from 'src/app/services/popup/popup-manager';
import { MockPopupManager } from 'src/app/services/popup/mock';
import { DiffPageDataSourceServer } from 'src/app/services/frame-connection/frames/diff-page-datasource-server.service';
import { GraphPageDataSourceServer } from 'src/app/services/frame-connection/frames/graph-page-datasource-server.service';

describe('AppComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AppComponent],
      imports: [
        MatDialogModule,
        HeaderModule,
        TimelineModule,
        LogModule,
        DiffModule,
        MatIconModule,
        MatSnackBarModule,
      ],
      providers: [
        DataLoadSourceExtension,
        InspectionDataLoaderService,
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
        {
          provide: POPUP_MANAGER,
          useValue: new MockPopupManager(),
        },
        provideHttpClient(),
        DiffPageDataSourceServer,
        GraphPageDataSourceServer,
      ],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });
});
