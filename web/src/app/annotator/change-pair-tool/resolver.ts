import { InjectionToken } from '@angular/core';
import { AnnotatorResolver } from '../annotator';
import { ResourceRevisionChangePair } from 'src/app/store/timeline';

export const CHANGE_PAIR_TOOL_ANNOTATOR_RESOLVER =
  new InjectionToken<ChangePairAnnotatorResolver>(
    'CHANGE_PAIR_TOOL_ANNOTATOR_RESOLVER',
  );

export const CHANGE_PAIR_TOOL_ANNOTATOR_FOR_FLOATING_PAGE_RESOLVER =
  new InjectionToken<ChangePairAnnotatorResolver>(
    'CHANGE_PAIR_TOOL_ANNOTATOR_FOR_FLOATING_PAGE_RESOLVER',
  );

export class ChangePairAnnotatorResolver extends AnnotatorResolver<ResourceRevisionChangePair> {}
