package structure

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"

// ReaderFactory instanciate the instance of Reader.
type ReaderFactory struct {
	store structuredatastore.StructureDataStore
}

func NewReaderFactory(store structuredatastore.StructureDataStore) *ReaderFactory {
	return &ReaderFactory{
		store: store,
	}
}

func (f *ReaderFactory) NewReader(adapter ReaderDataAdapter) (*Reader, error) {
	return adapter.GetReaderBackedByStore(f.store)
}
