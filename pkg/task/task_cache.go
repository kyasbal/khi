package task

import "sync"

type TaskVariableCache interface {
	Store(key string, value any)
	Load(key string) (any, bool)
	// this must be an atomic operation. When the key is missing, store the value and return it.
	// This should behave same with sync.Map.LoadOrStore
	LoadOrStore(key string, value any) (any, bool)
}

// LocalTaskVariableCache implements TaskVariableCache and it stores the data inside of the struct.
// The cache is not shared with the other task graph. This is used for testing.
type LocalTaskVariableCache struct {
	cache sync.Map
}

// Load implements TaskVariableCache.
func (ltc *LocalTaskVariableCache) Load(key string) (any, bool) {
	return ltc.cache.Load(key)
}

// Store implements TaskVariableCache.
func (ltc *LocalTaskVariableCache) Store(key string, value any) {
	ltc.cache.Swap(key, value)
}

func (ltc *LocalTaskVariableCache) LoadOrStore(key string, value any) (any, bool) {
	return ltc.cache.LoadOrStore(key, value)
}

var _ TaskVariableCache = (*LocalTaskVariableCache)(nil)

func NewLocalTaskVariableCache() *LocalTaskVariableCache {
	return &LocalTaskVariableCache{
		cache: sync.Map{},
	}
}

var globalTaskVariableCacheMemory sync.Map = sync.Map{}

type GlobalTaskVariableCache struct {
}

// LoadOrStore implements TaskVariableCache.
func (*GlobalTaskVariableCache) LoadOrStore(key string, value any) (any, bool) {
	return globalTaskVariableCacheMemory.LoadOrStore(key, value)
}

// Load implements TaskVariableCache.
func (*GlobalTaskVariableCache) Load(key string) (any, bool) {
	return globalTaskVariableCacheMemory.Load(key)
}

// Store implements TaskVariableCache.
func (*GlobalTaskVariableCache) Store(key string, value any) {
	globalTaskVariableCacheMemory.Swap(key, value)
}

var _ TaskVariableCache = (*GlobalTaskVariableCache)(nil)
