import { InjectionToken } from '@angular/core';
import { AnnotatorResolver } from '../annotator';
import { TimelineEntry } from 'src/app/store/timeline';

export const NAVIGATOR_ANNOTATOR_RESOLVER =
  new InjectionToken<NavigatorAnnotatorResolver>(
    'NAVIGATOR_ANNOTATOR_RESOLVER',
  );

export class NavigatorAnnotatorResolver extends AnnotatorResolver<TimelineEntry> {}
