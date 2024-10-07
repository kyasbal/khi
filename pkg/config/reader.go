package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var DEFAULT_CONFIG_LOCATION = "./resources/config.yml"

var parsedConfigCache *ConfigFile = nil

func DefaultConfig() *ConfigFile {
	location, found := os.LookupEnv("KHI_CONFIG_LOCATION")
	if !found {
		location = DEFAULT_CONFIG_LOCATION
	}
	if parsedConfigCache == nil {
		file, err := os.Open(location)
		if err != nil {
			panic("Failed to open default configuration")
		}
		config, err := readConfig(file)
		if err != nil {
			panic(err)
		}
		parsedConfigCache = config
	}
	return parsedConfigCache
}

func readConfig(r io.Reader) (*ConfigFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration\n%v", err)
	}
	var config ConfigFile
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data\n%v", err)
	}
	return &config, nil
}
