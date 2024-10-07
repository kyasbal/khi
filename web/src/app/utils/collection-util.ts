export function iterToArr<T>(iter: IterableIterator<T>): T[] {
  const result: T[] = [];
  for (const elem of iter) {
    result.push(elem);
  }
  return result;
}

export function subtractSet<T>(
  baseSet: Set<T>,
  subtractingSet: Set<T>,
): Set<T> {
  const result = new Set<T>();
  for (const elem of baseSet.values()) {
    if (!subtractingSet.has(elem)) {
      result.add(elem);
    }
  }
  return result;
}

export function filteElementsByIncludedSubstring(
  candidates: Iterable<string>,
  query: string,
): string[] {
  const middleResult: { value: string; index: number }[] = [];
  for (const candidate of candidates) {
    const index = candidate.indexOf(query);
    if (index != -1) {
      middleResult.push({
        value: candidate,
        index,
      });
    }
  }
  return middleResult
    .sort((a, b) => {
      const diff = a.index - b.index;
      return diff == 0 ? a.value.localeCompare(b.value) : diff;
    })
    .map((a) => a.value);
}
