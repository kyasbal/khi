import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { GraphRenderer } from './graph/renderer';
import { emptyGraphData } from '../../../common/schema/graph-schema';
import { DownloadService } from '../services/donwload-service';
import { GraphPageDataSource } from 'src/app/services/frame-connection/frames/graph-page-datasource.service';
@Component({
  selector: 'graph-architecture-graph',
  templateUrl: './architecture-graph.component.html',
  styleUrls: ['./architecture-graph.component.sass'],
})
export class ArchitectureGraphComponent implements AfterViewInit {
  constructor(
    private dataStore: GraphPageDataSource,
    private downloadService: DownloadService,
  ) {}

  @ViewChild('graphContainer')
  graphContainer!: ElementRef<HTMLDivElement>;

  graphRenderer!: GraphRenderer;

  ngAfterViewInit(): void {
    this.graphRenderer = new GraphRenderer(this.graphContainer.nativeElement);
    this.graphRenderer.updateGraphData(emptyGraphData());
    this.dataStore.data$.subscribe((d) => {
      this.graphRenderer.updateGraphData(d);
    });
    this.downloadService.registerRenderer(this.graphRenderer);
  }
}
