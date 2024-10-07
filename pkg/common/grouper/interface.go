package grouper

// Grouper is a base type to split the source group into named list of results.
type Grouper[T any, K comparable] interface {
	Group(input []T) map[K][]T
}
