package query

import (
	"testing"

	metadata_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/metadata"
	"github.com/google/go-cmp/cmp"
)

func TestQueryConformance(t *testing.T) {
	metadata_test.ConformanceMetadataTypeTest(t, &QueryMetadata{
		Queries: []*QueryItem{
			{
				Id:    "foo",
				Query: "foo-body",
			},
			{
				Id:    "bar",
				Query: "bar-body",
			},
		},
	})
}

func TestQuerySerializeInSortedOrder(t *testing.T) {
	query := QueryMetadata{
		Queries: []*QueryItem{
			{Id: "a"},
			{Id: "c"},
			{Id: "b"},
			{Id: "e"},
			{Id: "d"},
		},
	}

	expected := []*QueryItem{
		{Id: "a"},
		{Id: "b"},
		{Id: "c"},
		{Id: "d"},
		{Id: "e"},
	}
	if diff := cmp.Diff(query.ToSerializable(), expected); diff != "" {
		t.Errorf("Query info serialization result was not in the sorted order\n%s", diff)
	}
}
