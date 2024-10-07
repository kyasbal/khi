package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MustOpen(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

func TestDefaultConfigurationIsValid(t *testing.T) {
	config, err := readConfig(MustOpen("../../resources/config.yml"))
	assert.Nil(t, err)
	assert.NotNil(t, config)
}
