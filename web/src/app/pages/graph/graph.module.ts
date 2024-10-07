import { NgModule } from '@angular/core';

import { GraphComponent } from './graph.component';
import { ArchitectureGraphComponent } from './architecture-graph/architecture-graph.component';
import { MatToolbarModule } from '@angular/material/toolbar';
import { NgxEnvModule } from '@ngx-env/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { DownloadService } from './services/donwload-service';
import { CommonModule } from '@angular/common';
import { HeaderModule } from 'src/app/header/header.module';
import { GraphPageDataSource } from 'src/app/services/frame-connection/frames/graph-page-datasource.service';

@NgModule({
  declarations: [GraphComponent, ArchitectureGraphComponent],
  imports: [
    CommonModule,
    MatToolbarModule,
    NgxEnvModule,
    MatButtonModule,
    MatIconModule,
    HeaderModule,
  ],
  providers: [GraphPageDataSource, DownloadService],
})
export class GraphPageModule {}
