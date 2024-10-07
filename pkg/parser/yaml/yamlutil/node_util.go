package yamlutil

import (
	"errors"

	"gopkg.in/yaml.v3"
)

var NodeNotFoundError = errors.New("Node not found")

func NewEmptyMapNode() *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, 0),
	}
}

func NewMapElementWithScalarValue(key string, value string) []*yaml.Node {
	return []*yaml.Node{
		{
			Kind:    yaml.ScalarNode,
			Content: make([]*yaml.Node, 0),
			Value:   key,
		},
		{
			Kind:    yaml.ScalarNode,
			Content: make([]*yaml.Node, 0),
			Value:   value,
		},
	}
}

func NewScalarNode(value string) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.ScalarNode,
		Content: make([]*yaml.Node, 0),
		Value:   value,
	}
}

func DecomposeMapElement(mapNode *yaml.Node, index int) (string, *yaml.Node) {
	mapIndex := index * 2
	key := mapNode.Content[mapIndex]
	value := mapNode.Content[mapIndex+1]
	return key.Value, value
}

func GetMapElement(mapNode *yaml.Node, key string) (*yaml.Node, error) {
	for i := 0; i < GetMapLength(mapNode); i++ {
		fieldName, node := DecomposeMapElement(mapNode, i)
		if key == fieldName {
			return node, nil
		}
	}
	return nil, NodeNotFoundError
}

func GetMapLength(mapNode *yaml.Node) int {
	return len(mapNode.Content) / 2
}
