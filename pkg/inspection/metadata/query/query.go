package query

import (
	"slices"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var QueryMetadataKey = "query"

type QueryItem struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

type QueryMetadata struct {
	Queries []*QueryItem
	lock    sync.Mutex
}

// Labels implements metadata.Metadata.
func (*QueryMetadata) Labels() *task.LabelSet {
	return task.NewLabelSet(metadata.IncludeInDryRunResult(), metadata.IncludeInRunResult())
}

// ToSerializable implements metadata.Metadata.
func (q *QueryMetadata) ToSerializable() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	slices.SortFunc(q.Queries, func(a, b *QueryItem) int { return strings.Compare(a.Id, b.Id) })
	return q.Queries
}

func (q *QueryMetadata) SetQuery(id string, name string, queryString string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for _, qi := range q.Queries {
		if qi.Id == id {
			qi.Name = name
			qi.Query = queryString
			return
		}
	}
	q.Queries = append(q.Queries, &QueryItem{
		Id:    id,
		Name:  name,
		Query: queryString,
	})
}

var _ metadata.Metadata = (*QueryMetadata)(nil)

type QueryMetadataFactory struct{}

// Instanciate implements metadata.MetadataFactory.
func (q *QueryMetadataFactory) Instanciate() metadata.Metadata {
	return &QueryMetadata{}
}

var _ metadata.MetadataFactory = (*QueryMetadataFactory)(nil)
