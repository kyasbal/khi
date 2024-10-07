import { NgModule } from '@angular/core';
import { LogBodyComponent } from './body.component';
import { LogHeaderComponent } from './header.component';
import { LogViewComponent } from './log-view.component';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { ScrollingModule } from '@angular/cdk/scrolling';
import { HighlightModule, provideHighlightOptions } from 'ngx-highlightjs';
import { KHICommonModule } from '../common/common.module';
import { CommonModule, NgComponentOutlet } from '@angular/common';
import { IconToggleButtonComponent } from './icon-toggle-button.component';
import { MatTooltipModule } from '@angular/material/tooltip';
import { LOG_ANNOTATOR_RESOLVER } from '../annotator/log/resolver';
import { HighlightLineNumbers } from 'ngx-highlightjs/line-numbers';
import { getDefaultLogAnnotatorResolver } from '../annotator/log/default';
import { LOG_TOOL_ANNOTATOR_RESOLVER } from '../annotator/log-tool/resolver';
import { getDefaultLogToolAnnotatorResolver } from '../annotator/log-tool/default';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { LogViewLogLineComponent } from './log-view-log-line.component';

@NgModule({
  declarations: [
    LogBodyComponent,
    LogHeaderComponent,
    LogViewComponent,
    LogViewLogLineComponent,
    IconToggleButtonComponent,
  ],
  imports: [
    CommonModule,
    KHICommonModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatSlideToggleModule,
    FormsModule,
    ReactiveFormsModule,
    ScrollingModule,
    HighlightModule,
    HighlightLineNumbers,
    ClipboardModule,
    MatTooltipModule,
    NgComponentOutlet,
    MatSnackBarModule,
  ],
  providers: [
    provideHighlightOptions({
      coreLibraryLoader: () => import('highlight.js/lib/core'),
      lineNumbersLoader: () => import('ngx-highlightjs/line-numbers'),
      languages: {
        yaml: () => import('highlight.js/lib/languages/yaml'),
      },
    }),
    {
      provide: LOG_ANNOTATOR_RESOLVER,
      useValue: getDefaultLogAnnotatorResolver(),
    },
    {
      provide: LOG_TOOL_ANNOTATOR_RESOLVER,
      useValue: getDefaultLogToolAnnotatorResolver(),
    },
  ],
  exports: [LogViewComponent, LogHeaderComponent],
})
export class LogModule {}
