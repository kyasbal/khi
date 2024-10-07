import { Component } from '@angular/core';
import { DownloadService } from '../pages/graph/services/donwload-service';

@Component({
  selector: 'khi-graph-menu',
  templateUrl: './graph-menu.component.html',
  styleUrls: ['./graph-menu.component.sass'],
})
export class GraphMenuComponent {
  constructor(private downloadService: DownloadService) {}

  downloadAsPng() {
    this.downloadService.downloadAsPng();
  }

  downloadAsSvg() {
    this.downloadService.downloadAsSvg();
  }
}
