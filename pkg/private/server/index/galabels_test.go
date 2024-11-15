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

package index

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/parameters"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestGALabelTagGenerator(t *testing.T) {
	testCases := []struct {
		name       string
		before     func()
		after      func()
		wantResult []string
	}{
		{
			name: "With ga labels",
			before: func() {
				parameters.Private.GALabels = testutil.P("foo=bar,qux=quux")
			},
			after: func() {
				parameters.Private.GALabels = nil
			},
			wantResult: []string{`<meta id="ga-meta-foo" content="bar">`, `<meta id="ga-meta-qux" content="quux">`},
		},
		{
			name: "Without ga labels",
			before: func() {
				parameters.Private.GALabels = testutil.P("")
			},
			after: func() {
				parameters.Private.GALabels = nil
			},
			wantResult: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			got := (&GALabelTagGenerator{}).GenerateTags()
			if diff := cmp.Diff(tc.wantResult, got); diff != "" {
				t.Errorf("GenerateTags() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
