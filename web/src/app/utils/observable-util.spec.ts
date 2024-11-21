import { map, take, timer } from 'rxjs';
import { monitorElementHeight } from './observable-util';

describe('observable-util', () => {
  describe('monitorElementHeight', () => {
    it('emits heights on resize events', (done) => {
      const element = document.createElement('div');
      document.body.appendChild(element);
      element.style.height = '100px';

      const sizes = [200, 200, 100];
      const gotValues: number[] = [];
      const observable = monitorElementHeight(element);
      observable.subscribe((v) => {
        gotValues.push(v);
      });
      timer(0, 100)
        .pipe(
          take(3),
          map((v) => `${sizes[v]}px`),
        )
        .subscribe((size) => {
          element.style.height = size;
        });

      timer(1000).subscribe(() => {
        expect(gotValues).toEqual([100, 200, 100]);
        done();
      });
    });
  });
});
