import { Component, inject, OnInit } from '@angular/core';
import { DataLoadSourceExtension } from '../../extensions/data-loader/extension';
import { animationFrames, BehaviorSubject, map } from 'rxjs';
import { StartupDialogComponent } from 'src/app/dialogs/startup/startup.component';
import { MatDialog } from '@angular/material/dialog';
import {
  POPUP_MANAGER,
  PopupManager,
} from 'src/app/services/popup/popup-manager';
import {
  AdditionalInputPopupComponent,
  AdditionalInputPopupDialogRequest,
} from 'src/app/dialogs/additional-input-popup/additional-input-popup.component';
import { NotificationManager } from 'src/app/services/notification/notification';
import { ResizingCalculator } from 'src/app/common/resizable-pane/resizing-calculator';
import { DiffPageDataSourceServer } from 'src/app/services/frame-connection/frames/diff-page-datasource-server.service';
import { GraphPageDataSourceServer } from 'src/app/services/frame-connection/frames/graph-page-datasource-server.service';

@Component({
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.sass'],
})
export class AppComponent implements OnInit {
  readonly showLogPane = new BehaviorSubject<boolean>(true);
  readonly showHistoryPane = new BehaviorSubject<boolean>(true);
  readonly popupManager: PopupManager = inject(POPUP_MANAGER);
  readonly diffPageSourceSender: DiffPageDataSourceServer = inject(
    DiffPageDataSourceServer,
  );
  readonly graphPageSourceSender: GraphPageDataSourceServer = inject(
    GraphPageDataSourceServer,
  );
  readonly notificationManager: NotificationManager =
    inject(NotificationManager);
  readonly resizer = new ResizingCalculator([
    {
      id: 'explorer-view',
      initialSize: 300,
      minSizeInPx: 300,
      resizeRatio: 0,
    },
    {
      id: 'explorer-view-expander',
      initialSize: 5,
      minSizeInPx: 5,
      resizeRatio: 0,
    },
    {
      id: 'timeline-view',
      initialSize: 300,
      minSizeInPx: 300,
      resizeRatio: 1,
    },
    {
      id: 'log-view-expander',
      initialSize: 5,
      resizeRatio: 0,
      minSizeInPx: 5,
    },
    {
      id: 'log-view',
      initialSize: 300,
      minSizeInPx: 200,
      resizeRatio: 0,
    },
    {
      id: 'history-view-expander',
      initialSize: 5,
      minSizeInPx: 5,
      resizeRatio: 0,
    },
    {
      id: 'history-view',
      initialSize: 300,
      minSizeInPx: 200,
      resizeRatio: 0,
    },
  ]);

  constructor(
    private loaderExtension: DataLoadSourceExtension,
    private dialog: MatDialog,
  ) {}

  ngOnInit() {
    if (!this._processFileLoaderExtension()) {
      this.dialog.open(StartupDialogComponent, {
        maxWidth: '100vw',
        panelClass: 'startup-modalbox',
        disableClose: true,
      });
    }
    // Start monitoring popup request from server
    this.popupManager.requests().subscribe((formRequest) => {
      this.notificationManager.notify({
        title: 'KHI requests additional parameter',
        body: `Please supply ${formRequest.title} to proceed tasks`,
      });
      this.dialog.open<
        AdditionalInputPopupComponent,
        AdditionalInputPopupDialogRequest
      >(AdditionalInputPopupComponent, {
        data: {
          formRequest,
        },
      });
    });
    animationFrames()
      .pipe(map(() => document.body.getBoundingClientRect().width))
      .subscribe((width) => {
        this.resizer.setContainerSizeInPx(width);
      });
    this.diffPageSourceSender.activate();
    this.graphPageSourceSender.activate();
  }

  togglePane(pane: 'log' | 'history') {
    switch (pane) {
      case 'log':
        this.showLogPane.next(!this.showLogPane.value);
        break;
      case 'history':
        this.showHistoryPane.next(!this.showHistoryPane.value);
        break;
    }
  }

  /**
   *
   * @returns if file loader extension loads the file or not
   */
  private _processFileLoaderExtension(): boolean {
    if (window.location.hash.length > 1) {
      const hash = window.location.hash.substring(1);
      this.loaderExtension.load(hash);
      return true;
    }
    return false;
  }
}
