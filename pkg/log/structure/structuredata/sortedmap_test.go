package structuredata

import (
	"fmt"
	"math/rand"
	"testing"

	"gopkg.in/yaml.v3"
)

func shufleFields(fields []string) []string {
	for i := len(fields) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		fields[i], fields[j] = fields[j], fields[i]
	}
	return fields
}

func TestUnorderedMarshal(t *testing.T) {
	sm := newSortedMap()
	fields := []string{}
	for i := 0; i < 10; i++ {
		fields = append(fields, fmt.Sprintf("field-%d", i))
	}
	shufleFields(fields)
	yamlStr := ""
	for i, field := range fields {
		yamlStr += fmt.Sprintf("%s: %d\n", field, i)
		sm.AddNextField(field, i)
	}

	result, err := yaml.Marshal(sm)
	if err != nil {
		t.Fatal(err)
	}

	resultYaml := string(result)
	if yamlStr != resultYaml {
		t.Errorf("Result is not matching with the input YAML data\nEXPECTED:\n\n%s\n\nACTUAL:\n\n%s", yamlStr, resultYaml)
	}
}

func TestUnorderedMarshalWithNilField(t *testing.T) {
	sm := newSortedMap()
	fields := []string{}
	for i := 0; i < 10; i++ {
		fields = append(fields, fmt.Sprintf("field-%d", i))
	}
	shufleFields(fields)
	yamlStr := ""
	for _, field := range fields {
		yamlStr += fmt.Sprintf("%s: null\n", field)
		sm.AddNextField(field, nil)
	}

	result, err := yaml.Marshal(sm)
	if err != nil {
		t.Fatal(err)
	}

	resultYaml := string(result)
	if yamlStr != resultYaml {
		t.Errorf("Result is not matching with the input YAML data\nEXPECTED:\n\n%s\n\nACTUAL:\n\n%s", yamlStr, resultYaml)
	}
}
