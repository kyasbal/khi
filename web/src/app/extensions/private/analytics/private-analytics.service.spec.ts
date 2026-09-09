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

import { TestBed } from '@angular/core/testing';
import { Client } from '@connectrpc/connect';
import { PageType } from 'src/app/extensions/extension-common/extension-types/lifecycle-hook';
import { FrontendAnalyticsWithGA } from 'src/app/extensions/private/analytics/ga';
import {
  PRIVATE_ANALYTICS_CONNECT_CLIENT,
  PrivateAnalyticsService,
} from 'src/app/extensions/private/analytics/private-analytics.service';
import {
  KHIAnalyticsActivityMetadata,
  KHIAnalyticsActivityType,
} from 'src/app/extensions/private/analytics/types';
import {
  PrivateAnalyticsService as PrivateAnalyticsConnectService,
  ReportActivityResponse,
} from 'src/app/generated/api/v1/private_analytics_pb';
import { VERSION } from 'src/environments/version';

describe('PrivateAnalyticsService', () => {
  let service: PrivateAnalyticsService;
  let mockGA: jasmine.SpyObj<FrontendAnalyticsWithGA>;
  let mockClient: jasmine.SpyObj<Client<typeof PrivateAnalyticsConnectService>>;
  let sharedMetadata: KHIAnalyticsActivityMetadata;

  beforeEach(() => {
    sharedMetadata = {};
    mockGA = jasmine.createSpyObj<FrontendAnalyticsWithGA>(
      'FrontendAnalyticsWithGA',
      ['init', 'report', 'globalMetadata'],
    );
    mockGA.globalMetadata.and.returnValue(sharedMetadata);

    mockClient = jasmine.createSpyObj<
      Client<typeof PrivateAnalyticsConnectService>
    >('PrivateAnalyticsConnectService', ['reportActivity']);
    mockClient.reportActivity.and.returnValue(
      Promise.resolve({} as ReportActivityResponse),
    );

    TestBed.configureTestingModule({
      providers: [
        { provide: FrontendAnalyticsWithGA, useValue: mockGA },
        { provide: PRIVATE_ANALYTICS_CONNECT_CLIENT, useValue: mockClient },
        PrivateAnalyticsService,
      ],
    });

    service = TestBed.inject(PrivateAnalyticsService);
  });

  it('should send INIT event to backend via Connect-RPC and call ga.init() on init()', () => {
    service.init(PageType.Main);

    expect(mockGA.init).toHaveBeenCalledWith(PageType.Main);
    expect(mockClient.reportActivity).toHaveBeenCalledWith({
      payload: {
        case: 'init',
        value: {
          pageType: String(PageType.Main),
          frontendVersion: VERSION,
        },
      },
    });
  });

  it('should not send unrecognized activity types to backend via Connect-RPC but still call ga.report()', () => {
    const unknownEvent = 'UNKNOWN_EVENT' as KHIAnalyticsActivityType;
    service.report(unknownEvent, { key: 'val' });

    expect(mockGA.report).toHaveBeenCalledWith(unknownEvent, { key: 'val' });
    expect(mockClient.reportActivity).not.toHaveBeenCalled();
  });

  it('should send INSPECT event to backend via Connect-RPC and call ga.report() on report()', () => {
    service.report(KHIAnalyticsActivityType.Inspect, {});

    expect(mockGA.report).toHaveBeenCalledWith(
      KHIAnalyticsActivityType.Inspect,
      {},
    );
    expect(mockClient.reportActivity).toHaveBeenCalledWith({
      payload: {
        case: 'inspect',
        value: {},
      },
    });
  });

  it('should send OPEN_INSPECTION_DATA event with merged metadata to backend via Connect-RPC on report()', () => {
    sharedMetadata['inspectionDataHash'] = 'hash-abc';
    sharedMetadata['openId'] = 'open-xyz';

    service.report(KHIAnalyticsActivityType.OpenInspectionData, {
      logLength: 500,
      decompressedTextBufferLength: 10000,
      revisionCount: 8,
      eventCount: 20,
    });

    expect(mockGA.report).toHaveBeenCalledWith(
      KHIAnalyticsActivityType.OpenInspectionData,
      {
        logLength: 500,
        decompressedTextBufferLength: 10000,
        revisionCount: 8,
        eventCount: 20,
      },
    );
    expect(mockClient.reportActivity).toHaveBeenCalledWith({
      payload: {
        case: 'openInspectionData',
        value: {
          inspectionDataHash: 'hash-abc',
          openId: 'open-xyz',
          logLength: 500n,
          decompressedTextBufferLength: 10000n,
          revisionCount: 8n,
          eventCount: 20n,
        },
      },
    });
  });

  it('should handle Connect-RPC failure gracefully without throwing and log a warning on report()', async () => {
    const warnSpy = spyOn(console, 'warn');
    mockClient.reportActivity.and.returnValue(
      Promise.reject(new Error('Network error')),
    );

    service.report(KHIAnalyticsActivityType.Inspect, {});
    await Promise.resolve();

    expect(warnSpy).toHaveBeenCalledWith(
      '[PrivateAnalyticsService] Failed to report inspect activity to backend:',
      jasmine.any(Error),
    );
  });

  it('should handle Connect-RPC failure gracefully without throwing and log a warning on init()', async () => {
    const warnSpy = spyOn(console, 'warn');
    mockClient.reportActivity.and.returnValue(
      Promise.reject(new Error('Network error on init')),
    );

    service.init(PageType.Main);
    await Promise.resolve();

    expect(warnSpy).toHaveBeenCalledWith(
      '[PrivateAnalyticsService] Failed to report init activity to backend:',
      jasmine.any(Error),
    );
  });

  it('should share globalMetadata reference with gaAnalytics', () => {
    expect(service.globalMetadata()).toBe(sharedMetadata);
  });
});
