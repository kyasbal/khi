import { Observable } from 'rxjs';
import {
  GetInspectionTasksResponse,
  GetInspectionTypesResponse,
} from 'src/app/common/schema/api-types';

/**
 * BackendConnectionService communicates the backend continuously and emit the latest information.
 */
export interface BackendConnectionService {
  /**
   * Return an observable to monitor the available task types on tge backend.
   */
  inspectionTypes(): Observable<GetInspectionTypesResponse>;

  /**
   * Return an observable to monitor the task lists on the backend.
   */
  tasks(): Observable<GetInspectionTasksResponse>;
}
