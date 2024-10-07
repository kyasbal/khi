import { Annotator } from '../annotator';
import { CommonFieldAnnotatorComponent } from '../common-field-annotator.component';
import { RelationshipAnnotatorComponent } from './relationship-annotator.component';
import { LogAnnotatorResolver } from './resolver';
import { TypeSeverityAnnotatorComponent } from './type-severity-annotator.component';

export function getDefaultLogAnnotatorResolver(): LogAnnotatorResolver {
  return new LogAnnotatorResolver([
    new Annotator(
      TypeSeverityAnnotatorComponent,
      TypeSeverityAnnotatorComponent.inputMapper,
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimestamp(
        'schedule',
        'Timestamp',
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.annotationDeciderForLogBodyField(
        'fingerprint',
        'InsertId',
        (l) => l['insertId'],
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForSummary(
        'summarize',
        'Summary',
      ),
    ),
    new Annotator(
      RelationshipAnnotatorComponent,
      RelationshipAnnotatorComponent.inputMapper,
    ),
  ]);
}
