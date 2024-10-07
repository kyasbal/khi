import { Injectable } from '@angular/core';
import { OAuthTokenAPI } from './oauth-token-api';

@Injectable({ providedIn: 'root' })
export class GoogleDriveAPI {
  constructor(private _oauthTokenAPI: OAuthTokenAPI) {}

  public getFileAsText(fileId: string): Promise<ArrayBuffer> {
    return this._oauthTokenAPI.requestAPIWithOAuth(
      (token) => this.getFileAsTextWithToken(fileId, token),
      (xhr) => xhr.status < 500 && xhr.status >= 400,
    );
  }

  private getFileAsTextWithToken(
    fileId: string,
    accessToken: string,
  ): Promise<ArrayBuffer> {
    return new Promise((resolve, error) => {
      const xhr = new XMLHttpRequest();
      xhr.responseType = 'arraybuffer';
      xhr.open(
        'GET',
        'https://www.googleapis.com/drive/v3/files/' + fileId + '?alt=media',
      );
      xhr.setRequestHeader('Authorization', 'Bearer ' + accessToken);
      xhr.onload = () => {
        if (xhr.status >= 400 && xhr.status < 600) {
          error(xhr);
        }
        resolve(xhr.response);
      };
      xhr.onerror = () => {
        error(xhr);
      };
      xhr.send();
    });
  }
}
