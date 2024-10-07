import { GraphObjectWithElementType } from './base-containers';

/**
 * Graph wrapper that can be placed just under <defs>
 */
export class GraphDefsChildItem<
  T extends SVGElement,
> extends GraphObjectWithElementType<T> {}

/**
 * Graph wrapper for <pattern>
 */
export class GraphPattern extends GraphDefsChildItem<SVGPatternElement> {
  constructor(width: number, height: number) {
    super(document.createElementNS('http://www.w3.org/2000/svg', 'pattern'));
    this.transform.onSvgElementUpdate = (pos, size) => {
      this.withMinSize(size.width, size.height);
      return;
    };
    this.withMinSize(width, height);
    this.withStyle({
      width: `${width}px`,
      height: `${height}px`,
      viewBox: `0,0,${width},${height}`,
      patternUnits: 'userSpaceOnUse',
    });
  }
}
