import {
  distinctUntilKeyChanged,
  filter,
  interval,
  map,
  Observable,
  retry,
  shareReplay,
  switchMap,
  throwError,
} from 'rxjs';
import {
  PopupClient,
  PopupFormRequestWithClient,
  PopupManager,
} from './popup-manager';
import {
  PopupAnswerResponse,
  PopupAnswerValidationResult,
  PopupFormRequest,
} from 'src/app/common/schema/api-types';
import { BACKEND_API, BackendAPI } from '../api/backend-api-interface';
import { Inject, Injectable } from '@angular/core';

@Injectable({ providedIn: 'any' })
export class PopupManagerImpl implements PopupManager {
  private popupRequest = interval(1000).pipe(
    switchMap(() => this.backendAPI.getPopup() as Observable<PopupFormRequest>),
    filter((pr) => !!pr),
    distinctUntilKeyChanged('id'),
    retry(),
    shareReplay({
      bufferSize: 1,
      refCount: true,
    }),
  );

  constructor(@Inject(BACKEND_API) private backendAPI: BackendAPI) {}

  requests(): Observable<PopupFormRequestWithClient> {
    return this.popupRequest.pipe(
      map((request) => ({
        ...request,
        client: new PopupClientImpl(request.id, this.backendAPI),
      })),
    );
  }
}

export class PopupClientImpl implements PopupClient {
  constructor(
    public readonly popupId: string,
    private backendAPI: BackendAPI,
  ) {}
  validate(data: PopupAnswerResponse): Observable<PopupAnswerValidationResult> {
    if (data.id === this.popupId) {
      return this.backendAPI.validatePopupAnswer(data);
    } else {
      return throwError(() => {
        return 'the popup id is not for this client';
      });
    }
  }
  answer(data: PopupAnswerResponse): Observable<void> {
    if (data.id === this.popupId) {
      return this.backendAPI.answerPopup(data);
    } else {
      return throwError(() => {
        return 'the popup id is not for this client';
      });
    }
  }
}
