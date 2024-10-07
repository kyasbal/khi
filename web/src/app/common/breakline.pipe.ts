import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'breakline',
})
export class BreaklinePipe implements PipeTransform {
  transform(value: string): string {
    return value.replace(/\n/g, `<br/>`);
  }
}
