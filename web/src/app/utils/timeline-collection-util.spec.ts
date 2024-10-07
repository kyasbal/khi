import { ParentRelationship } from '../generated';
import { TimelineEntry } from '../store/timeline';
import { SelectOnlyDeeperOrEqual } from './timeline-collection-util';

describe('TimelineCollectionUtility', () => {
  it('SelectOnlyDeeperOrEqual', () => {
    const timelines = [
      new TimelineEntry('1', [], [], ParentRelationship.RelationshipChild),
      new TimelineEntry('1#1', [], [], ParentRelationship.RelationshipChild),
      new TimelineEntry('1#2', [], [], ParentRelationship.RelationshipChild),
      new TimelineEntry('1#2#1', [], [], ParentRelationship.RelationshipChild),
      new TimelineEntry(
        '1#2#2',

        [],
        [],
        ParentRelationship.RelationshipChild,
      ),
      new TimelineEntry('2', [], [], ParentRelationship.RelationshipChild),
      new TimelineEntry('3', [], [], ParentRelationship.RelationshipChild),
    ];
    const result = SelectOnlyDeeperOrEqual(timelines, 2);
    expect(result.length).toBe(4);
    expect(result[0].resourcePath).toBe('1');
    expect(result[1].resourcePath).toBe('1#2');
    expect(result[2].resourcePath).toBe('1#2#1');
    expect(result[3].resourcePath).toBe('1#2#2');
  });
});
