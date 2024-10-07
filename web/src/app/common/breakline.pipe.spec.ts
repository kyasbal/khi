import { BreaklinePipe } from './breakline.pipe';

describe('BreaklinePipe', () => {
  it('create an instance', () => {
    const pipe = new BreaklinePipe();
    expect(pipe).toBeTruthy();
  });
  it('convert string not containing breakline', () => {
    const pipe = new BreaklinePipe();
    expect(pipe.transform('foo')).toBe('foo');
  });
  it('convert string containing breaklines', () => {
    const pipe = new BreaklinePipe();
    expect(pipe.transform('a\nb\nc')).toBe('a<br/>b<br/>c');
  });
});
