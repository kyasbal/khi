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

package iamtoken

import (
	"net/http"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc/metadata"
)

func TestIamTokenCallOptionInjectorOption_ApplyToCallContext(t *testing.T) {
	testCases := []struct {
		desc      string
		prepare   func() *IAMTokenCallOptionInjectorOption
		container googlecloud.ResourceContainer
		wantToken string
	}{
		{
			desc: "from container specific token",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("foo"), "foo-token")
				return option
			},
			container: googlecloud.Project("foo"),
			wantToken: "foo-token",
		},
		{
			desc: "no token for container",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("bar"), "bar-token")
				return option
			},
			container: googlecloud.Project("foo"),
			wantToken: "",
		},
		{
			desc:      "nil container",
			prepare:   NewInjector,
			container: nil,
			wantToken: "",
		},
	}
	for _, tc := range testCases {
		option := tc.prepare()
		t.Run("ApplyToCallContext", func(t *testing.T) {
			t.Run(tc.desc, func(t *testing.T) {
				ctx := t.Context()
				ctx = option.ApplyToCallContext(ctx, tc.container)

				md, ok := metadata.FromOutgoingContext(ctx)
				if !ok {
					if tc.wantToken == "" {
						return
					}
					t.Fatal("metadata not found in context")
				}

				tokens := md.Get(iamTokenKey)
				if tc.wantToken == "" {
					if len(tokens) != 0 {
						t.Errorf("expected no token, got %v", tokens)
					}
					return
				}

				if len(tokens) != 1 {
					t.Fatalf("expected 1 token, got %d", len(tokens))
				}

				if tokens[0] != tc.wantToken {
					t.Errorf("expected tc.wantToken %q, got %q", tc.wantToken, tokens[0])
				}
			})
		})

		t.Run("ApplyToRawHTTPHeader", func(t *testing.T) {
			t.Run(tc.desc, func(t *testing.T) {
				header := http.Header{}

				option.ApplyToRawHTTPHeader(header, tc.container)

				if got := header.Get(iamTokenKey); got != tc.wantToken {
					t.Errorf("expected token %q, got %q", tc.wantToken, got)
				}
			})
		})

	}
}

func TestIamTokenCallOptionInjectorOption_HasTokenFor(t *testing.T) {
	testCases := []struct {
		name      string
		prepare   func() *IAMTokenCallOptionInjectorOption
		container googlecloud.ResourceContainer
		want      bool
	}{
		{
			name: "has token for specific container",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("foo"), "foo-token")
				return option
			},
			container: googlecloud.Project("foo"),
			want:      true,
		},
		{
			name: "does not have token for container",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("bar"), "bar-token")
				return option
			},
			container: googlecloud.Project("foo"),
			want:      false,
		},
		{
			name:      "nil container",
			prepare:   NewInjector,
			container: nil,
			want:      false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := tc.prepare()
			got := option.HasTokenFor(tc.container)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("HasTokenFor() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIamTokenCallOptionInjectorOption_SetTokenFor(t *testing.T) {
	testCases := []struct {
		name      string
		prepare   func() *IAMTokenCallOptionInjectorOption
		container googlecloud.ResourceContainer
		token     string
		check     func(t *testing.T, option *IAMTokenCallOptionInjectorOption)
	}{
		{
			name:      "set token for new container",
			prepare:   NewInjector,
			container: googlecloud.Project("foo"),
			token:     "my-token",
			check: func(t *testing.T, option *IAMTokenCallOptionInjectorOption) {
				if got := option.getTokenFor(googlecloud.Project("foo")); got != "my-token" {
					t.Errorf("expected token %q, got %q", "my-token", got)
				}
			},
		},
		{
			name: "overwrite existing token",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("foo"), "old-token")
				return option
			},
			container: googlecloud.Project("foo"),
			token:     "new-token",
			check: func(t *testing.T, option *IAMTokenCallOptionInjectorOption) {
				if got := option.getTokenFor(googlecloud.Project("foo")); got != "new-token" {
					t.Errorf("expected token %q, got %q", "new-token", got)
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := tc.prepare()
			option.SetTokenFor(tc.container, tc.token)
			tc.check(t, option)
		})
	}
}

func TestIamTokenCallOptionInjectorOption_getTokenFor(t *testing.T) {
	testCases := []struct {
		name      string
		prepare   func() *IAMTokenCallOptionInjectorOption
		container googlecloud.ResourceContainer
		want      string
	}{
		{
			name: "has token for specific container",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("foo"), "foo-token")
				return option
			},
			container: googlecloud.Project("foo"),
			want:      "foo-token",
		},
		{
			name: "does not have token for container",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("bar"), "bar-token")
				return option
			},
			container: googlecloud.Project("foo"),
			want:      "",
		},
		{
			name:      "nil container",
			prepare:   NewInjector,
			container: nil,
			want:      "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := tc.prepare()
			got := option.getTokenFor(tc.container)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("getTokenFor() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIamTokenCallOptionInjectorOption_RegisteredProjectIDs(t *testing.T) {
	testCases := []struct {
		name    string
		prepare func() *IAMTokenCallOptionInjectorOption
		want    []string
	}{
		{
			name: "multiple projects registered",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := NewInjector()
				option.SetTokenFor(googlecloud.Project("proj-a"), "token-a")
				option.SetTokenFor(googlecloud.Project("proj-b-tp"), "token-b")
				return option
			},
			want: []string{"proj-a", "proj-b-tp"},
		},
		{
			name:    "no tokens registered",
			prepare: NewInjector,
			want:    []string{},
		},
		{
			name: "nil receiver",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				return nil
			},
			want: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := tc.prepare()
			got := option.RegisteredProjectIDs()
			if diff := cmp.Diff(tc.want, got, cmpopts.SortSlices(func(a, b string) bool { return a < b }), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("RegisteredProjectIDs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
