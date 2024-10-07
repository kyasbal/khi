package rawlogs

import (
	"time"
)

// Repository stores logs after fetching logs
// Implementation shouldn't keep the given data on memory not to consume huge memory.
type Repository interface {
	Write(timestamp time.Time, data []byte) error
	IterateInSortedOrder() LogIterator
	Dispose() error
}

type LogIterator interface {
	HasNext() bool
	Next() ([]byte, error)
	Reset()
}
