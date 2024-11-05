import { Observable, of } from 'rxjs';
import {
  PopupAnswerResponse,
  PopupAnswerValidationResult,
} from 'src/app/common/schema/api-types';
import {
  PopupClient,
  PopupManager,
  PopupFormRequestWithClient,
} from './popup-manager';

/**
 * MockPopupClient is a mock PopupClient implementation for testing purpose.
 * This client will return no validation error only when the given string is `valid`, otherwise this client will return validation error.
 */
export class MockPopupClient implements PopupClient {
  validate(data: PopupAnswerResponse): Observable<PopupAnswerValidationResult> {
    if (data.value === 'valid') {
      return of<PopupAnswerValidationResult>({
        id: data.id,
        validationError: '',
      });
    }
    return of<PopupAnswerValidationResult>({
      id: data.id,
      validationError: `${data.value} isn't valid`,
    });
  }
  answer(): Observable<void> {
    return of(void 0);
  }
}

/**
 * MockPopupManager is a PopupManager implementation for testing purpose. It always returns PopupFormRequest for testing.
 */
export class MockPopupManager implements PopupManager {
  requests(): Observable<PopupFormRequestWithClient> {
    return of({
      id: 'test',
      title: 'Google Admin inspection token',
      type: 'text',
      description:
        'Google Admin token seems to be expired. Please copy the KHI command again from Google Admin and paste it here.',
      placeholder: `gcloud auth configure-docker --quiet
docker run --rm --pull always -p 8080:8080 -it -e KHI_FIXED_PROJECT_ID="XXXXXXXXX" -e GCP_DEFAULT_PROJECT="XXXXXXXXX" -e GCP_ACCESS_TOKEN=\`gcloud auth print-access-token\` -e GCP_IDENTITY_TOKEN=\`gcloud auth print-identity-token\` -e KHI_GA_LABELS="justification=vector/123456,user=XXXXXXXXX" -e IAM_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" gcr.io/kubernetes-history-inspector/standalone:latest`,
      client: new MockPopupClient(),
      options: {},
    });
  }
}
