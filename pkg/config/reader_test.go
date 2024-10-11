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
		t.Errorf("Failed to read configuration\n%v", err)
	}
	if config == nil {
		t.Errorf("Failed to read configuration. Config is nil")
	}
}
