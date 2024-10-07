import { TimelineEntry } from '../store/timeline';

/**
 * Filter non used upper layers when there were no any layer lower or equal to given depth
 * @param timelines
 * @param depth
 * @returns
 */
export function SelectOnlyDeeperOrEqual(
  timelines: TimelineEntry[],
  depth: number,
): TimelineEntry[] {
  const result: TimelineEntry[] = [];
  const retainIndicies: Set<number> = new Set();
  for (let i = 0; i < timelines.length; i++) {
    const timeline = timelines[i];
    if (timeline.layer >= depth) {
      let prevLayer = timeline.layer;
      for (let j = 0; j <= i && prevLayer >= 0; j++) {
        if (timelines[i - j].layer == prevLayer) {
          if (retainIndicies.has(i - j)) break;
          retainIndicies.add(i - j);
          prevLayer -= 1;
        }
      }
    }
  }
  for (let i = 0; i < timelines.length; i++) {
    if (retainIndicies.has(i)) {
      result.push(timelines[i]);
    }
  }
  return result;
}
