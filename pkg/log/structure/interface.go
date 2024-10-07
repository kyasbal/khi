package structure

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

// ReaderDataAdapter convert a type to a structuredata.StructureData with storing the data in the StructureDataStore instanciated from StuctureDataFactory.
type ReaderDataAdapter interface {
	// GetReaderBackedByStore the given data into a structuredata.StructureData
	GetReaderBackedByStore(store structuredatastore.StructureDataStore) (*Reader, error)
}
