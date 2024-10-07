import { BehaviorSubject, ReplaySubject, Subject, filter } from 'rxjs';

/**
 * An abstract data storage got from inter-frame connection.
 */
export abstract class InterframeDatasource<T> {
  /**
   * Determine wheather the data update request from main frame should be accepted or not.
   */
  public readonly bound$ = new BehaviorSubject(true);

  /**
   * Must be updated by the child class.
   * The updated data sent from main frame should be routed to here.
   */
  protected readonly rawUpdateRequest$ = new Subject<T>();

  /**
   * An observable to monitor the data change.
   * Data update won't be reported when bound$ is false.
   */
  public readonly data$ = new ReplaySubject<T>(1);

  constructor() {
    this.rawUpdateRequest$
      .pipe(filter(() => this.bound$.value))
      .subscribe((data) => {
        this.data$.next(data);
      });
  }

  abstract enable(): void;

  abstract disable(): void;
}
