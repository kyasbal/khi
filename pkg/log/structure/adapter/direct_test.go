package adapter

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

func TestDirectAdapter(t *testing.T) {
	store := structuredatastore.OnMemoryStructureDataStore{}
	sd, err := structuredata.DataFromYaml("textPayload: hello world")
	if err != nil {
		t.Errorf(err.Error())
	}
	direct := Direct(sd)
	reader, err := direct.GetReaderBackedByStore(&store)
	if err != nil {
		t.Errorf(err.Error())
	}
	if reader.ReadStringOrDefault("textPayload", "") != "hello world" {
		t.Errorf("expected hello world, got %s", reader.ReadStringOrDefault("textPayload", ""))
	}
}
