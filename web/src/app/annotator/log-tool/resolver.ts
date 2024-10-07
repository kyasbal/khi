import { InjectionToken } from '@angular/core';
import { LogAnnotatorResolver } from '../log/resolver';

export const LOG_TOOL_ANNOTATOR_RESOLVER =
  new InjectionToken<LogAnnotatorResolver>('LOG_TOOL_ANNOTATOR_RESOLVER');
