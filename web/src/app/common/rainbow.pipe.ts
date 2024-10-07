import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'rainbow',
})
export class RainbowPipe implements PipeTransform {
  transform(
    value: number,
    baseHue: number,
    saturation: string,
    brightness: string,
  ): string {
    const PRIME_MULTIPLIER = 103;
    const hue = (baseHue + value * PRIME_MULTIPLIER) % 360;
    return `hsl(${hue}deg,${saturation},${brightness})`;
  }
}
