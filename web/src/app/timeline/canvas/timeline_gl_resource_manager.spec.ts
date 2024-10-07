import { LRULifetimeManager } from './timeline_gl_resource_manager';

describe('LRULifetimeManager', () => {
  it("returns null when it's not hitting the capacity", () => {
    const lru = new LRULifetimeManager<number>(3);
    expect(lru.touch(1)).toBe(null);
    expect(lru.touch(2)).toBe(null);
    expect(lru.touch(3)).toBe(null);
  });

  it("returns null when it's not hitting the capacity and ignore the same element", () => {
    const lru = new LRULifetimeManager<number>(1);
    expect(lru.touch(1)).toBe(null);
    expect(lru.touch(1)).toBe(null);
    expect(lru.touch(1)).toBe(null);
  });

  it("returns the element to remove when it's hitting the capacity", () => {
    const lru = new LRULifetimeManager<number>(2);
    expect(lru.touch(1)).toBe(null);
    expect(lru.touch(1)).toBe(null);
    expect(lru.touch(2)).toBe(null);
    expect(lru.touch(3)).toBe(1);
  });
});
