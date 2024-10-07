package testlog

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parser/yaml/yamlutil"
	"gopkg.in/yaml.v3"
)

func BaseYaml(yamlStr string) TestLogOpt {
	return func(original *yaml.Node) (*yaml.Node, error) {
		if original != nil {
			return nil, fmt.Errorf("BaseYaml expects no previous TestLogOpt is given. But an instance of node was given")
		}
		if yamlStr == "" {
			return yamlutil.NewEmptyMapNode(), nil
		}
		var node yaml.Node
		err := yaml.Unmarshal([]byte(yamlStr), &node)
		if err != nil {
			return nil, err
		}
		return &node, err
	}
}
