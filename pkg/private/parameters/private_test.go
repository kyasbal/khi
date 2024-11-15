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

package parameters

import (
	"flag"
	"os"
	"testing"

	khiflag "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/flag"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil"
	"github.com/google/go-cmp/cmp"
)

func prepareFlagParsingTest(t *testing.T) {
	t.Helper()
	khiflag.Reset()
}

func TestPrivateParameterss(t *testing.T) {
	testCases := []struct {
		name   string
		want   *PrivateParameters
		before func()
	}{
		{
			name: "default",
			want: &PrivateParameters{
				InspectionMode:   testutil.P(false),
				IAMToken:         testutil.P(""),
				GALabels:         testutil.P(""),
				DisableAnalytics: testutil.P(false),
				AnalyticsDebug:   testutil.P(false),
			},
			before: func() {
				os.Args = []string{os.Args[0]}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
		},
		{
			name: "with iam token",
			want: &PrivateParameters{
				InspectionMode:   testutil.P(true),
				IAMToken:         testutil.P("foo"),
				GALabels:         testutil.P(""),
				DisableAnalytics: testutil.P(false),
				AnalyticsDebug:   testutil.P(false),
			},
			before: func() {
				os.Args = []string{os.Args[0], "--iam-token=foo"}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prepareFlagParsingTest(t)
			tc.before()
			store := &PrivateParameters{}
			parameters.ResetStore()
			parameters.AddStore(store)
			err := parameters.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, store); diff != "" {
				t.Errorf("unexpected result (-want +got)\n%s", diff)
			}
		})
	}
}

func TestPrivateParameters_GetMapOfGALabels(t *testing.T) {
	testCases := []struct {
		name   string
		want   map[string]string
		params *PrivateParameters
	}{
		{
			name: "empty",
			want: map[string]string{},
			params: &PrivateParameters{
				GALabels: testutil.P(""),
			},
		},
		{
			name: "white space",
			want: map[string]string{},
			params: &PrivateParameters{
				GALabels: testutil.P("  "),
			},
		},
		{
			name: "single",
			want: map[string]string{"key1": "value1"},
			params: &PrivateParameters{
				GALabels: testutil.P("key1=value1"),
			},
		},
		{
			name: "multiple",
			want: map[string]string{"key1": "value1", "key2": "value2"},
			params: &PrivateParameters{
				GALabels: testutil.P("key1=value1,key2=value2"),
			},
		},
		{
			name: "multiple with spaces",
			want: map[string]string{"key1": "value1", "key2": "value2"},
			params: &PrivateParameters{
				GALabels: testutil.P("key1=value1 ,  key2=value2"),
			},
		},
		{
			name: "with null",
			want: map[string]string{"key1": "value1", "key2": "null"},
			params: &PrivateParameters{
				GALabels: testutil.P("key1=value1,key2"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prepareFlagParsingTest(t)
			got := tc.params.GetMapOfGALabels()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected result (-want +got)\n%s", diff)
			}
		})
	}
}
