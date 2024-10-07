package gcp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidLogQuery(t *testing.T) {
	err := IsValidLogQuery("\"")
	assert.NotNil(t, err)
}
