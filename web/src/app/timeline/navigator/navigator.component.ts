import { Component, EnvironmentInjector, inject } from '@angular/core';
import { Observable, map, of, shareReplay } from 'rxjs';
import { ResolvedAnnotator } from 'src/app/annotator/annotator';
import { NAVIGATOR_ANNOTATOR_RESOLVER } from 'src/app/annotator/navigator/resolver';
import { SelectionManagerService } from 'src/app/services/selection-manager.service';
import { TimelineEntry } from 'src/app/store/timeline';

interface NavigatorLayer {
  label: string;
  icon: string;
  isLast: boolean;
  annotators: Observable<ResolvedAnnotator[]>;
}

/**
 * NavigatorComponent is a control shown bottom left of timeline chart.
 * Contains information for currently selected timeline.
 */
@Component({
  templateUrl: './navigator.component.html',
  styleUrl: './navigator.component.sass',
  selector: 'khi-timeline-navigator',
})
export class NavigatorComponent {
  private readonly envInjector = inject(EnvironmentInjector);
  private readonly selectionManager = inject(SelectionManagerService);
  private readonly navigatorAnnotatorResolver = inject(
    NAVIGATOR_ANNOTATOR_RESOLVER,
  );
  selectedTimeline = this.selectionManager.selectedTimeline;

  layerTimelines = this.selectedTimeline.pipe(
    map((tl) => {
      const layers: TimelineEntry[] = [];
      while (tl) {
        layers.push(tl);
        tl = tl.parent;
      }
      return layers.reverse();
    }),
    shareReplay(1),
  );

  layers = this.layerTimelines.pipe(
    map((tls) =>
      tls.map(
        (tl, index) =>
          ({
            label: tl.name,
            icon: '',
            isLast: index == tls.length - 1,
            annotators: this.navigatorAnnotatorResolver.getResolvedAnnotators(
              of(tl),
              this.envInjector,
            ),
          }) as NavigatorLayer,
      ),
    ),
  );
}
