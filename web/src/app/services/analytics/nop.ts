import { FrontendAnalytics, KHIAnalyticsActivityMetadata } from './types';

/**
 * NopFrontendAnalytics is an implementation of FrontendAnalytics for testing purpose.
 */
export class NopFrontendAnalytics implements FrontendAnalytics {
  init(): void {
    return;
  }
  report(): void {
    return;
  }
  globalMetadata(): KHIAnalyticsActivityMetadata {
    return {};
  }
}
