package structuredatastore

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"

type OnMemoryStructureDataStoreRef struct {
	store *OnMemoryStructureDataStore
	data  structuredata.StructureData
}

// GetStore implements StructureDataStore.
func (d *OnMemoryStructureDataStoreRef) GetStore() StructureDataStore {
	return d.store
}

// Get implements StructureDataStore.
func (d *OnMemoryStructureDataStoreRef) Get() (structuredata.StructureData, error) {
	return d.data, nil
}

var _ StructureDataStorageRef = (*OnMemoryStructureDataStoreRef)(nil)

type OnMemoryStructureDataStore struct{}

// Instantiate implements StructureDataStore.
func (d *OnMemoryStructureDataStore) StoreStructureData(sd structuredata.StructureData) (StructureDataStorageRef, error) {
	return &OnMemoryStructureDataStoreRef{
		store: d,
		data:  sd,
	}, nil
}

var _ StructureDataStore = (*OnMemoryStructureDataStore)(nil)
