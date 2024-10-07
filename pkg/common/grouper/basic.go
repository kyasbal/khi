package grouper

// BasicGrouper implements Grouper
type BasicGrouper[T any, K comparable] struct {
	GroupingFunc func(input T) K
}

func (b *BasicGrouper[T, K]) Group(input []T) map[K][]T {
	result := map[K][]T{}
	for _, v := range input {
		key := b.GroupingFunc(v)
		if _, found := result[key]; !found {
			result[key] = []T{}
		}
		result[key] = append(result[key], v)
	}
	return result
}

func NewBasicGrouper[T any, K comparable](groupingFunction func(input T) K) *BasicGrouper[T, K] {
	return &BasicGrouper[T, K]{
		GroupingFunc: groupingFunction,
	}
}
