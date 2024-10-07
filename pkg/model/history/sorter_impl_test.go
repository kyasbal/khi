package history

import (
	"errors"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/google/go-cmp/cmp"
)

type resourceChunkSortStrategyTestCase struct {
	Name          string
	Chunk         []*Resource
	Parents       []*Resource
	ExpectedChunk []*Resource
	ExpectedError error
}

func testResourceChunkSortStrategy(t *testing.T, name string, sortStrategy ResourceChunkSortStrategy, testCases ...resourceChunkSortStrategyTestCase) {
	t.Run(name, func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				relationship := enum.RelationshipChild
				if len(tc.Chunk) > 0 {
					relationship = tc.Chunk[0].Relationship
				}
				actual, err := sortStrategy.SortChunk(nil, tc.Parents, relationship, tc.Chunk)
				if tc.ExpectedError != nil {
					if !errors.Is(err, tc.ExpectedError) {
						t.Errorf("error is not matching with the expected error.expected:%s,actual:%s", tc.ExpectedError.Error(), err.Error())
					}
				} else {
					if err != nil {
						t.Errorf("unexpected error %s", err.Error())
					}
					if diff := cmp.Diff(tc.ExpectedChunk, actual); diff != "" {
						t.Errorf("non matching result\n%s", diff)
					}
				}
			})
		}
	})
}

func TestNameSortStrategy(t *testing.T) {
	testResourceChunkSortStrategy(t, "with root layer", NewNameSortStrategy(0, []string{
		"@a", "@b", "@c",
	}), resourceChunkSortStrategyTestCase{
		Name: "simple sort without special keys",
		Chunk: []*Resource{
			newResourceForTesting("c", enum.RelationshipChild),
			newResourceForTesting("b", enum.RelationshipChild),
			newResourceForTesting("a", enum.RelationshipChild),
		},
		Parents: []*Resource{},
		ExpectedChunk: []*Resource{
			newResourceForTesting("a", enum.RelationshipChild),
			newResourceForTesting("b", enum.RelationshipChild),
			newResourceForTesting("c", enum.RelationshipChild),
		},
	},
		resourceChunkSortStrategyTestCase{
			Name: "simple sort with special keys",
			Chunk: []*Resource{
				newResourceForTesting("@c", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipChild),
				newResourceForTesting("a", enum.RelationshipChild),
			},
			Parents: []*Resource{},
			ExpectedChunk: []*Resource{
				newResourceForTesting("@c", enum.RelationshipChild),
				newResourceForTesting("a", enum.RelationshipChild),
				newResourceForTesting("b", enum.RelationshipChild),
			},
		},
		resourceChunkSortStrategyTestCase{
			Name: "for different layer",
			Chunk: []*Resource{
				newResourceForTesting("a", enum.RelationshipChild),
			},
			Parents:       []*Resource{newResourceForTesting("parent", enum.RelationshipChild)},
			ExpectedError: ErrorSortSkipped,
		})
}

func TestUnreachableSortStrategy(t *testing.T) {
	testResourceChunkSortStrategy(t, "-", &UnreachableSortStrategy{}, resourceChunkSortStrategyTestCase{
		Name: "simple sort",
		Chunk: []*Resource{
			newResourceForTesting("c", enum.RelationshipChild),
			newResourceForTesting("b", enum.RelationshipChild),
			newResourceForTesting("a", enum.RelationshipChild),
		},
		Parents: []*Resource{},
		ExpectedChunk: []*Resource{
			newResourceForTesting("a", enum.RelationshipChild),
			newResourceForTesting("b", enum.RelationshipChild),
			newResourceForTesting("c", enum.RelationshipChild),
		},
	})
}
