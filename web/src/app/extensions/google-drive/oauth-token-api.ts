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

import { Inject, Injectable, InjectionToken } from '@angular/core';
import { Subject, Subscription } from 'rxjs';

interface OAuthTokenResponse {
  access_token: string;
}

/**
 * Partial type definition of https://accounts.google.com/gsi/client
 */
interface GoogleOAuth2ClientLibrary {
  initTokenClient(request: {
    client_id: string;
    scope: string;
    callback: (tokenResponse: OAuthTokenResponse) => void;
  }): GoogleOAuth2Client;
}

interface GoogleOAuth2Client {
  requestAccessToken(): void;
}

export const GOOGLE_OAUTH2_LIB = new InjectionToken<GoogleOAuth2Client>(
  'khi-google-oauth2-library',
);

export const LOCALSTORAGE = new InjectionToken<typeof window.localStorage>(
  'khi-localStorage',
);

@Injectable({
  providedIn: 'root',
})
export class OAuthTokenAPI {
  public static readonly LOCALSTORAGE_OAUTH_TOKEN_KEY = 'khi-google-oauth-key';
  private static readonly OAUTH_SCOPE =
    'https://www.googleapis.com/auth/drive.readonly';
  private readonly oauthObservable: Subject<OAuthTokenResponse> = new Subject();
  private oAuthClient: GoogleOAuth2Client;

  constructor(
    @Inject(GOOGLE_OAUTH2_LIB)
    oAuth2Lib: GoogleOAuth2ClientLibrary,
    @Inject(LOCALSTORAGE)
    private localStorage: typeof window.localStorage,
  ) {
    this.oAuthClient = oAuth2Lib.initTokenClient({
      client_id:
        '173883494332-f5o6gkmpb5fku5u3vp43279mqbj24pe5.apps.googleusercontent.com', // TODO: Client ID is not secret and this code won't be released as a part of OSS. This should be specified on some configuration file, but I will keep this here.
      scope: OAuthTokenAPI.OAUTH_SCOPE,
      callback: (tokenResponse: OAuthTokenResponse) => {
        this.oauthObservable.next(tokenResponse);
      },
    });
  }

  public get tokenFromCache(): string | null {
    return this.localStorage.getItem(
      OAuthTokenAPI.LOCALSTORAGE_OAUTH_TOKEN_KEY,
    );
  }

  /**
   * Attempt to obtain the access token with OAuth
   * Must be called with user interaction to show the popup
   */
  public async requestAccessToken(disableCache: boolean): Promise<string> {
    return new Promise((resolve, error) => {
      if (!disableCache && this.tokenFromCache !== null) {
        resolve(this.tokenFromCache);
        return;
      }
      let oauthTokenSubscription: Subscription | null = null;
      const onTokenReceive = (response: OAuthTokenResponse) => {
        if (oauthTokenSubscription) oauthTokenSubscription.unsubscribe();
        if (response && response.access_token) {
          this.saveToken(response.access_token);
          resolve(response.access_token);
          return;
        }
        this.saveToken(null);
        error(
          'Failed to request access token. Received response:\n' +
            JSON.stringify(response),
        );
      };
      oauthTokenSubscription = this.oauthObservable.subscribe(onTokenReceive);
      this.oAuthClient.requestAccessToken();
    });
  }

  /**
   * Utility function to call Web API with access token given from OAuth.
   * If the first call fail and regarded as retry, the access token will be obtained from OAuth without using cache.
   * @param apiCall Directly calling the API
   * @param shouldAuthRetry Check if the raised error should be retried or not
   * @returns Result of the API call
   */
  public async requestAPIWithOAuth<T>(
    apiCall: (token: string) => Promise<T>,
    shouldAuthRetry: (err: XMLHttpRequest) => boolean,
  ): Promise<T> {
    const token = await this.requestAccessToken(false);
    try {
      return await apiCall(token);
    } catch (e) {
      if (shouldAuthRetry(e as XMLHttpRequest)) {
        try {
          const refreshedToken = await this.requestAccessToken(true);
          return await apiCall(refreshedToken);
        } catch (e) {
          throw new Error(
            'Failed to call api with refreshed OAuth token.\n' + e,
          );
        }
      }
      throw new Error(
        'Failed to call api with existing OAuth token. This error was regarded not to retry.\n' +
          e,
      );
    }
  }

  public clearTokenCache() {
    this.localStorage.removeItem(OAuthTokenAPI.LOCALSTORAGE_OAUTH_TOKEN_KEY);
  }

  private saveToken(accessToken: string | null) {
    if (accessToken == null) {
      this.clearTokenCache();
      return;
    }
    this.localStorage.setItem(
      OAuthTokenAPI.LOCALSTORAGE_OAUTH_TOKEN_KEY,
      accessToken,
    );
  }
}
