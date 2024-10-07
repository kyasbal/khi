package adapter

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/merger"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

// MergeYamlAdapter implements ReaderDataAdapter to get the Reader of a merged Yaml from 2 yaml strings.
type MergeYamlAdapter struct {
	prevYaml            string
	currentYaml         string
	mergeConfigResolver *merger.MergeConfigResolver
}

func MergeYaml(prevYaml string, currentYaml string, mergeConfigResolver *merger.MergeConfigResolver) *MergeYamlAdapter {
	return &MergeYamlAdapter{
		prevYaml:            prevYaml,
		currentYaml:         currentYaml,
		mergeConfigResolver: mergeConfigResolver,
	}
}

// GetReaderBackedByStore implements structure.ReaderDataAdapter.
func (y *MergeYamlAdapter) GetReaderBackedByStore(store structuredatastore.StructureDataStore) (*structure.Reader, error) {
	prevStructureData, err := structuredata.DataFromYaml(y.prevYaml)
	if err != nil {
		return nil, err
	}
	currentStructureData, err := structuredata.DataFromYaml(y.currentYaml)
	if err != nil {
		return nil, err
	}
	merged := merger.NewStrategicMergedStructureData("", prevStructureData, currentStructureData, y.mergeConfigResolver)
	storeRef, err := store.StoreStructureData(merged)
	if err != nil {
		return nil, err
	}
	return structure.NewReader(storeRef), nil
}

var _ structure.ReaderDataAdapter = (*MergeYamlAdapter)(nil)
