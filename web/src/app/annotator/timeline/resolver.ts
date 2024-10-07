import { InjectionToken } from '@angular/core';
import { AnnotatorResolver } from '../annotator';
import { TimelineEntry } from 'src/app/store/timeline';

export const TIMELINE_ANNOTATOR_RESOLVER =
  new InjectionToken<TimelineAnnotatorResolver>('TIMELINE_ANNOTATION_RESOLVER');

export class TimelineAnnotatorResolver extends AnnotatorResolver<TimelineEntry> {}
