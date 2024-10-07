import { Pipe, PipeTransform } from '@angular/core';

/**
 * Convert given string to the valid class name in CSS.
 */
@Pipe({
  name: 'cssClassFormat',
})
export class CssClassFormatPipe implements PipeTransform {
  transform(value: string): string {
    value = value.toLowerCase();
    value = value.replace('(', ' ');
    value = value.replace(')', ' ');
    value = value.trim();
    return value.replace(' ', '-');
  }
}
