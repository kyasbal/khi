import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'parsePath',
})
export class ParsePathPipe implements PipeTransform {
  transform(value: string, index: number): string {
    return value.split('/')[index] ?? '';
  }
}
