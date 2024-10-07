import { InjectionToken } from '@angular/core';
import { AnnotatorResolver } from '../annotator';
import { LogEntry } from 'src/app/store/log';

export const LOG_ANNOTATOR_RESOLVER = new InjectionToken<LogAnnotatorResolver>(
  'LOG_ANNOTATOR_RESOLVER',
);

export class LogAnnotatorResolver extends AnnotatorResolver<LogEntry> {}
