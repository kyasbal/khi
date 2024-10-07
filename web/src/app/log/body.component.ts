import { Component, EnvironmentInjector, Input, inject } from '@angular/core';
import { LOG_TOOL_ANNOTATOR_RESOLVER } from '../annotator/log-tool/resolver';
import { Subject, map, shareReplay, startWith, withLatestFrom } from 'rxjs';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';

@Component({
  selector: 'khi-log-body',
  templateUrl: './body.component.html',
  styleUrls: ['./body.component.sass'],
})
export class LogBodyComponent {
  private readonly dataStore = inject(InspectionDataStoreService);

  private readonly envInjector = inject(EnvironmentInjector);

  private readonly logToolAnnotatorResolver = inject(
    LOG_TOOL_ANNOTATOR_RESOLVER,
  );

  @Input()
  public set logIndex(index: number) {
    this.logIndexObservable.next(index);
  }

  private logIndexObservable = new Subject<number>();

  public logEntryObservable = this.logIndexObservable.pipe(
    startWith(0),
    withLatestFrom(this.dataStore.$allLogs),
    map(([i, all]) => all[i]),
    shareReplay(1),
  );

  public logAnnotators = this.logToolAnnotatorResolver.getResolvedAnnotators(
    this.logEntryObservable,
    this.envInjector,
  );
}
