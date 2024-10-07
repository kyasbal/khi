import { Annotator } from '../annotator';
import { ChangePairAnnotatorResolver } from '../change-pair-tool/resolver';
import { CommonWarningMessageComponent } from '../common-warning-message.component';

export function getDefaultChangePairAnnotatorResolver(): ChangePairAnnotatorResolver {
  return new ChangePairAnnotatorResolver([
    new Annotator(
      CommonWarningMessageComponent,
      CommonWarningMessageComponent.inputMapperForRevisionPair(
        'warning',
        () => 'This is a deletion request',
        (p) => p.current && p.current.isDeletion,
      ),
    ),
    new Annotator(
      CommonWarningMessageComponent,
      CommonWarningMessageComponent.inputMapperForRevisionPair(
        'warning',
        () => 'No change made by this request',
        (p) => p.previous?.resourceContent === p.current.resourceContent,
      ),
    ),
  ]);
}
