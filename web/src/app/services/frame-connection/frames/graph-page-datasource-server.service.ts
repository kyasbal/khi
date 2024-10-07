import {
  GRAPH_PAGE_OPEN,
  UPDATE_GRAPH_DATA,
} from 'src/app/common/schema/inter-window-messages';
import { GraphDataConverterService } from '../../graph-converter.service';
import { SelectionManagerService } from '../../selection-manager.service';
import { WindowConnectorService } from '../window-connector.service';
import { withLatestFrom } from 'rxjs';
import { InspectionDataStoreService } from '../../inspection-data-store.service';
import { Injectable } from '@angular/core';
import { UpdateGraphMessage } from './graph-page-datasource.service';

@Injectable()
export class GraphPageDataSourceServer {
  constructor(
    private graphConverter: GraphDataConverterService,
    private connector: WindowConnectorService,
    private selectionManager: SelectionManagerService,
    private dataStore: InspectionDataStoreService,
  ) {}

  public activate() {
    this.connector
      .receiver(GRAPH_PAGE_OPEN)
      .pipe(
        withLatestFrom(
          this.selectionManager.selectedLog,
          this.dataStore.$filteredTimelines,
        ),
      )
      .subscribe(([message, log, timeline]) => {
        if (log && timeline) {
          const graphData = this.graphConverter.getGraphDataAt(
            timeline,
            log.time,
          );
          this.connector.unicast<UpdateGraphMessage>(
            UPDATE_GRAPH_DATA,
            {
              graphData,
            },
            message.sourceFrameId!,
          );
        }
      });
  }
}
