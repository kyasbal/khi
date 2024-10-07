import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  PopupAnswerResponse,
  PopupAnswerValidationResult,
  PopupFormRequest,
} from 'src/app/common/schema/api-types';

/**
 * The injection token to get the actual PopupManager implementation.
 */
export const POPUP_MANAGER = new InjectionToken<PopupManager>('POPUP_MANAGER');

export interface PopupManager {
  /**
   * Get the observable stream to monitor popup form requests.
   */
  requests(): Observable<PopupFormRequestWithClient>;
}

export interface PopupFormRequestWithClient extends PopupFormRequest {
  client: PopupClient;
}

export interface PopupClient {
  /**
   * Validate if the content is valid or not.
   * @param data the data to verify as the response of popup request
   */
  validate(data: PopupAnswerResponse): Observable<PopupAnswerValidationResult>;

  /**
   * Send the answer for the popup request.
   * @param data
   */
  answer(data: PopupAnswerResponse): Observable<void>;
}
