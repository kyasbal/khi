import { CommonModule } from '@angular/common';
import { Component, Input, inject } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import {
  NEVER,
  Observable,
  filter,
  map,
  of,
  shareReplay,
  withLatestFrom,
} from 'rxjs';
import {
  KHIFileTextReference,
  LogAnnotationTypeResourceRef,
} from 'src/app/common/schema/khi-file-types';
import { InspectionDataStoreService } from 'src/app/services/inspection-data-store.service';
import { SelectionManagerService } from 'src/app/services/selection-manager.service';
import { AnnotationDecider, DECISION_HIDDEN } from '../annotator';
import { LogEntry } from 'src/app/store/log';

interface ResourceRefAnnotationViewModel {
  label: string;
  path: string;
}

@Component({
  standalone: true,
  templateUrl: './relationship-annotator.component.html',
  styleUrl: './relationship-annotator.component.sass',
  imports: [CommonModule, MatIconModule],
})
export class RelationshipAnnotatorComponent {
  private readonly selectionManager = inject(SelectionManagerService);

  @Input()
  refs: Observable<ResourceRefAnnotationViewModel[]> = NEVER;

  currentSelectedTimelinePath = this.selectionManager.selectedTimeline.pipe(
    map((t) => t?.resourcePath ?? ''),
    shareReplay(1),
  );

  public selectResource(resourcePath: string) {
    this.selectionManager.onSelectTimeline(resourcePath);
  }

  public highlightResource(resourcePath: string) {
    this.selectionManager.onHighlightTimeline(resourcePath);
  }

  public static inputMapper: AnnotationDecider<LogEntry> = (
    l?: LogEntry | null,
  ) => {
    if (!l) {
      return DECISION_HIDDEN;
    }
    const dataStore = inject(InspectionDataStoreService);
    const pathReferences: KHIFileTextReference[] = [];
    for (const annotation of l.annotations) {
      if (annotation.type == LogAnnotationTypeResourceRef) {
        const pathReference = annotation['path'] as KHIFileTextReference;
        pathReferences.push(pathReference);
      }
    }
    if (pathReferences.length == 0) return DECISION_HIDDEN;
    return {
      inputs: {
        refs: of(pathReferences).pipe(
          withLatestFrom(dataStore.textBufferSource.pipe(filter((tb) => !!tb))),
          map(([refs, bufferLoader]) =>
            [...new Set(refs.map((ref) => bufferLoader!.getText(ref)))].map(
              (refPath) => {
                const splittedPath = refPath.split('#');
                const resourceRefLabel = `${splittedPath[splittedPath.length - 1]} of ${splittedPath[splittedPath.length - 2]}`;
                return {
                  label: resourceRefLabel,
                  path: refPath,
                } as ResourceRefAnnotationViewModel;
              },
            ),
          ),
        ),
      },
    };
  };
}
