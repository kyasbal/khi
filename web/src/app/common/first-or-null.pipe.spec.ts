import { FirstOrUndefined } from './first-or-null.pipe';

describe('FirstOrNullPipe', () => {
  it('create an instance', () => {
    const pipe = new FirstOrUndefined();
    expect(pipe).toBeTruthy();
  });
  it('returns null when null was passed', () => {
    const pipe = new FirstOrUndefined();
    expect(pipe.transform(null)).toBeUndefined();
  });
  it('returns null when empty array was passed', () => {
    const pipe = new FirstOrUndefined();
    expect(pipe.transform([])).toBeUndefined();
  });
  it('returns first element', () => {
    const pipe = new FirstOrUndefined();
    expect(pipe.transform([1, 2, 3])).toBe(1);
  });
});
