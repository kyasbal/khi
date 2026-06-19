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

package privategkemaster_impl

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khierrors"
)

// readGoStructFromString finds the struct part of a specific structName in the given string and returns its fields as a map.
func readGoStructFromString(message string, structName string) map[string]string {
	splitted := strings.Split(message, structName)
	if len(splitted) > 1 {
		laterPart := splitted[1]
		if len(laterPart) == 0 {
			return map[string]string{}
		}
		if laterPart[0] == '{' {
			laterPart = laterPart[1:]
		}
		structPart := strings.Split(laterPart, "}")[0]
		fields := strings.Split(structPart, ",")
		result := map[string]string{}
		for _, field := range fields {
			keyValue := strings.Split(field, ":")
			if len(keyValue) == 2 {
				result[keyValue[0]] = keyValue[1]
			}
		}
		return result
	}
	return map[string]string{}
}

// readNextQuotedString extracts the content of the first double-quoted string found in the input message.
func readNextQuotedString(msg string) string {
	splitted := strings.Split(msg, "\"")
	if len(splitted) > 2 {
		return splitted[1]
	} else {
		return ""
	}
}

// slashSplittedPodNameToNamespaceAndName converts a slash-separated pod name (e.g., "kube-system/kube-dns-abcd") into its namespace and name components.
func slashSplittedPodNameToNamespaceAndName(name string) (string, string, error) {
	nameSplitted := strings.Split(name, "/")
	if len(nameSplitted) == 2 {
		return nameSplitted[0], nameSplitted[1], nil
	}
	return "", "", fmt.Errorf("invalid pod name format %q: %w", name, khierrors.ErrInvalidInput)
}
