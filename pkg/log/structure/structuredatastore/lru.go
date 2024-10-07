package structuredatastore

import (
	"crypto/md5"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/cache"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"
)

type LRUStructureDataStoreRef struct {
	store    *LRUStructureDataStoreFactory
	cacheRef cache.Cache[structuredata.StructureData]
	hashKey  string
}

// GetStore implements StructureDataStorageRef.
func (l *LRUStructureDataStoreRef) GetStore() StructureDataStore {
	return l.store
}

// Get implements StructureDataStoreRef.
func (l *LRUStructureDataStoreRef) Get() (structuredata.StructureData, error) {
	return l.cacheRef.Get(l.hashKey)
}

var _ StructureDataStorageRef = (*LRUStructureDataStoreRef)(nil)

type LRUStructureDataStoreFactory struct {
	cache cache.Cache[structuredata.StructureData]
}

type LogEntryConverter struct{}

// Default implements cache.CacheItemConverter.
func (l *LogEntryConverter) Default() structuredata.StructureData {
	return nil
}

// Deserialize implements cache.CacheItemConverter.
func (l *LogEntryConverter) Deserialize(source []byte) (structuredata.StructureData, error) {
	yamlString := string(source)
	sd, err := structuredata.DataFromYaml(yamlString)
	if err != nil {
		return nil, err
	}
	return sd, nil
}

// Serialize implements cache.CacheItemConverter.
func (l *LogEntryConverter) Serialize(item structuredata.StructureData) ([]byte, error) {
	yamlStr, err := structuredata.ToYaml(item)
	if err != nil {
		return nil, err
	}
	return []byte(yamlStr), nil
}

var _ cache.CacheItemConverter[structuredata.StructureData] = (*LogEntryConverter)(nil)

// Convert implements cache.CacheItemConverter.]

func NewLRUStructureDataStoreFactory() *LRUStructureDataStoreFactory {
	return &LRUStructureDataStoreFactory{
		cache: cache.NewShardingCache(
			func() cache.CacheItemStorageProvider {
				return cache.NewGZipCacheItemStorageProvider(cache.NewShardingMapStorageProvider(common.NewSuffixShardingProvider(128, 4)))
			},
			&LogEntryConverter{},
			func() cache.CacheItemReleaseStrategy {
				return cache.NewLRUCacheItemReleaseStrategy(50)
			},
			common.NewSuffixShardingProvider(128, 4),
		),
	}
}

// StoreStructureData implements StructureDataStoreFactory.
func (l *LRUStructureDataStoreFactory) StoreStructureData(sd structuredata.StructureData) (StructureDataStorageRef, error) {
	yamlStr, err := structuredata.ToYaml(sd)
	if err != nil {
		return nil, err
	}
	md5 := md5.Sum([]byte(yamlStr))
	hash := fmt.Sprintf("%x", md5)
	l.cache.SetAsync(hash, []byte(yamlStr))
	return &LRUStructureDataStoreRef{
		store:    l,
		cacheRef: l.cache,
		hashKey:  hash,
	}, nil
}

var _ StructureDataStore = (*LRUStructureDataStoreFactory)(nil)
