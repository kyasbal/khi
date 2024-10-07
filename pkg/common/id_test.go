package common

import "testing"

func TestNewUUID(t *testing.T) {
	uuid1 := NewUUID()
	if uuid1 == "" {
		t.Error("NewUUID returned an empty string")
	}

	uuid2 := NewUUID()
	if uuid1 == uuid2 {
		t.Error("NewUUID returned the same UUID twice")
	}
}
