package structuredatastore

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"
)

func TestLRUStructureDataStoreFactory(t *testing.T) {
	ITEM_COUNT := 10000
	lru := NewLRUStructureDataStoreFactory()
	for i := 0; i < ITEM_COUNT; i++ {
		data := fmt.Sprintf("textPayload: hello-%d\n", i)
		sd, err := structuredata.DataFromYaml(data)
		if err != nil {
			t.Fatal(err)
		}
		d, err := lru.StoreStructureData(sd)
		if err != nil {
			t.Errorf(err.Error())
		}
		_, err = d.Get()
		if err != nil {
			t.Errorf(err.Error())
		}
		// Needs to check Get() call twice to verify it's on the cache not to read from the storage
		sd, err = d.Get()
		if err != nil {
			t.Errorf(err.Error())
		}
		yaml, err := structuredata.ToYaml(sd)
		if err != nil {
			t.Errorf(err.Error())
		}
		if yaml != data {
			t.Errorf("expected %s, got %s", data, yaml)
		}
	}
}
