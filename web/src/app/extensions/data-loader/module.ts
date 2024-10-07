import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { GoogleDriveAPI } from './google-drive-api';
import { LoginDialogComponent } from './login-dialog/login.component';
import {
  GOOGLE_OAUTH2_LIB,
  LOCALSTORAGE,
  OAuthTokenAPI,
} from './oauth-token-api';

@NgModule({
  declarations: [LoginDialogComponent],
  imports: [CommonModule, MatDialogModule, MatButtonModule],
  providers: [
    {
      provide: GOOGLE_OAUTH2_LIB,
      useValue: 'google' in window ? google.accounts.oauth2 : null,
    },
    {
      provide: LOCALSTORAGE,
      useValue: window.localStorage,
    },
    GoogleDriveAPI,
    OAuthTokenAPI,
  ],
  exports: [LoginDialogComponent],
})
export class GoogleDriveDataLoaderModule {}
