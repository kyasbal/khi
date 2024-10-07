import { GraphData } from 'src/app/common/schema/graph-schema';
import { InterframeDatasource } from '../inter-frame-datasource.service';
import { Inject, Injectable } from '@angular/core';
import { WindowConnectorService } from '../window-connector.service';
import {
  GRAPH_PAGE_OPEN,
  UPDATE_GRAPH_DATA,
} from 'src/app/common/schema/inter-window-messages';
import {
  FRONTEND_ANALYTICS,
  FrontendAnalytics,
  KHIAnalyticsActivityType,
} from '../../analytics/types';

export interface UpdateGraphMessage {
  graphData: GraphData;
}

@Injectable()
export class GraphPageDataSource extends InterframeDatasource<GraphData> {
  private enabled = false;

  constructor(
    private connector: WindowConnectorService,
    @Inject(FRONTEND_ANALYTICS) private analytics: FrontendAnalytics,
  ) {
    super();
  }
  override enable(): void {
    if (this.enabled) {
      return;
    }
    this.connector
      .receiver<UpdateGraphMessage>(UPDATE_GRAPH_DATA)
      .subscribe((graphData) => {
        this.data$.next(graphData.data.graphData);
        this.analytics.report(KHIAnalyticsActivityType.RenderGraph, {});
      });
    this.connector.broadcast(GRAPH_PAGE_OPEN, {});
  }
  override disable(): void {}
}
