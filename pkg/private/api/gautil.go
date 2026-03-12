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

package api

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
)

const operationName = "all-adapters"

// Justification types used in Google Admin URL parameter
// https://source.corp.google.com/piper///depot/google3/abuse/admin/justification/justification.proto;l=24;cl=869971733
const (
	// justificationTypeBuganizer is the justification type for Buganizer bug IDs.
	justificationTypeBuganizer = 1
	// justificationTypeVector is the justification type for Vector case numbers.
	justificationTypeVector = 28
)

// ToGoogleAdminLink generates a URL to search for the given resource container
// in Google Admin. The search uses the provided justification constraint.
// Currently, only googlecloud.ResourceContainerProject is supported. The justification
// string must be prefixed with "b/" for Buganizer bug IDs or "vector/" for Vector case numbers.
func ToGoogleAdminLink(container googlecloud.ResourceContainer, justification string) (string, error) {
	query, err := getAdminQuery(container)
	if err != nil {
		return "", err
	}
	jt, jv, err := getAdminJustificationTypeAndValue(justification)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://ga/x/operation/%s/search?query=%s&jt=%d&jv=%s", operationName, query, jt, jv), nil
}

// ToGKEAdminLink generates a URL to search for the given resource container
// in the GKE section of Google Admin. The search uses the provided justification constraint.
// Currently, only googlecloud.ResourceContainerProject is supported. The justification
// string must be prefixed with "b/" for Buganizer bug IDs or "vector/" for Vector case numbers.
func ToGKEAdminLink(container googlecloud.ResourceContainer, justification string) (string, error) {
	query, err := getAdminQuery(container)
	if err != nil {
		return "", err
	}

	jt, jv, err := getAdminJustificationTypeAndValue(justification)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://ga/gke/search?query=%s&jt=%d&jv=%s", query, jt, jv), nil
}

func getAdminQuery(container googlecloud.ResourceContainer) (string, error) {
	query := ""
	switch container.GetType() {
	case googlecloud.ResourceContainerProject:
		query, _ = strings.CutPrefix(container.Identifier(), "projects/")
	default:
		return "", fmt.Errorf("unsupported container type: %v", container.GetType())
	}
	return query, nil
}

func getAdminJustificationTypeAndValue(justification string) (jt int, jv string, err error) {
	switch {
	case strings.HasPrefix(justification, "b/"):
		jt = justificationTypeBuganizer
		jv = strings.TrimPrefix(justification, "b/")
	case strings.HasPrefix(justification, "vector/"):
		jt = justificationTypeVector
		jv = strings.TrimPrefix(justification, "vector/")
	default:
		return 0, "", fmt.Errorf("unsupported justification type: %s. Currently vector case number or buganizer ID is supported", justification)
	}
	return jt, jv, nil
}
