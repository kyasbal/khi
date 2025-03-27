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

import { environment } from 'src/environments/environment';
import { VERSION } from 'src/environments/version';
import {
  FrontendAnalytics,
  KHIAnalyticsActivityMetadata,
  KHIAnalyticsActivityType,
} from './types';
import { PageType } from '../../extension-common/extension-types/lifecycle-hook';

/**
 * Google Analytics for monitoring acitivities over KHI
 * Only for internal usage. Will be removed on OSSing process.
 */
export class FrontendAnalyticsWithGA implements FrontendAnalytics {
  private globalMetadataMap: KHIAnalyticsActivityMetadata = {};

  static GTAG_ADDRESS = 'https://www.googletagmanager.com/gtag/js';

  private _enabled = false;

  public init(pageType: PageType) {
    const gtagId = environment.options['GTAG_ID'] as string;
    if (gtagId && typeof gtagId === 'string' && typeof gtag === 'function') {
      this._enabled = true;
    } else {
      console.warn('GA is disabled');
    }
    if (this._enabled) {
      this.injectGTagCode();
      const viewerMode = environment.options['VIEWER_MODE'];
      const gaLabels: { [key: string]: string | boolean } = {
        debug_mode: !environment.production,
        version: VERSION + (viewerMode ? '-ro' : ''),
        pageType: pageType,
        ...this.gatherGAMetaTagLabels(),
      };
      gtag('js', new Date());
      gtag('config', gtagId, gaLabels);
    }

    this.report(KHIAnalyticsActivityType.Init);
  }

  public report(
    event: KHIAnalyticsActivityType,
    eventMetadata: KHIAnalyticsActivityMetadata = {},
  ): void {
    const analyticsMetadata = { ...this.globalMetadataMap, ...eventMetadata };
    if (this._enabled) {
      gtag('event', event, analyticsMetadata);
    } else {
      console.info(
        `[GA-event]event=${event},metadata = ${JSON.stringify(analyticsMetadata, null, 2)}`,
      );
    }
  }

  public globalMetadata(): KHIAnalyticsActivityMetadata {
    return this.globalMetadataMap;
  }

  private injectGTagCode(): void {
    const gtagId = environment.options['GTAG_ID'];
    const scriptTag = document.createElement('script');
    scriptTag.async = true;
    scriptTag.src = `${FrontendAnalyticsWithGA.GTAG_ADDRESS}?id=${gtagId}`;
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
