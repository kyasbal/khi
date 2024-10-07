import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

export function RegexValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    try {
      new RegExp(control.value);
    } catch (e) {
      return { regex: e };
    }
    return null;
  };
}
