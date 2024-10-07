import {
  Component,
  EnvironmentInjector,
  Inject,
  Input,
  inject,
} from '@angular/core';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';
import { Subject, map, shareReplay, startWith, withLatestFrom } from 'rxjs';
import {
  LOG_ANNOTATOR_RESOLVER,
  LogAnnotatorResolver,
} from '../annotator/log/resolver';

@Component({
  selector: 'khi-log-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.sass'],
})
export class LogHeaderComponent {
  private readonly envInjector = inject(EnvironmentInjector);

  @Input()
  public set logIndex(index: number) {
    this.logIndexObservable.next(index);
  }

  private logIndexObservable = new Subject<number>();

  public logEntryObservable = this.logIndexObservable.pipe(
    startWith(0),
    withLatestFrom(this._inspectionDataStore.$allLogs),
    map(([i, all]) => all[i]),
    shareReplay(1),
  );

  public logAnnotators = this.logAnnotatorResolver.getResolvedAnnotators(
    this.logEntryObservable,
    this.envInjector,
  );

  constructor(
    @Inject(LOG_ANNOTATOR_RESOLVER)
    private logAnnotatorResolver: LogAnnotatorResolver,
    private _inspectionDataStore: InspectionDataStoreService,
  ) {}
}
