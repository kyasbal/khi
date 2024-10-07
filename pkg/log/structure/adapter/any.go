package adapter

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parser/yaml/yamlutil"
)

type AnyAdapter struct {
	objectRef any
}

// GetReaderBackedByStore implements structure.ReaderDataAdapter.
func (a *AnyAdapter) GetReaderBackedByStore(store structuredatastore.StructureDataStore) (*structure.Reader, error) {
	yamlString, err := yamlutil.MarshalToYamlString(a.objectRef)
	if err != nil {
		return nil, err
	}
	return Yaml(yamlString).GetReaderBackedByStore(store)
}

func Any(objectRef any) *AnyAdapter {
	return &AnyAdapter{objectRef: objectRef}
}

var _ structure.ReaderDataAdapter = (*AnyAdapter)(nil)
