import { Component, Input } from '@angular/core';
import { WindowConnectorService } from '../services/frame-connection/window-connector.service';
import { map } from 'rxjs';

@Component({
  selector: 'khi-title',
  templateUrl: './titlebar.component.html',
  styleUrls: ['./titlebar.component.sass'],
})
export class TitleBarComponent {
  @Input()
  pageName = 'N/A';

  mainPageConenctionEstablished =
    this.windowConnector.mainPageConenctionEstablished;

  sessionId = this.windowConnector.sessionEstablished.pipe(
    map(() => this.windowConnector.sessionId),
  );

  sessionPages = this.windowConnector.sessionPages;

  constructor(private readonly windowConnector: WindowConnectorService) {}

  focusWindow(frameId: string) {
    this.windowConnector.focusWindow(frameId);
  }
}
