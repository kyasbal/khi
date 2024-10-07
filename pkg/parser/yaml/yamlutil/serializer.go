package yamlutil

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func MarshalToYamlString(data any) (string, error) {
	result, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return escapeInvalidYamlMarshalResult(string(result)), nil
}

func escapeInvalidYamlMarshalResult(yamlStr string) string {
	// https://github.com/go-yaml/yaml/issues/940
	// Find /4- and replace it to /-
	yamlStr = strings.ReplaceAll(yamlStr, "- |4", "- |")
	yamlStr = strings.ReplaceAll(yamlStr, ": |4", ": |")
	return yamlStr
}
