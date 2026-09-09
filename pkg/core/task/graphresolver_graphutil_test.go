// Copyright 2026 Google LLC
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

package coretask

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCloneOutgoingGraph(t *testing.T) {
	testCases := []struct {
		name          string
		input         map[string][]string
		want          map[string][]string
		mutateKey     string
		mutateVal     string
		checkOrigKey  string
		checkOrigWant string
	}{
		{
			name:  "empty map",
			input: map[string][]string{},
			want:  map[string][]string{},
		},
		{
			name: "multi-node map with deep copy verification",
			input: map[string][]string{
				"A": {"B", "C"},
				"B": {"C"},
			},
			want: map[string][]string{
				"A": {"B", "C"},
				"B": {"C"},
			},
			mutateKey:     "A",
			mutateVal:     "Z",
			checkOrigKey:  "A",
			checkOrigWant: "B",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cloned := cloneOutgoingGraph(tc.input)
			if diff := cmp.Diff(tc.want, cloned); diff != "" {
				t.Errorf("cloneOutgoingGraph() mismatch (-want +got):\n%s", diff)
			}
			if tc.mutateKey != "" {
				cloned[tc.mutateKey][0] = tc.mutateVal
				if tc.input[tc.checkOrigKey][0] != tc.checkOrigWant {
					t.Errorf("cloneOutgoingGraph() shallow copy detected, modifying cloned modified original: got %q, want %q", tc.input[tc.checkOrigKey][0], tc.checkOrigWant)
				}
			}
		})
	}
}

func TestVerifyAcyclic(t *testing.T) {
	testCases := []struct {
		name          string
		outgoing      map[string][]string
		inDegree      map[string]int
		nodeCount     int
		wantErr       bool
		wantSubstring string
	}{
		{
			name:      "empty graph is acyclic",
			outgoing:  map[string][]string{},
			inDegree:  map[string]int{},
			nodeCount: 0,
			wantErr:   false,
		},
		{
			name: "single node without edges is acyclic",
			outgoing: map[string][]string{
				"A": nil,
			},
			inDegree: map[string]int{
				"A": 0,
			},
			nodeCount: 1,
			wantErr:   false,
		},
		{
			name: "linear DAG A -> B -> C is acyclic",
			outgoing: map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": nil,
			},
			inDegree: map[string]int{
				"A": 0,
				"B": 1,
				"C": 1,
			},
			nodeCount: 3,
			wantErr:   false,
		},
		{
			name: "diamond DAG is acyclic",
			outgoing: map[string][]string{
				"A": {"B", "C"},
				"B": {"D"},
				"C": {"D"},
				"D": nil,
			},
			inDegree: map[string]int{
				"A": 0,
				"B": 1,
				"C": 1,
				"D": 2,
			},
			nodeCount: 4,
			wantErr:   false,
		},
		{
			name: "self-loop cycle A -> A fails",
			outgoing: map[string][]string{
				"A": {"A"},
			},
			inDegree: map[string]int{
				"A": 1,
			},
			nodeCount:     1,
			wantErr:       true,
			wantSubstring: "cyclic dependency",
		},
		{
			name: "2-node cycle A -> B -> A fails",
			outgoing: map[string][]string{
				"A": {"B"},
				"B": {"A"},
			},
			inDegree: map[string]int{
				"A": 1,
				"B": 1,
			},
			nodeCount:     2,
			wantErr:       true,
			wantSubstring: "cyclic dependency",
		},
		{
			name: "cycle with upstream entry S -> A -> B -> A fails",
			outgoing: map[string][]string{
				"S": {"A"},
				"A": {"B"},
				"B": {"A"},
			},
			inDegree: map[string]int{
				"S": 0,
				"A": 2,
				"B": 1,
			},
			nodeCount:     3,
			wantErr:       true,
			wantSubstring: "cyclic dependency",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyAcyclic(tc.outgoing, tc.inDegree, tc.nodeCount)
			if (err != nil) != tc.wantErr {
				t.Fatalf("verifyAcyclic() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr && tc.wantSubstring != "" {
				if !strings.Contains(err.Error(), tc.wantSubstring) {
					t.Errorf("verifyAcyclic() error %q does not contain %q", err.Error(), tc.wantSubstring)
				}
			}
		})
	}
}

func TestIsReachable(t *testing.T) {
	outgoing := map[string][]string{
		"A": {"B", "C"},
		"B": {"D"},
		"C": {"E"},
		"D": {"F"},
		"E": nil,
		"F": nil,
		"X": {"Y"},
		"Y": nil,
	}

	testCases := []struct {
		name       string
		startNode  string
		targetNode string
		want       bool
	}{
		{
			name:       "direct reachability A -> B",
			startNode:  "A",
			targetNode: "B",
			want:       true,
		},
		{
			name:       "transitive reachability A -> B -> D -> F",
			startNode:  "A",
			targetNode: "F",
			want:       true,
		},
		{
			name:       "transitive reachability A -> C -> E",
			startNode:  "A",
			targetNode: "E",
			want:       true,
		},
		{
			name:       "reverse direction is not reachable B -> A",
			startNode:  "B",
			targetNode: "A",
			want:       false,
		},
		{
			name:       "disjoint component is not reachable A -> X",
			startNode:  "A",
			targetNode: "X",
			want:       false,
		},
		{
			name:       "same node without self edge is not reachable",
			startNode:  "A",
			targetNode: "A",
			want:       false,
		},
		{
			name:       "node with no outgoing edges cannot reach anything",
			startNode:  "F",
			targetNode: "A",
			want:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isReachable(tc.startNode, tc.targetNode, outgoing)
			if got != tc.want {
				t.Errorf("isReachable(%q, %q) = %v, want %v", tc.startNode, tc.targetNode, got, tc.want)
			}
		})
	}
}

func TestExtractCyclicDependencyPath(t *testing.T) {
	testCases := []struct {
		name          string
		outgoing      map[string][]string
		inDegree      map[string]int
		wantSubstring string
	}{
		{
			name: "simple cycle A -> B -> A",
			outgoing: map[string][]string{
				"A": {"B"},
				"B": {"A"},
			},
			inDegree: map[string]int{
				"A": 1,
				"B": 1,
			},
			wantSubstring: "A -> B",
		},
		{
			name:     "no unresolved nodes returns empty string",
			outgoing: map[string][]string{},
			inDegree: map[string]int{
				"A": 0,
			},
			wantSubstring: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := extractCyclicDependencyPath(tc.outgoing, tc.inDegree)
			if tc.wantSubstring == "" {
				if got != "" {
					t.Errorf("extractCyclicDependencyPath() = %q, want empty string", got)
				}
			} else if !strings.Contains(got, tc.wantSubstring) {
				t.Errorf("extractCyclicDependencyPath() = %q, want containing %q", got, tc.wantSubstring)
			}
		})
	}
}
