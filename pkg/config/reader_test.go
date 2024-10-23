package config

import (
	"os"
	"testing"
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
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if config == nil {
		t.Errorf("got nil, want config object")
	}
}
