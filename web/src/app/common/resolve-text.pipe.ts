import { Pipe, PipeTransform } from '@angular/core';
import { InspectionDataStoreService } from '../services/inspection-data-store.service';
import { Observable, map } from 'rxjs';
import { KHIFileTextReference } from './schema/khi-file-types';

/**
 * A pipe to resolve KHIFileTextReference type with data store.
 * Large texts shouldn't be kept in view models, resolve text from buffer source in data store.
 */
@Pipe({
  name: 'resolveText',
})
export class ResolveTextPipe implements PipeTransform {
  constructor(private dataStore: InspectionDataStoreService) {}
  transform(value: KHIFileTextReference): Observable<string> {
    return this.dataStore.textBufferSource.pipe(
      map((bs) => bs?.getText(value) ?? 'error'),
    );
  }
}
