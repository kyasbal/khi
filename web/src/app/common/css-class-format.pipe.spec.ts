import { CssClassFormatPipe } from './css-class-format.pipe';

describe('CssClassFormatPipe', () => {
  it('create an instance', () => {
    const pipe = new CssClassFormatPipe();
    expect(pipe).toBeTruthy();
  });
  it('convert upper cased strings to lower case', () => {
    const pipe = new CssClassFormatPipe();
    expect(pipe.transform('FoO')).toBe('foo');
  });
  it('escapes removes invalid character in the end replace it with - in the middle', () => {
    const pipe = new CssClassFormatPipe();
    expect(pipe.transform('foo(bar)')).toBe('foo-bar');
  });
});
