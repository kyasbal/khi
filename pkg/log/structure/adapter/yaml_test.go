package adapter

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

func TestYamlAdapter(t *testing.T) {
	store := structuredatastore.OnMemoryStructureDataStore{}
	yamlAdapter := Yaml("textPayload: hello world")
	reader, err := yamlAdapter.GetReaderBackedByStore(&store)
	if err != nil {
		t.Errorf(err.Error())
	}
	if reader.ReadStringOrDefault("textPayload", "") != "hello world" {
		t.Errorf("expected hello world, got %s", reader.ReadStringOrDefault("textPayload", ""))
	}
}
