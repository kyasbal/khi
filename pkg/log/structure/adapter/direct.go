package adapter

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

// DirectAdapter implements ReaderDataAdapter to pass StructureData directly into the Reader.
type DirectAdapter struct {
	sd structuredata.StructureData
}

func Direct(sd structuredata.StructureData) *DirectAdapter {
	return &DirectAdapter{sd: sd}
}

// GetReaderBackedByStore implements structure.ReaderDataAdapter.
func (d *DirectAdapter) GetReaderBackedByStore(store structuredatastore.StructureDataStore) (*structure.Reader, error) {
	sd, err := store.StoreStructureData(d.sd)
	if err != nil {
		return nil, err
	}
	return structure.NewReader(sd), nil
}

var _ structure.ReaderDataAdapter = (*DirectAdapter)(nil)
