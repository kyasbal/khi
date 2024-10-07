import { AbstractControl } from '@angular/forms';
import { RegexValidator } from './regex-validator';

describe('RegexValidator', () => {
  it('Should accept the valid regex', () => {
    const input = '\\s*[0-9a-z]*';
    const actual = RegexValidator()({ value: input } as AbstractControl);
    expect(actual).toBeNull();
  });
  it("Shouldn't accept the invalid regex", () => {
    const input = '(\\s';
    const actual = RegexValidator()({ value: input } as AbstractControl);
    expect(actual).not.toBeNull();
  });
});
