import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TimelineComponent } from './timeline.component';
import { ScrollingModule } from '@angular/cdk/scrolling';
import { OverlayModule } from '@angular/cdk/overlay';
import { MatIconModule } from '@angular/material/icon';
import { KHICommonModule } from '../common/common.module';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { NgxEnvModule } from '@ngx-env/core';
import { NavigatorComponent } from './navigator/navigator.component';
import { NAVIGATOR_ANNOTATOR_RESOLVER } from '../annotator/navigator/resolver';
import { getDefaultNavigatorAnnotatorResolver } from '../annotator/navigator/default';

@NgModule({
  declarations: [TimelineComponent, NavigatorComponent],
  imports: [
    CommonModule,
    KHICommonModule,
    ScrollingModule,
    OverlayModule,
    MatIconModule,
    MatButtonModule,
    MatTooltipModule,
    NgxEnvModule,
  ],
  providers: [
    {
      provide: NAVIGATOR_ANNOTATOR_RESOLVER,
      useValue: getDefaultNavigatorAnnotatorResolver(),
    },
  ],
  exports: [TimelineComponent],
})
export class TimelineModule {}
