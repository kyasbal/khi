import { $alignedGroup, $sized_rect } from '../builder-alias';
import { graphRootIt } from '../test/graph-test-utiility';
import { Direction } from './base-containers';
import { ArrowHead } from './element';

describe('Sized rect', () => {
  graphRootIt('Sized rect should retain the size', 200, (r) => {
    r.withChildren([
      $alignedGroup(Direction.Horizontal)
        .withGap(30)
        .withChildren([
          $sized_rect(100, 100).withId('small'),
          $sized_rect(100, 300).withId('big'),
        ]),
    ]).render();

    const small = r.find('small')?.element;
    const big = r.find('big')?.element;

    expect(small).not.toBeNull();
    expect(big).not.toBeNull();
    expect(small?.getAttribute('width')).toBe('100');
    expect(small?.getAttribute('height')).toBe('100');
    expect(small?.getAttribute('x')).toBe('0');
    expect(small?.getAttribute('y')).toBe('0');

    expect(big?.getAttribute('width')).toBe('100');
    expect(big?.getAttribute('height')).toBe('300');
    expect(big?.getAttribute('x')).toBe('130');
    expect(big?.getAttribute('y')).toBe('0');
  });
});

describe('Arrow head', () => {
  graphRootIt('Arrow head should be rendered correctly', 100, (r) => {
    r.withChildren([new ArrowHead(10, 0)]).render();
  });
});
