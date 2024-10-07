import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { HeaderComponent } from './header.component';
import { ToolbarComponent } from './toolbar.component';
import { NgxEnvModule } from '@ngx-env/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RegexInputComponent } from './regex-input.component';
import { SetInputComponent } from './set-input.component';
import { OverlayModule } from '@angular/cdk/overlay';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatChipsModule } from '@angular/material/chips';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { KHICommonModule } from '../common/common.module';
import { CommonModule } from '@angular/common';
import { TitleBarComponent } from './titlebar.component';
import { MatMenuModule } from '@angular/material/menu';
import { MainMenuComponent } from './main-menu.component';
import { GraphMenuComponent } from './graph-menu.component';
import { MatTooltipModule } from '@angular/material/tooltip';

@NgModule({
  declarations: [
    HeaderComponent,
    ToolbarComponent,
    SetInputComponent,
    RegexInputComponent,
    TitleBarComponent,
    MainMenuComponent,
    GraphMenuComponent,
  ],
  imports: [
    CommonModule,
    KHICommonModule,
    MatButtonModule,
    MatIconModule,
    MatToolbarModule,
    MatFormFieldModule,
    MatAutocompleteModule,
    MatChipsModule,
    MatInputModule,
    MatMenuModule,
    ReactiveFormsModule,
    FormsModule,
    OverlayModule,
    NgxEnvModule,
    MatTooltipModule,
  ],
  exports: [HeaderComponent, TitleBarComponent, GraphMenuComponent],
})
export class HeaderModule {}
