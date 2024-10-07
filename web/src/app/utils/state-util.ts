/**
 * Returns the default value when the passed value was empty or null.
 * @param value
 * @param defaultValue
 * @returns
 */
export function nonEmptyOrDefaultString(
  value: string,
  defaultValue: string,
): string {
  if (value === undefined || value == null || value.trim() == '') {
    return defaultValue;
  }
  return value;
}
