import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { DataLoadSourceExtension } from '../../extensions/data-loader/extension';
import {
  animationFrames,
  BehaviorSubject,
  map,
  Subject,
  takeUntil,
} from 'rxjs';
import { StartupDialogComponent } from 'src/app/dialogs/startup/startup.component';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import {
  POPUP_MANAGER,
  PopupManager,
} from 'src/app/services/popup/popup-manager';
import {
  RequestUserActionPopupComponent,
  RequestUserActionPopupRequest,
} from 'src/app/dialogs/request-user-action-popup/request-user-action-popup.component';
import { NotificationManager } from 'src/app/services/notification/notification';
import { ResizingCalculator } from 'src/app/common/resizable-pane/resizing-calculator';
import { DiffPageDataSourceServer } from 'src/app/services/frame-connection/frames/diff-page-datasource-server.service';
import { GraphPageDataSourceServer } from 'src/app/services/frame-connection/frames/graph-page-datasource-server.service';
import { NilPopupFormRequest } from 'src/app/services/popup/popup-manager-impl';

@Component({
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.sass'],
})
export class AppComponent implements OnInit, OnDestroy {
  readonly destroyed = new Subject<void>();
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
    let lastDialogRef: MatDialogRef<RequestUserActionPopupComponent> | null =
      null;
    this.popupManager
      .requests()
      .pipe(takeUntil(this.destroyed))
      .subscribe((formRequest) => {
        // The last opened dialog will be closed automatically When the popup was cancelled from server side,
        if (formRequest.id === NilPopupFormRequest.id) {
          lastDialogRef?.close();
          lastDialogRef = null;
          return;
        }
        lastDialogRef = this.dialog.open<
          RequestUserActionPopupComponent,
          RequestUserActionPopupRequest
        >(RequestUserActionPopupComponent, {
          data: {
            formRequest,
          },
        });
        this.notificationManager.notify({
          title: 'KHI requests additional parameter',
          body: `Please supply ${formRequest.title} to proceed tasks`,
        });
      });
    animationFrames()
      .pipe(
        takeUntil(this.destroyed),
        map(() => document.body.getBoundingClientRect().width),
      )
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
   * Attempts loading Google Drive file ID from URL hash and load.
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

  ngOnDestroy(): void {
    this.destroyed.next();
  }
}
