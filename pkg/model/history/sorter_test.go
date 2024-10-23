package history

import (
	"slices"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/google/go-cmp/cmp"
)

type TestResourceChunkSortStrategy struct {
}

// SortChunk implements ResourceChunkSortStrategy.
func (t *TestResourceChunkSortStrategy) SortChunk(builder *Builder, parents []*Resource, groupedRelationship enum.ParentRelationship, chunk []*Resource) ([]*Resource, error) {
	slices.SortFunc(
		chunk, func(a, b *Resource) int {
			return strings.Compare(a.ResourceName, b.ResourceName)
		},
	)
	return chunk, nil
}

var _ ResourceChunkSortStrategy = (*TestResourceChunkSortStrategy)(nil)

type testAllSkipChunkSortStrategy struct {
}

// SortChunk implements ResourceChunkSortStrategy.
func (t *testAllSkipChunkSortStrategy) SortChunk(builder *Builder, parents []*Resource, groupedRelationship enum.ParentRelationship, chunk []*Resource) ([]*Resource, error) {
	return nil, ErrorSortSkipped
}

var _ ResourceChunkSortStrategy = (*testAllSkipChunkSortStrategy)(nil)

func newResourceForTesting(name string, rel enum.ParentRelationship, children ...*Resource) *Resource {
	if children == nil {
		children = make([]*Resource, 0)
	}
	return &Resource{
		ResourceName: name,
		Relationship: rel,
		Children:     children,
	}
}

func TestResourceSort(t *testing.T) {
	type testCase struct {
		name     string
		input    []*Resource
		expected []*Resource
		strategy []ResourceChunkSortStrategy
	}
	testCases := []testCase{
		{
			name: "single layer with name",
			input: []*Resource{
				newResourceForTesting("c", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipChild),
				newResourceForTesting("a", enum.RelationshipChild),
			},
			expected: []*Resource{
				newResourceForTesting("a", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipChild),
				newResourceForTesting("c", enum.RelationshipChild),
			},
			strategy: []ResourceChunkSortStrategy{&TestResourceChunkSortStrategy{}},
		},
		{
			name: "single layer with name and skipped sorter",
			input: []*Resource{
				newResourceForTesting("c", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipChild),
				newResourceForTesting("a", enum.RelationshipChild),
			},
			expected: []*Resource{
				newResourceForTesting("a", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipChild),
				newResourceForTesting("c", enum.RelationshipChild),
			},
			strategy: []ResourceChunkSortStrategy{&TestResourceChunkSortStrategy{}},
		},
		{
			name: "single layer with relationship",
			input: []*Resource{
				newResourceForTesting("c", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipContainer),
				newResourceForTesting("a", enum.RelationshipChild),
			},
			expected: []*Resource{
				newResourceForTesting("a", enum.RelationshipChild),
				newResourceForTesting("c", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipContainer),
			},
			strategy: []ResourceChunkSortStrategy{&testAllSkipChunkSortStrategy{}, &TestResourceChunkSortStrategy{}},
		},
		{
			name: "multiple layer with name",
			input: []*Resource{
				newResourceForTesting("c", enum.RelationshipChild,
					newResourceForTesting("c", enum.RelationshipChild),
					newResourceForTesting("b", enum.RelationshipChild),
					newResourceForTesting("a", enum.RelationshipChild)),
				newResourceForTesting("b", enum.RelationshipChild,
					newResourceForTesting("c", enum.RelationshipChild),
					newResourceForTesting("b", enum.RelationshipChild),
					newResourceForTesting("a", enum.RelationshipChild)),
				newResourceForTesting("a", enum.RelationshipChild,
					newResourceForTesting("c", enum.RelationshipChild),
					newResourceForTesting("b", enum.RelationshipChild),
					newResourceForTesting("a", enum.RelationshipChild)),
			},
			expected: []*Resource{
				newResourceForTesting("a", enum.RelationshipChild,
					newResourceForTesting("a", enum.RelationshipChild),
					newResourceForTesting("b", enum.RelationshipChild),
					newResourceForTesting("c", enum.RelationshipChild)),
				newResourceForTesting("b", enum.RelationshipChild,
					newResourceForTesting("a", enum.RelationshipChild),
					newResourceForTesting("b", enum.RelationshipChild),
					newResourceForTesting("c", enum.RelationshipChild)),
				newResourceForTesting("c", enum.RelationshipChild,
					newResourceForTesting("a", enum.RelationshipChild),
					newResourceForTesting("b", enum.RelationshipChild),
					newResourceForTesting("c", enum.RelationshipChild)),
			},
			strategy: []ResourceChunkSortStrategy{&TestResourceChunkSortStrategy{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sorter := NewResourceSorter(tc.strategy...)
			result, err := sorter.SortAll(nil, tc.input)
			if err != nil {
				t.Errorf("unexpected error %s", err.Error())
			}
			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("sort result mismatch\n%s", diff)
			}
		})
	}
}
