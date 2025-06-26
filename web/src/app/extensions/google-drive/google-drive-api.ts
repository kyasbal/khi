/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Injectable, inject } from '@angular/core';
import { OAuthTokenAPI } from './oauth-token-api';

@Injectable({ providedIn: 'root' })
export class GoogleDriveAPI {
  private _oauthTokenAPI = inject(OAuthTokenAPI);

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
