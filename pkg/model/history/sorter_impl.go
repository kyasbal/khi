package history

import (
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

type NameSortStrategy struct {
	PrioritizedKeys []string
	Layer           int
}

// SortChunk implements ResourceChunkSortStrategy.
func (n *NameSortStrategy) SortChunk(builder *Builder, parents []*Resource, groupedRelationship enum.ParentRelationShip, chunk []*Resource) ([]*Resource, error) {
	if len(parents) != n.Layer {
		return nil, ErrorSortSkipped
	}
	sortResult := slices.Clone(chunk)
	keyToIndex := make(map[string]int)
	for index, key := range n.PrioritizedKeys {
		keyToIndex[key] = index
	}
	priority := func(key string) int {
		index, ok := keyToIndex[key]
		if ok {
			return index
		}
		return len(n.PrioritizedKeys)
	}
	sort.Slice(sortResult, func(i, j int) bool {
		a := priority(sortResult[i].ResourceName)
		b := priority(sortResult[j].ResourceName)
		reltypeA := sortResult[i].Relationship
		relTypeB := sortResult[j].Relationship
		switch {
		case a != b:
			return a < b
		case reltypeA != relTypeB:
			return reltypeA < relTypeB
		default:
			return sortResult[i].ResourceName < sortResult[j].ResourceName
		}
	})
	return sortResult, nil
}

var _ ResourceChunkSortStrategy = (*NameSortStrategy)(nil)

func NewNameSortStrategy(layer int, prioritizedKeys []string) *NameSortStrategy {
	return &NameSortStrategy{
		Layer:           layer,
		PrioritizedKeys: prioritizedKeys,
	}
}

// UnreachableSortStrategy is the default sort strategy catches all.
// Hitting this sorter is unexpected but implemented not to crush because of bad output from parsers.
type UnreachableSortStrategy struct {
}

// SortChunk implements ResourceChunkSortStrategy.
func (u *UnreachableSortStrategy) SortChunk(builder *Builder, parents []*Resource, groupedRelationship enum.ParentRelationShip, chunk []*Resource) ([]*Resource, error) {
	cloned := slices.Clone(chunk)
	slices.SortFunc(cloned, func(a, b *Resource) int {
		return strings.Compare(a.ResourceName, b.ResourceName)
	})
	for _, c := range chunk {
		slog.Warn(fmt.Sprintf("hitting unreachable sorter: %s", c.FullResourcePath))
	}
	return cloned, nil
}

var _ ResourceChunkSortStrategy = (*UnreachableSortStrategy)(nil)

type FirstRevisionTimeSortStrategy struct {
	TargetRelationship enum.ParentRelationShip
}

// SortChunk implements ResourceChunkSortStrategy.
func (b *FirstRevisionTimeSortStrategy) SortChunk(builder *Builder, parents []*Resource, groupedRelationship enum.ParentRelationShip, chunk []*Resource) ([]*Resource, error) {
	if groupedRelationship != b.TargetRelationship {
		return nil, ErrorSortSkipped
	}
	cloned := slices.Clone(chunk)
	slices.SortFunc(cloned, func(a, b *Resource) int {
		abuilder := builder.GetTimelineBuilder(a.FullResourcePath)
		bbuilder := builder.GetTimelineBuilder(b.FullResourcePath)
		arevs := abuilder.timeline.Revisions
		brevs := bbuilder.timeline.Revisions
		if len(arevs) == 0 {
			return -1
		}
		if len(brevs) == 0 {
			return 1
		}
		return int(arevs[0].ChangeTime.Sub(brevs[0].ChangeTime))
	})
	return cloned, nil
}

var _ ResourceChunkSortStrategy = (*FirstRevisionTimeSortStrategy)(nil)
