package adapter

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/merger"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
)

func TestYamlMergeAdapterTest(t *testing.T) {
	store := structuredatastore.OnMemoryStructureDataStore{}
	yamlAdapter := MergeYaml("foo: hello", "bar: world", &merger.MergeConfigResolver{})
	reader, err := yamlAdapter.GetReaderBackedByStore(&store)
	if err != nil {
		t.Errorf(err.Error())
	}
	if reader.ReadStringOrDefault("foo", "") != "hello" {
		t.Errorf("expected hello world, got %s", reader.ReadStringOrDefault("foo", ""))
	}
	if reader.ReadStringOrDefault("bar", "") != "world" {
		t.Errorf("expected world, got %s", reader.ReadStringOrDefault("bar", ""))
	}
}
