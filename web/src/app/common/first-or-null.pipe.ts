import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'firstOrUndefined',
})
export class FirstOrUndefined implements PipeTransform {
  transform<T>(value: T[] | null): T | undefined {
    if (value == null || value.length == 0) {
      return undefined;
    }
    return value[0];
  }
}
