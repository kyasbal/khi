import { Injectable } from '@angular/core';
import {
  GetInspectionTypesResponse,
  CreateInspectionTaskResponse,
  GetInspectionTaskFeatureResponse,
  PutInspectionTaskFeatureRequest,
  InspectionFeature,
  InspectionDryRunResponse,
  GetInspectionTasksResponse,
  InspectionMetadataResponse,
  InspectionDryRunRequest,
  InspectionRunRequest,
  PopupAnswerResponse,
  PopupAnswerValidationResult,
  PopupFormRequest,
} from '../../common/schema/api-types';
import { HttpClient, HttpEventType, HttpRequest } from '@angular/common/http';
import {
  Observable,
  ReplaySubject,
  Subject,
  combineLatest,
  concat,
  debounceTime,
  last,
  map,
  of,
  shareReplay,
  switchMap,
  tap,
  withLatestFrom,
} from 'rxjs';
import { ViewStateService } from '../view-state.service';
import { BackendAPI, DownloadProgressReporter } from './backend-api-interface';

const BASE_URL = process.env['NG_APP_BACKEND_ROOT_URL']
  ? process.env['NG_APP_BACKEND_ROOT_URL']
  : '/';

@Injectable({
  providedIn: 'root',
})
export class BackendAPIImpl implements BackendAPI {
  constructor(
    private http: HttpClient,
    private readonly viewState: ViewStateService,
  ) {}
  public getInspectionTypes() {
    const url = BASE_URL + 'api/v2/inspection/types';
    return this.http.get<GetInspectionTypesResponse>(url);
  }

  public getTaskStatuses() {
    const url = BASE_URL + 'api/v2/inspection/tasks';
    return this.http.get<GetInspectionTasksResponse>(url);
  }

  public createInspection(
    inspectionTypeId: string,
  ): Observable<InspectionTaskClient> {
    const url = BASE_URL + 'api/v2/inspection/types/' + inspectionTypeId;
    return this.http
      .post<CreateInspectionTaskResponse>(url, null)
      .pipe(
        map(
          (response) =>
            new InspectionTaskClient(
              this,
              response.inspectionId,
              this.viewState,
            ),
        ),
      );
  }

  public getFeatureList(taskId: string) {
    const url = BASE_URL + `api/v2/inspection/tasks/${taskId}/features`;
    return this.http.get<GetInspectionTaskFeatureResponse>(url);
  }

  public setEnabledFeatures(taskId: string, featureIds: string[]) {
    const url = BASE_URL + `api/v2/inspection/tasks/${taskId}/features`;
    const request: PutInspectionTaskFeatureRequest = {
      features: featureIds,
    };
    return this.http.put(url, request, {
      responseType: 'text',
    }) as Observable<unknown> as Observable<void>;
  }

  public getInspectionMetadata(taskId: string) {
    const url = BASE_URL + `api/v2/inspection/tasks/${taskId}/metadata`;
    return this.http.get<InspectionMetadataResponse>(url);
  }

  public runTask(
    taskId: string,
    request: InspectionRunRequest,
  ): Observable<void> {
    const url = BASE_URL + `api/v2/inspection/tasks/${taskId}/run`;
    return this.http
      .post(url, request, { responseType: 'text' })
      .pipe(map(() => void 0));
  }

  public dryRunTask(
    taskId: string,
    request: InspectionDryRunRequest,
  ): Observable<InspectionDryRunResponse> {
    const url = BASE_URL + `api/v2/inspection/tasks/${taskId}/dryrun`;
    return this.http.post<InspectionDryRunResponse>(url, request);
  }

  public getInspectionData(taskId: string, reporter: DownloadProgressReporter) {
    const url = BASE_URL + `api/v2/inspection/tasks/${taskId}/data`;
    const httpRequest = new HttpRequest('GET', url, null, {
      reportProgress: true,
      responseType: 'blob',
    });
    return this.http.request<Blob>(httpRequest).pipe(
      tap((event) => {
        if (event.type === HttpEventType.DownloadProgress) {
          reporter(event.loaded, event.total ?? 0);
        }
      }),
      last(),
      map((event) => {
        if (event.type === HttpEventType.Response) {
          return event.body;
        } else {
          throw new Error('unreachable. Unexpected last event of http request');
        }
      }),
    );
  }

  public getPopup(): Observable<PopupFormRequest | null> {
    const url = BASE_URL + `api/v2/popup`;
    return this.http.get<PopupFormRequest | null>(url);
  }

  public validatePopupAnswer(
    answer: PopupAnswerResponse,
  ): Observable<PopupAnswerValidationResult> {
    const url = BASE_URL + `api/v2/popup/validate`;
    return this.http.post<PopupAnswerValidationResult>(url, answer);
  }
  public answerPopup(answer: PopupAnswerResponse): Observable<void> {
    const url = BASE_URL + `api/v2/popup/answer`;
    return this.http.post(url, answer).pipe(map(() => {}));
  }

  public cancelInspection(taskId: string) {
    const url = BASE_URL + `api/v2/inspection/tasks/${taskId}/cancel`;
    return this.http
      .post(url, null, { responseType: 'text' })
      .pipe(map(() => {}));
  }
}

export class InspectionTaskClient {
  private static DRYRUN_DEBOUNCE_DURATION = 100;

  public features = new ReplaySubject<InspectionFeature[]>(1);

  private dryRunParameter = new Subject<InspectionDryRunRequest>();

  private nonFormParameters = concat(this.viewState.timezoneShift).pipe(
    map((tzShift) => ({
      timezoneShift: tzShift,
    })),
    shareReplay(1),
  );

  public dryRunResult = this.dryRunParameter.pipe(
    debounceTime(InspectionTaskClient.DRYRUN_DEBOUNCE_DURATION),
    switchMap((param) => this.dryrunDirect(param)),
    shareReplay(1),
  );

  constructor(
    private readonly api: BackendAPI,
    public readonly taskId: string,
    private readonly viewState: ViewStateService,
  ) {
    this.downloadFeatureList();
  }

  public downloadFeatureList() {
    return this.api
      .getFeatureList(this.taskId)
      .pipe(map((r) => r.features))
      .subscribe((features) => this.features.next(features));
  }

  public setFeatures(featureIds: string[]) {
    return this.api
      .setEnabledFeatures(this.taskId, featureIds)
      .subscribe(() => {
        this.downloadFeatureList();
      });
  }

  public run(request: InspectionRunRequest) {
    return this.getRunParameter(request).pipe(
      switchMap((request) => {
        return this.api.runTask(this.taskId, request);
      }),
      map(() => {}),
    );
  }

  public dryrun(request: InspectionDryRunRequest) {
    this.dryRunParameter.next(request);
  }

  /**
   * dryrunDirect calls the dryrun API directly without debouncing.
   * This method is public for testing purpose. Use dryrun method instead.
   */
  public dryrunDirect(request: InspectionDryRunRequest) {
    return this.getRunParameter(request).pipe(
      switchMap((request) => this.api.dryRunTask(this.taskId, request)),
    );
  }

  private getRunParameter(
    request: InspectionRunRequest | InspectionDryRunRequest,
  ): Observable<{ [key: string]: unknown }> {
    return of(request).pipe(
      withLatestFrom(this.nonFormParameters),
      map(([request, nonForm]) => ({
        ...request,
        ...nonForm,
      })),
      tap((param) => {
        console.log(param);
      }),
    );
  }
}

/**
 * Utility functions using BackendAPI interface
 */
export class BackendAPIUtil {
  /**
   * Save the inspection data as a file
   */
  public static downloadInspectionDataAsFile(api: BackendAPI, taskId: string) {
    return combineLatest([
      api.getInspectionMetadata(taskId),
      api.getInspectionData(taskId, (done, all) => {
        console.log(`Downloading inspection data ${done}/${all}`);
      }),
    ]).pipe(
      map(([metadata, blob]) => {
        if (blob === null) return;
        const link = document.createElement('a');
        link.download = metadata.header.suggestedFilename;
        link.href = window.URL.createObjectURL(blob);
        link.style.display = 'none';
        document.body.appendChild(link);
        link.click();
        link.remove();
        return metadata.header.suggestedFilename;
      }),
    );
  }
}
