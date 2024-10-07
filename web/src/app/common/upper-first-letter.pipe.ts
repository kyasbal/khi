import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'upperFirstLetter',
})
export class UpperFirstLetterPipe implements PipeTransform {
  transform(value: string): string | undefined {
    if (!value) return '';
    return value[0].toUpperCase() + value.slice(1);
  }
}
