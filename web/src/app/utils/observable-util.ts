import {
  BehaviorSubject,
  Observable,
  animationFrames,
  endWith,
  map,
  takeWhile,
} from 'rxjs';

/**
 *
 * @deprecated prefer not to store the actual value with BehaviorSubject. Use pipe and shareReplay instead to support the late registered observable.
 */
export function asBehaviorSubject<T>(
  ob: Observable<T>,
  initialValue: T,
): BehaviorSubject<T> {
  const bs = new BehaviorSubject<T>(initialValue);
  ob.subscribe(bs);
  return bs;
}

/**
 * Returns an observable emitting number of tweening between start to end with duration milliseconds.
 */
export function tweenNumber(start: number, end: number, duration: number) {
  const diff = end - start;
  return animationFrames().pipe(
    map(({ elapsed }) => elapsed / duration),
    takeWhile((v) => v < 1),
    endWith(1),
    map((v) => v * diff + start),
  );
}
