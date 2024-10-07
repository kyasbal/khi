package structuredatastore

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"

// structure.Reader read its data from StructureDataStorageRef.
// Implementations may store its data in storage or compressed on its demand
type StructureDataStorageRef interface {
	// Get the current structuredata.StructureData
	// This may read the data from storage or compressed memory by its implementation
	Get() (structuredata.StructureData, error)

	// GetStore returns the reference to the store hold the actual data of this reference.
	GetStore() StructureDataStore
}

// StructureDataStore is a factory instanciating StructureDataStore
type StructureDataStore interface {
	// Store the given StructureData and return the reference to it.
	StoreStructureData(sd structuredata.StructureData) (StructureDataStorageRef, error)
}
