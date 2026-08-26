/**
 * Copyright 2026 Google LLC
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

import { Injectable, InjectionToken, inject } from '@angular/core';
import { Client, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { PageType } from 'src/app/extensions/extension-common/extension-types/lifecycle-hook';
import { FrontendAnalyticsWithGA } from 'src/app/extensions/private/analytics/ga';
import {
  FrontendAnalytics,
  KHIAnalyticsActivityMetadata,
  KHIAnalyticsActivityType,
} from 'src/app/extensions/private/analytics/types';
import { PrivateAnalyticsService as PrivateAnalyticsConnectService } from 'src/app/generated/api/v1/private_analytics_pb';
import { ApiPathUtil } from 'src/app/services/api/api-path-util';

/**
 * Injection token for providing a custom or mock Connect-RPC client for PrivateAnalyticsService.
 */
export const PRIVATE_ANALYTICS_CONNECT_CLIENT = new InjectionToken<
  Client<typeof PrivateAnalyticsConnectService>
>('PRIVATE_ANALYTICS_CONNECT_CLIENT');

/**
 * Service to manage analytics reporting to both the KHI backend (via Connect-RPC) and Google Analytics.
 */
@Injectable()
export class PrivateAnalyticsService implements FrontendAnalytics {
  private readonly gaAnalytics =
    inject(FrontendAnalyticsWithGA, { optional: true }) ??
    new FrontendAnalyticsWithGA();

  private readonly client =
    inject(PRIVATE_ANALYTICS_CONNECT_CLIENT, { optional: true }) ??
    createClient(
      PrivateAnalyticsConnectService,
      createConnectTransport({
        baseUrl: ApiPathUtil.getServerBaseUrl(),
        useBinaryFormat: true,
      }),
    );

  /**
   * Initializes analytics providers for the current page.
   * @param pageType The type of the loaded page.
   */
  public init(pageType: PageType): void {
    this.gaAnalytics.init(pageType);
    this.sendInitToBackend(pageType);
  }

  /**
   * Reports an event to both the KHI backend and Google Analytics.
   * @param event The type of activity being recorded.
   * @param metadata Supplementary metadata for this event.
   */
  public report(
    event: KHIAnalyticsActivityType,
    metadata: KHIAnalyticsActivityMetadata = {},
  ): void {
    this.gaAnalytics.report(event, metadata);
    this.sendReportToBackend(event, metadata);
  }

  /**
   * Returns the global metadata map shared across all events.
   */
  public globalMetadata(): KHIAnalyticsActivityMetadata {
    return this.gaAnalytics.globalMetadata();
  }

  /**
   * Sends page load init activity to the backend asynchronously.
   * @param pageType The type of the loaded page.
   */
  private sendInitToBackend(pageType: PageType): void {
    this.client
      .reportActivity({
        payload: {
          case: 'init',
          value: {
            pageType: String(pageType),
          },
        },
      })
      .catch((error: unknown) => {
        console.warn(
          '[PrivateAnalyticsService] Failed to report init activity to backend:',
          error,
        );
      });
  }

  /**
   * Sends activity reporting payload to the backend asynchronously.
   * @param event The type of activity being recorded.
   * @param metadata Supplementary metadata for this event.
   */
  private sendReportToBackend(
    event: KHIAnalyticsActivityType,
    metadata: KHIAnalyticsActivityMetadata,
  ): void {
    const mergedMetadata: KHIAnalyticsActivityMetadata = {
      ...this.globalMetadata(),
      ...metadata,
    };

    switch (event) {
      case KHIAnalyticsActivityType.Inspect:
        this.client
          .reportActivity({
            payload: {
              case: 'inspect',
              value: {},
            },
          })
          .catch((error: unknown) => {
            console.warn(
              '[PrivateAnalyticsService] Failed to report inspect activity to backend:',
              error,
            );
          });
        break;

      case KHIAnalyticsActivityType.OpenInspectionData:
        this.client
          .reportActivity({
            payload: {
              case: 'openInspectionData',
              value: {
                inspectionDataHash: String(
                  mergedMetadata['inspectionDataHash'] ?? '',
                ),
                openId: String(mergedMetadata['openId'] ?? ''),
                logLength: BigInt(Number(mergedMetadata['logLength'] ?? 0)),
                decompressedTextBufferLength: BigInt(
                  Number(mergedMetadata['decompressedTextBufferLength'] ?? 0),
                ),
                revisionCount: BigInt(
                  Number(mergedMetadata['revisionCount'] ?? 0),
                ),
                eventCount: BigInt(Number(mergedMetadata['eventCount'] ?? 0)),
              },
            },
          })
          .catch((error: unknown) => {
            console.warn(
              '[PrivateAnalyticsService] Failed to report openInspectionData activity to backend:',
              error,
            );
          });
        break;

      default:
        break;
    }
  }
}
