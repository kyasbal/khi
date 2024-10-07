/**
 * Calculate unique number from selectedLog and highlightedLogIndices. This is used for deciding if interaction buffers are needed to be updated or not.
 * @param selectedLogIndex
 * @param highlightedLogIndices
 * @returns
 */
export function calcInteractionDigest(
  selectedLogIndex: number,
  highlightedLogIndices: Set<number>,
): number {
  let hash = 0;
  for (const index of highlightedLogIndices) {
    hash ^= index;
  }
  hash ^= ~selectedLogIndex;
  return hash;
}
