package adapter

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

// YamlAdapter implements ReaderDataAdapter to get the Reader from YAML string.
type YamlAdapter struct {
	sourceYaml string
}

// Returns adapter for parsing source yaml.
func Yaml(sourceYaml string) *YamlAdapter {
	return &YamlAdapter{
		sourceYaml: sourceYaml,
	}
}

// GetReaderBackedByStore implements StructureDataAdapter.
func (y *YamlAdapter) GetReaderBackedByStore(store structuredatastore.StructureDataStore) (*structure.Reader, error) {
	sd, err := structuredata.DataFromYaml(y.sourceYaml)
	if err != nil {
		return nil, err
	}
	sdstore, err := store.StoreStructureData(sd)
	if err != nil {
		return nil, err
	}
	return structure.NewReader(sdstore), nil
}

var _ structure.ReaderDataAdapter = (*YamlAdapter)(nil)
