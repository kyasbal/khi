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

import "errors"

type NodeType int

const (
	ScalarNodeType   = 0
	SequenceNodeType = 1
	MapNodeType      = 2
)

// NodeChildrenIterator is a type to represent the iterator returned from the Children method of Node interface.
type NodeChildrenIterator = func(func(key NodeChildrenKey, value Node) bool)

type Node interface {
	Type() NodeType
	NodeScalarValue() (any, error)
	Children() NodeChildrenIterator
}

// NodeChildrenElement represents an item of Chidlren of a Node.
type NodeChildrenKey struct {
	// Index is the index of the children
	Index int
	// Key is the key of the children in the map.
	// This value is empty when the Node is a sequence and not a map.
	Key string
}

var ErrNonScalarNode = errors.New("this is not a scalar node but called NodeScalarValue method")
