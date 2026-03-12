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
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/google/go-cmp/cmp"
)

type invalidResourceContainer struct{}

func (i *invalidResourceContainer) GetType() googlecloud.ResourceContainerType {
	return googlecloud.ResourceContainerInvalid
}

func (i *invalidResourceContainer) Identifier() string {
	return "invalid"
}

func TestToGoogleAdminLink(t *testing.T) {
	testCases := []struct {
		name          string
		container     googlecloud.ResourceContainer
		justification string
		want          string
		wantErr       bool
	}{
		{
			name:          "project container with buganizer justification",
			container:     googlecloud.Project("my-project-123"),
			justification: "b/123456",
			want:          "http://ga/x/operation/all-adapters/search?query=my-project-123&jt=1&jv=123456",
			wantErr:       false,
		},
		{
			name:          "project container with vector justification",
			container:     googlecloud.Project("my-project-123"),
			justification: "vector/987654",
			want:          "http://ga/x/operation/all-adapters/search?query=my-project-123&jt=28&jv=987654",
			wantErr:       false,
		},
		{
			name:          "unsupported justification type",
			container:     googlecloud.Project("my-project-123"),
			justification: "unknown/123",
			want:          "",
			wantErr:       true,
		},
		{
			name:          "unsupported container type",
			container:     &invalidResourceContainer{},
			justification: "b/123456",
			want:          "",
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToGoogleAdminLink(tc.container, tc.justification)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ToGoogleAdminLink() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("ToGoogleAdminLink() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestToGKEAdminLink(t *testing.T) {
	testCases := []struct {
		name          string
		container     googlecloud.ResourceContainer
		justification string
		want          string
		wantErr       bool
	}{
		{
			name:          "project container with buganizer justification",
			container:     googlecloud.Project("my-project-123"),
			justification: "b/123456",
			want:          "http://ga/gke/search?query=my-project-123&jt=1&jv=123456",
			wantErr:       false,
		},
		{
			name:          "project container with vector justification",
			container:     googlecloud.Project("my-project-123"),
			justification: "vector/987654",
			want:          "http://ga/gke/search?query=my-project-123&jt=28&jv=987654",
			wantErr:       false,
		},
		{
			name:          "unsupported justification type",
			container:     googlecloud.Project("my-project-123"),
			justification: "unknown/123",
			want:          "",
			wantErr:       true,
		},
		{
			name:          "unsupported container type",
			container:     &invalidResourceContainer{},
			justification: "b/123456",
			want:          "",
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToGKEAdminLink(tc.container, tc.justification)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ToGKEAdminLink() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("ToGKEAdminLink() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
