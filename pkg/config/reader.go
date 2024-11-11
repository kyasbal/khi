// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if parsedConfigCache == nil {
		file, err := os.Open(DEFAULT_CONFIG_LOCATION)
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
