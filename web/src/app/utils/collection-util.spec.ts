import {
  filteElementsByIncludedSubstring,
  subtractSet,
} from './collection-util';

describe('Collection util', () => {
  it('subtractSet', () => {
    const subtractFrom = new Set(['a', 'b', 'c', 'd']);
    const subtracting = new Set(['a', 'c', 'e']);
    const result = subtractSet(subtractFrom, subtracting);
    expect(result.has('a')).toBeFalse();
    expect(result.has('b')).toBeTrue();
    expect(result.has('c')).toBeFalse();
    expect(result.has('d')).toBeTrue();
    expect(result.has('e')).toBeFalse();

    const subtractFrom2 = new Set(['a']);
    const subtracting2 = new Set([]);
    const result2 = subtractSet(subtractFrom2, subtracting2);
    expect(result2.has('a')).toBeTrue();
  });

  it('filteElementsByIncludedSubstring', () => {
    const result = filteElementsByIncludedSubstring(
      [
        'keyword',
        'keywordB',
        'not-key-word',
        'Akeyword',
        'keywordA',
        'Bkeyword',
        'yet-another-not-key-word',
      ],
      'keyword',
    );
    expect(result.length).toBe(5);
    expect(result[0]).toBe('keyword');
    expect(result[1]).toBe('keywordA');
    expect(result[2]).toBe('keywordB');
    expect(result[3]).toBe('Akeyword');
    expect(result[4]).toBe('Bkeyword');
  });
});
