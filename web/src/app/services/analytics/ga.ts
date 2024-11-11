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

import {
  FrontendAnalytics,
  KHIAnalyticsActivityMetadata,
  KHIAnalyticsActivityType,
  KHIAnalyticsPageType,
} from './types';

/**
 * Google Analytics for monitoring acitivities over KHI
 * Only for internal usage. Will be removed on OSSing process.
 */
export class FrontendAnalyticsWithGA implements FrontendAnalytics {
  private globalMetadataMap: KHIAnalyticsActivityMetadata = {};

  static GTAG_ADDRESS = 'https://www.googletagmanager.com/gtag/js';

  private _enabled = false;

  public init(pageType: KHIAnalyticsPageType) {
    if (process.env['NG_APP_GTAG_ID'] && typeof gtag === 'function') {
      this._enabled = true;
    } else {
      console.warn('GA is disabled');
      return;
    }
    this.injectGTagCode();

    const isViewer =
      process.env['NG_APP_VIEWER_MODE'] &&
      process.env['NG_APP_VIEWER_MODE'] != 'false';
    const gaLabels: { [key: string]: string | boolean } = {
      debug_mode: process.env['NG_APP_GTAG_DEBUG'] === 'true',
      version: process.env['NG_APP_VERSION'] + (isViewer ? '-ro' : ''),
      pageType: pageType,
      ...this.gatherGAMetaTagLabels(),
    };
    gtag('js', new Date());
    gtag('config', process.env['NG_APP_GTAG_ID'], gaLabels);

    this.report(KHIAnalyticsActivityType.Init);
  }

  public report(
    event: KHIAnalyticsActivityType,
    metadata: KHIAnalyticsActivityMetadata = {},
  ): void {
    if (this._enabled)
      gtag('event', event, { ...this.globalMetadataMap, ...metadata });
  }

  public globalMetadata(): KHIAnalyticsActivityMetadata {
    return this.globalMetadataMap;
  }

  private injectGTagCode(): void {
    const scriptTag = document.createElement('script');
    scriptTag.async = true;
    scriptTag.src = `${FrontendAnalyticsWithGA.GTAG_ADDRESS}?id=${process.env['NG_APP_GTAG_ID']}`;
    document.head.appendChild(scriptTag);
  }

  /**
   * Gather meta values injected from backend to send Google Analytics.
   */
  private gatherGAMetaTagLabels(): { [key: string]: string } {
    const metaTagParameters: { [key: string]: string } = {};
    const metaTags = document.getElementsByTagName('meta');
    for (let i = 0; i < metaTags.length; i++) {
      const metaTag = metaTags[i];
      if (!metaTag.id.startsWith('ga-meta-')) continue;
      // If the metadata id was `ga-meta-user`, it will be treated as the `user_id` in Google Analytics
      if (
        metaTag.id === 'ga-meta-user' &&
        metaTag.getAttribute('content') !== null
      ) {
        metaTagParameters['user_id'] = metaTag.getAttribute('content')!;
      } else {
        metaTagParameters[metaTag.id] =
          metaTag.getAttribute('content') ?? '(empty)';
      }
    }
    return metaTagParameters;
  }
}
