import { NgModule } from '@angular/core';
import {
  LongTimestampFormatPipe,
  TimestampFormatPipe,
} from './timestamp-format.pipe';
import { FirstOrUndefined } from './first-or-null.pipe';
import { ParsePathPipe } from './parse-path.pipe';
import { RainbowPipe } from './rainbow.pipe';
import { SidePaneComponent } from './components/side-pane.component';
import { CommonModule } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { ResolveTextPipe } from './resolve-text.pipe';
import { MetaTableRowComponent } from './components/meta-table-row.component';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { MatTooltipModule } from '@angular/material/tooltip';
import { CssClassFormatPipe } from './css-class-format.pipe';
import { BreaklinePipe } from './breakline.pipe';
import { CaptureShiftKeyDirective } from './capture-shiftkey.directive';

@NgModule({
  declarations: [
    TimestampFormatPipe,
    LongTimestampFormatPipe,
    FirstOrUndefined,
    ParsePathPipe,
    RainbowPipe,
    CssClassFormatPipe,
    ResolveTextPipe,
    SidePaneComponent,
    MetaTableRowComponent,
    BreaklinePipe,
    CaptureShiftKeyDirective,
  ],
  imports: [
    CommonModule,
    MatToolbarModule,
    MatIconModule,
    ClipboardModule,
    MatTooltipModule,
  ],
  exports: [
    TimestampFormatPipe,
    LongTimestampFormatPipe,
    FirstOrUndefined,
    ParsePathPipe,
    RainbowPipe,
    ResolveTextPipe,
    CssClassFormatPipe,
    SidePaneComponent,
    MetaTableRowComponent,
    BreaklinePipe,
    CaptureShiftKeyDirective,
  ],
})
export class KHICommonModule {}
