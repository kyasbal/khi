export function conditionalModule<T>(
  condition: boolean,
  ...extensions: T[]
): T[] {
  return condition ? extensions : [];
}
