import { InjectionToken } from '@angular/core';
import { ChangePairAnnotatorResolver } from '../change-pair-tool/resolver';

export const CHANGE_PAIR_ANNOTATOR_RESOLVER =
  new InjectionToken<ChangePairAnnotatorResolver>(
    'CHANGE_PAIR_ANNOTATOR_RESOLVER',
  );
