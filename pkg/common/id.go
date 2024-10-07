package common

import "github.com/google/uuid"

// NewUUID returns the random UUID in string
func NewUUID() string {
	return uuid.Must(uuid.NewUUID()).String()
}
