package parameters

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrivateParameterss(t *testing.T) {
	testCases := []struct {
		name   string
		want   *PrivateParameters
		before func()
	}{
		{
			name: "default",
			want: &PrivateParameters{
				InspectionMode: wrapPointer(false),
				IAMToken:       wrapPointer(""),
				GALabels:       wrapPointer(""),
			},
			before: func() {
				os.Args = []string{os.Args[0]}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
		},
		{
			name: "with iam token",
			want: &PrivateParameters{
				InspectionMode: wrapPointer(true),
				IAMToken:       wrapPointer("foo"),
				GALabels:       wrapPointer(""),
			},
			before: func() {
				os.Args = []string{os.Args[0], "--iam-token=foo"}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			store := &PrivateParameters{}
			err := Parse(store)
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
				GALabels: wrapPointer(""),
			},
		},
		{
			name: "white space",
			want: map[string]string{},
			params: &PrivateParameters{
				GALabels: wrapPointer("  "),
			},
		},
		{
			name: "single",
			want: map[string]string{"key1": "value1"},
			params: &PrivateParameters{
				GALabels: wrapPointer("key1=value1"),
			},
		},
		{
			name: "multiple",
			want: map[string]string{"key1": "value1", "key2": "value2"},
			params: &PrivateParameters{
				GALabels: wrapPointer("key1=value1,key2=value2"),
			},
		},
		{
			name: "multiple with spaces",
			want: map[string]string{"key1": "value1", "key2": "value2"},
			params: &PrivateParameters{
				GALabels: wrapPointer("key1=value1 ,  key2=value2"),
			},
		},
		{
			name: "with null",
			want: map[string]string{"key1": "value1", "key2": "null"},
			params: &PrivateParameters{
				GALabels: wrapPointer("key1=value1,key2"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.params.GetMapOfGALabels()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected result (-want +got)\n%s", diff)
			}
		})
	}
}
