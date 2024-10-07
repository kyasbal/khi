import { InjectionToken } from '@angular/core';

export const FRONTEND_ANALYTICS = new InjectionToken<FrontendAnalytics>(
  'FRONTEND_ANAYTICS',
);

/**
 * Type of activity recorded in Analytics
 */
export enum KHIAnalyticsActivityType {
  Init = 'INIT',
  Inspect = 'INSPECT',
  OpenInspectionData = 'OPEN_INSPECTION_DATA',
  ClickRevision = 'CLICK_REVISION',
  ClickEvent = 'CLICK_EVENT',
  RenderGraph = 'RENDER_GRAPH',
}

export enum KHIAnalyticsPageType {
  Main = 'MAIN',
  GraphView = 'GRAPH_VIEW',
  DiffView = 'DIFF_VIEW',
}

export type KHIAnalyticsActivityMetadata = {
  [metadataKey: string]: number | string;
};

/**
 * FrontendAnalytics is an interface to report events.
 */
export interface FrontendAnalytics {
  /**
   * initialize analytics provider. This must be called when a page loaded.
   */
  init(pageType: KHIAnalyticsPageType): void;

  /**
   * Report the event.
   * @param event type of event
   * @param metadata supplimental information attached to this event.
   */
  report(
    event: KHIAnalyticsActivityType,
    metadata: KHIAnalyticsActivityMetadata,
  ): void;

  /**
   * Returns the reference to global metadata in addition to the metadata argument reported with `report` method.
   */
  globalMetadata(): KHIAnalyticsActivityMetadata;
}
