import { UpperFirstLetterPipe } from './upper-first-letter.pipe';

describe('UpperFirstLetterPipe', () => {
  it('create an instance', () => {
    const pipe = new UpperFirstLetterPipe();
    expect(pipe).toBeTruthy();
  });
  it('convert the first letter of given string to upper case', () => {
    const pipe = new UpperFirstLetterPipe();
    expect(pipe.transform('foo')).toBe('Foo');
  });
  it('keeps the first letter of given string to upper case when it was uppercase at first', () => {
    const pipe = new UpperFirstLetterPipe();
    expect(pipe.transform('Foo')).toBe('Foo');
  });
  it('returns the empty string when the empty string was given', () => {
    const pipe = new UpperFirstLetterPipe();
    expect(pipe.transform('')).toBe('');
  });
});
