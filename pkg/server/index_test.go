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

package server

import (
	"sort"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
	"github.com/google/go-cmp/cmp"
)

func wrapPointer[T any](t T) *T {
	return &t
}

func TestGenerateGaMetaTags(t *testing.T) {
	input := map[string]string{
		"foo": "bar",
		"qux": "quux",
	}

	result := generateGaMetaTags(input)
	sort.Strings(result)

	expect0 := `<meta id="ga-meta-foo" content="bar">`
	expect1 := `<meta id="ga-meta-qux" content="quux">`
	if len(result) != 2 {
		t.Errorf("expect len(result) = 2 but %d", len(result))
	}
	if result[0] != expect0 {
		t.Errorf("expect result[0] to be %s but %s", expect0, result[0])
	}
	if result[1] != expect1 {
		t.Errorf("expect result[1] to be %s but %s", expect1, result[1])
	}
}

func TestGetServerBasePathMetaTag(t *testing.T) {
	testCases := []struct {
		name       string
		before     func()
		after      func()
		wantResult string
	}{
		{
			name: "With server base path env",
			before: func() {
				parameters.Server.BasePath = wrapPointer("/foo/bar/")
			},
			after: func() {
				parameters.Server.BasePath = nil
			},
			wantResult: `<meta id="server-base-path" content="/foo/bar/">`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			got := getServerBasePathMetaTag()
			if got != tc.wantResult {
				t.Errorf("got %s, want %s", got, tc.wantResult)
			}
		})
	}
}

func TestGetBaseTag(t *testing.T) {
	testCases := []struct {
		name       string
		before     func()
		after      func()
		wantResult string
	}{
		{
			name: "With frontend resource base path env",
			before: func() {
				parameters.Server.FrontendResourceBasePath = wrapPointer("/foo/bar/")
			},
			after: func() {
				parameters.Server.FrontendResourceBasePath = nil
			},
			wantResult: `<base href="/foo/bar/">`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			got := getBaseTag()
			if got != tc.wantResult {
				t.Errorf("got %s, want %s", got, tc.wantResult)
			}
		})
	}
}

func TestReplaceDynamicPartOfInde(t *testing.T) {
	source := `<!--INJECT GENERATED CODE HERE FROM BACKEND--><base href="/" />`
	testCases := []struct {
		name       string
		before     func()
		after      func()
		wantSource string
	}{{
		name: "with basic values",
		before: func() {
			parameters.Private.GALabels = wrapPointer("key1=val1,key2=val2")
			parameters.Server.BasePath = wrapPointer("/basepath/")
			parameters.Server.FrontendResourceBasePath = wrapPointer("/frontend/")
		},
		after: func() {
			parameters.Private.GALabels = nil
			parameters.Server.BasePath = nil
			parameters.Server.FrontendResourceBasePath = nil
		},
		wantSource: `<base href="/frontend/">
<meta id="ga-meta-key1" content="val1">
<meta id="ga-meta-key2" content="val2">
<meta id="server-base-path" content="/basepath/">`,
	}, {
		name: "with the / base path(not to be removed with the patch code for the local dev server)",
		before: func() {
			parameters.Private.GALabels = wrapPointer("key1=val1,key2=val2")
			parameters.Server.BasePath = wrapPointer("/basepath/")
			parameters.Server.FrontendResourceBasePath = wrapPointer("/")
		},
		after: func() {
			parameters.Private.GALabels = nil
			parameters.Server.BasePath = nil
			parameters.Server.FrontendResourceBasePath = nil
		},
		wantSource: `<base href="/">
<meta id="ga-meta-key1" content="val1">
<meta id="ga-meta-key2" content="val2">
<meta id="server-base-path" content="/basepath/">`,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			got, err := replaceDynamicPartOfIndex(source)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}
			if diff := cmp.Diff(tc.wantSource, got); diff != "" {
				t.Errorf("the result is not matching with the expected response(-want +got)\n%s", diff)
			}
		})
	}
}
