/**
 * Generate SHA-512 hash string from given ArrayBuffer
 * @param source source of the hash
 * @returns SHA-512 hash in hex-string
 */
export async function sha512FromArrayBuffer(
  source: ArrayBuffer,
): Promise<string> {
  const hashBuffer = await crypto.subtle.digest('SHA-512', source);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map((b) => b.toString(16).padStart(2, '0')).join('');
}
