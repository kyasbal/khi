package adapter

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

func TestAnyAdapter(t *testing.T) {
	store := structuredatastore.OnMemoryStructureDataStore{}
	direct := Any(map[string]string{
		"textPayload": "hello world",
	})
	reader, err := direct.GetReaderBackedByStore(&store)
	if err != nil {
		t.Errorf(err.Error())
	}
	if reader.ReadStringOrDefault("textPayload", "") != "hello world" {
		t.Errorf("expected hello world, got %s", reader.ReadStringOrDefault("textPayload", ""))
	}
}
