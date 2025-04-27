// Copyright 2025 Google LLC
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

package structurev2

import (
	"errors"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrMultipleDocumentNodeFound = errors.New("multiple document node found in a yaml. FromYAML only supports a single document node")
var ErrAliasNodeNotSupported = errors.New("alias node is not supported in a yaml. FromYAML does not support alias node")
var ErrUnknownYAMLNodeKind = errors.New("unknown yaml node kind")

func FromYAML(yamlStr string) (Node, error) {
	// Parse yaml string as yaml.Node instead of `any` type to keep the order of the map keys in the original YAML.
	var root yaml.Node
	err := yaml.Unmarshal([]byte(yamlStr), &root)
	if err != nil {
		return nil, err
	}

	return fromYAMLNode(&root)
}

func fromYAMLNode(node *yaml.Node) (Node, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 1 {
			return fromYAMLNode(node.Content[0])
		}
		return nil, ErrMultipleDocumentNodeFound
	case yaml.SequenceNode:
		return fromSequenceYAMLNode(node)
	case yaml.MappingNode:
		return fromMappingYAMLNode(node)
	case yaml.ScalarNode:
		return fromScalarYAMLNode(node)
	case yaml.AliasNode:
		return nil, ErrAliasNodeNotSupported
	default:
		return nil, ErrUnknownYAMLNodeKind
	}
}

func fromSequenceYAMLNode(node *yaml.Node) (Node, error) {
	children := []Node{}
	for _, content := range node.Content {
		child, err := fromYAMLNode(content)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return &StandardSequenceNode{value: children}, nil
}

func fromMappingYAMLNode(node *yaml.Node) (Node, error) {
	keys := []string{}
	values := []Node{}
	for i, content := range node.Content {
		if i%2 == 0 { // yaml.Node holds its map key-values as the sequence of a structure like key1,value1,key2,value2,...etc
			keys = append(keys, content.Value)
		} else {
			child, err := fromYAMLNode(content)
			if err != nil {
				return nil, err
			}
			values = append(values, child)
		}
	}
	return &StandardMapNode{keys: keys, values: values}, nil
}

func fromScalarYAMLNode(node *yaml.Node) (Node, error) {
	// Scalar yaml.Node holds its value as string but Tag field contains its type.
	// https://github.com/go-yaml/yaml/blob/944c86a7d29391925ed6ac33bee98a0516f1287a/resolve.go#L71-L80
	switch node.Tag {
	case "!!null":
		return &StandardScalarNode[any]{value: nil}, nil
	case "!!bool":
		boolValue, err := strconv.ParseBool(node.Value)
		if err != nil {
			return nil, err
		}
		return &StandardScalarNode[bool]{value: boolValue}, nil
	case "!!str":
		return &StandardScalarNode[string]{value: node.Value}, nil
	case "!!int":
		intValue, err := strconv.Atoi(node.Value)
		if err != nil {
			return nil, err
		}
		return &StandardScalarNode[int]{value: intValue}, nil
	case "!!float":
		floatValue, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return nil, err
		}
		return &StandardScalarNode[float64]{value: floatValue}, nil
	case "!!timestamp":
		timestampValue, err := time.Parse(time.RFC3339, node.Value)
		if err != nil {
			return nil, err
		}
		return &StandardScalarNode[time.Time]{value: timestampValue}, nil
	default:
		return &StandardScalarNode[string]{value: node.Value}, nil
	}
}
