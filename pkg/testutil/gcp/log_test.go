package gcp_test

import (
	"testing"
)

func TestIsValidLogQuery(t *testing.T) {
	err := IsValidLogQuery("\"")
	if err == nil {
		t.Errorf("got nil, want invalid query error")
	}
}
