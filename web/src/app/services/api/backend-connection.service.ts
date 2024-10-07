import { Inject, Injectable, InjectionToken } from '@angular/core';
import { Observable, interval, retry, shareReplay, switchMap, tap } from 'rxjs';
import { BACKEND_API, BackendAPI } from './backend-api-interface';
import { BackendConnectionService } from './backend-connection-interface';
import {
  GetInspectionTypesResponse,
  GetInspectionTasksResponse,
} from 'src/app/common/schema/api-types';

/**
 * Angular injection token for BackendConnectionService.
 */
export const BACKEND_CONNECTION = new InjectionToken<BackendConnectionService>(
  'BACKEND_CONNECTION',
);

/**
 * BackendConnectionService provides observables with polling backend endpoints.
 */
@Injectable()
export class BackendConnectionServiceImpl implements BackendConnectionService {
  /**
   * Interval to poll task progresses.
   */
  static readonly PROGRESS_POLLING_INTERVAL = 1000;

  /**
   * Interval to poll the list of inspection types.
   */
  static readonly LIST_INSPECTION_TYPES_RETRY_TIME = 1000;

  private inspectionTypesObservable = interval(
    BackendConnectionServiceImpl.LIST_INSPECTION_TYPES_RETRY_TIME,
  ).pipe(
    switchMap(() => this.backendApi.getInspectionTypes()),
    retry(),
    shareReplay({
      bufferSize: 1,
      // refCount is explcitly false to prevent waiting the next poll when a new subscriber added when there is no subscriber registered.
      refCount: false,
    }),
  );

  private taskProgressObservable = interval(
    BackendConnectionServiceImpl.PROGRESS_POLLING_INTERVAL,
  ).pipe(
    switchMap(() => this.backendApi.getTaskStatuses()),
    tap({
      error: (err) => {
        console.warn(
          `Failed to refresh task progerss status:\n` + JSON.stringify(err),
        );
      },
    }),
    shareReplay({
      bufferSize: 1,
      // refCount is explcitly false to prevent waiting the next poll when a new subscriber added when there is no subscriber registered.
      refCount: false,
    }),
    retry(),
  );

  constructor(@Inject(BACKEND_API) private backendApi: BackendAPI) {}

  inspectionTypes(): Observable<GetInspectionTypesResponse> {
    return this.inspectionTypesObservable;
  }
  tasks(): Observable<GetInspectionTasksResponse> {
    return this.taskProgressObservable;
  }
}
