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

package privatecsmcp_contract

import (
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
)

var (
	pathJSONPayloadMessage   = structured.CompileFieldPath("jsonPayload.message")
	pathJSONPayloadTimestamp = structured.CompileFieldPath("jsonPayload.timestamp")
	pathLabelsInstanceID     = structured.CompileFieldPath("labels.instanceId")
)

// PodIdentifier holds the namespace, name, and an optional connection ID of a Pod.
type PodIdentifier struct {
	Namespace    string
	Name         string
	ConnectionID string
}

// CSMCPFieldSet holds structured log data extracted from the log entry.
type CSMCPFieldSet struct {
	Timestamp  *time.Time
	Message    string
	Pods       []PodIdentifier
	InstanceID string
}

// ExtractCSMCP extracts CSMCPFieldSet from a NodeReader.
func ExtractCSMCP(reader *structured.NodeReader) (CSMCPFieldSet, error) {
	if mock, ok := structured.GetMock[CSMCPFieldSet](reader); ok {
		return mock, nil
	}

	message := reader.ReadStringOrDefault(pathJSONPayloadMessage, "")

	var pods []PodIdentifier
	for _, word := range strings.Fields(message) {
		if p := parsePodIdentifier(word); p != nil {
			pods = append(pods, *p)
		}
	}

	var timestamp *time.Time
	if ts, err := reader.ReadTimestamp(pathJSONPayloadTimestamp); err == nil {
		timestamp = &ts
	}

	instanceID := reader.ReadStringOrDefault(pathLabelsInstanceID, "unknown")

	return CSMCPFieldSet{
		Timestamp:  timestamp,
		Message:    message,
		Pods:       pods,
		InstanceID: instanceID,
	}, nil
}

// parsePodIdentifier parses a word to extract Pod details (name, namespace, optional connection ID).
// It returns nil if the word does not match the expected format.
func parsePodIdentifier(word string) *PodIdentifier {
	word = strings.TrimPrefix(word, "node:")

	if strings.Count(word, ".") != 1 {
		return nil
	}

	dotIdx := strings.IndexByte(word, '.')
	if dotIdx <= 0 || dotIdx == len(word)-1 {
		return nil
	}

	podName := word[:dotIdx]
	if !isValidPodName(podName) {
		return nil
	}

	rest := word[dotIdx+1:]
	namespace, connectionID := extractNamespaceAndConnectionID(rest)
	if !isValidNamespace(namespace) {
		return nil
	}

	return &PodIdentifier{
		Name:         podName,
		Namespace:    namespace,
		ConnectionID: connectionID,
	}
}

// isValidPodName checks if the given string is a valid Kubernetes pod name component.
func isValidPodName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

// isValidNamespace checks if the given string is a valid Kubernetes namespace component.
func isValidNamespace(namespace string) bool {
	if len(namespace) == 0 || !(namespace[0] >= 'a' && namespace[0] <= 'z') {
		return false
	}
	for i := 0; i < len(namespace); i++ {
		c := namespace[i]
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

// extractNamespaceAndConnectionID separates the namespace from the connection ID if it exists.
func extractNamespaceAndConnectionID(rest string) (string, string) {
	lastDashIdx := strings.LastIndexByte(rest, '-')
	if lastDashIdx > 0 && lastDashIdx < len(rest)-1 {
		possibleConnID := rest[lastDashIdx+1:]
		if isAllDigits(possibleConnID) {
			return rest[:lastDashIdx], possibleConnID
		}
	}
	return rest, ""
}

// isAllDigits checks if the given string consists entirely of digits.
func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
