package common

// A thread-safe counter data structure
type ConcurrentCounter struct {
	counts *ShardingMap[int]
}

func NewConcurrentCounter(shardingMap *ShardingMap[int]) *ConcurrentCounter {
	return &ConcurrentCounter{
		counts: shardingMap,
	}
}

func NewDefaultConcurrentCounter(shardingProvider MapShardingProvider) *ConcurrentCounter {
	return NewConcurrentCounter(NewShardingMap[int](shardingProvider))
}

func (c *ConcurrentCounter) Get(key string) int {
	shard := c.counts.AcquireShardReadonly(key)
	defer c.counts.ReleaseShardReadonly(key)
	if count, found := shard[key]; found {
		return count
	} else {
		return 0
	}
}

func (c *ConcurrentCounter) Incr(key string) int {
	shard := c.counts.AcquireShard(key)
	defer c.counts.ReleaseShard(key)
	if count, found := shard[key]; found {
		shard[key] = count + 1
	} else {
		shard[key] = 1
	}
	return shard[key]
}
