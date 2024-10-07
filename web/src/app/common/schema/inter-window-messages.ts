import { TimelineEntry } from 'src/app/store/timeline';

export const DIFF_PAGE_OPEN = 'DIFF_PAGE_OPEN';
export const UPDATE_SELECTED_RESOURCE_MESSAGE_KEY = 'UPDATE_SELECTED_RESOURCE';
export const UPDATE_GRAPH_DATA = 'UPDATE_GRAPH_DATA';
export const GRAPH_PAGE_OPEN = 'GRAPH_PAGE_OPEN';
/**
 * Main window broadcast this message when another resource was selected.
 */
export interface UpdateSelectedResourceMessage {
  timeline: TimelineEntry;
  logIndex: number;
}
/**
 * A viewmodel for entire diff page.
 */
export interface DiffPageViewModel {
  timeline: TimelineEntry;
  logIndex: number;
}
