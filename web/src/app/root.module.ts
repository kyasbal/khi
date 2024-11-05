import { NgModule, importProvidersFrom } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './pages/main/main.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { provideHighlightOptions } from 'ngx-highlightjs';
import { RegisterGoogleDriveExtensionProvidersIfEnabled } from './extensions/data-loader/google-drive';
import { DataLoadSourceExtension } from './extensions/data-loader/extension';
import { InspectionDataLoaderService } from './services/data-loader.service';
import { TimelineSelectionService } from './services/timeline-selection.service';
import { InspectionDataStoreService } from './services/inspection-data-store.service';
import { SelectionManagerService } from './services/selection-manager.service';
import { HeaderModule } from './header/header.module';
import { DialogsModule } from './dialogs/dialogs.module';
import { GoogleDriveDataLoaderModule } from './extensions/data-loader/module';
import { LogModule } from './log/log.module';
import { DiffModule } from './diff/diff.module';
import { CommonModule } from '@angular/common';
import { TimelineModule } from './timeline/timeline.module';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { RouterModule, TitleStrategy } from '@angular/router';
import { KHIRoutes } from './app.route';
import { RootComponent } from './root.component';
import { GraphPageModule } from './pages/graph/graph.module';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectorService,
} from './services/frame-connection/window-connector.service';
import { BroadcastChannelWindowConnectionProvider } from './services/frame-connection/window-connection-provider.service';
import { KHITitleStrategy } from './services/title-strategy.service';
import { KHICommonModule } from './common/common.module';
import { MatIconModule, MatIconRegistry } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { HttpClientModule } from '@angular/common/http';
import { RequestUserActionPopupComponent } from './dialogs/request-user-action-popup/request-user-action-popup.component';
import { POPUP_MANAGER } from './services/popup/popup-manager';
import { PopupManagerImpl } from './services/popup/popup-manager-impl';
import { BACKEND_API } from './services/api/backend-api-interface';
import { BackendAPIImpl } from './services/api/backend-api.service';
import { NotificationManager } from './services/notification/notification';
import { FRONTEND_ANALYTICS } from './services/analytics/types';
import { FrontendAnalyticsWithGA } from './services/analytics/ga';
import { ProgressDialogService } from './services/progress/progress-dialog.service';
import {
  BACKEND_CONNECTION,
  BackendConnectionServiceImpl,
} from './services/api/backend-connection.service';
import { DiffPageDataSource } from './services/frame-connection/frames/diff-page-datasource.service';
import { DiffPageDataSourceServer } from './services/frame-connection/frames/diff-page-datasource-server.service';
import { GraphPageDataSourceServer } from './services/frame-connection/frames/graph-page-datasource-server.service';
@NgModule({
  declarations: [AppComponent, RootComponent],
  imports: [
    CommonModule,
    KHICommonModule,
    BrowserModule,
    BrowserAnimationsModule,
    HeaderModule,
    DialogsModule,
    LogModule,
    DiffModule,
    GraphPageModule,
    TimelineModule,
    MatSnackBarModule,
    RouterModule.forRoot(KHIRoutes),
    MatIconModule,
    MatButtonModule,

    // Extension modules
    GoogleDriveDataLoaderModule,

    // Standoalone components
    RequestUserActionPopupComponent,
  ],
  providers: [
    importProvidersFrom(HttpClientModule),
    provideHighlightOptions({
      coreLibraryLoader: () => import('highlight.js/lib/core'),
      lineNumbersLoader: () => import('ngx-highlightjs/line-numbers'),
      languages: {
        yaml: () => import('highlight.js/lib/languages/yaml'),
      },
    }),
    { provide: TitleStrategy, useClass: KHITitleStrategy },
    { provide: FRONTEND_ANALYTICS, useClass: FrontendAnalyticsWithGA },
    ...RegisterGoogleDriveExtensionProvidersIfEnabled(),
    ...ProgressDialogService.providers(),
    DataLoadSourceExtension,
    InspectionDataLoaderService,
    DiffPageDataSourceServer,
    GraphPageDataSourceServer,

    TimelineSelectionService,
    InspectionDataStoreService,
    SelectionManagerService,
    WindowConnectorService,
    {
      provide: WINDOW_CONNECTION_PROVIDER,
      useValue: new BroadcastChannelWindowConnectionProvider(),
    },
    {
      provide: BACKEND_API,
      useClass: BackendAPIImpl,
    },
    {
      provide: BACKEND_CONNECTION,
      useClass: BackendConnectionServiceImpl,
    },
    {
      provide: POPUP_MANAGER,
      useClass: PopupManagerImpl,
    },
    NotificationManager,
    DiffPageDataSource,
  ],
  bootstrap: [RootComponent],
})
export class RootModule {
  constructor(
    iconRegistry: MatIconRegistry,
    notificationManager: NotificationManager,
  ) {
    iconRegistry.setDefaultFontSetClass('material-symbols-outlined');
    notificationManager.initialize();
  }
}
