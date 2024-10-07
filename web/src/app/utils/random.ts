/**
 * Returns random string with 11 characters.
 * @returns generated random string
 */
export function randomString(): string {
  return Math.random().toString(36).substring(2);
}
