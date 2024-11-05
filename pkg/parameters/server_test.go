package parameters

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestServerParameters(t *testing.T) {
	testCases := []struct {
		name   string
		want   *ServerParameters
		before func()
	}{
		{
			before: func() {
				os.Args = []string{os.Args[0]}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
			name: "default",
			want: &ServerParameters{
				ViewerMode:               wrapPointer(false),
				Port:                     wrapPointer(8080),
				Host:                     wrapPointer("localhost"),
				BasePath:                 wrapPointer("/"),
				FrontendResourceBasePath: wrapPointer("/"),
				FrontendAssetFolder:      wrapPointer("./web"),
			},
		},
		{
			before: func() {
				os.Args = []string{os.Args[0], "--base-path", "/foo/bar"}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
			name: "FrontendResourceBasePath uses BasePath when not set",
			want: &ServerParameters{
				ViewerMode:               wrapPointer(false),
				Port:                     wrapPointer(8080),
				Host:                     wrapPointer("localhost"),
				BasePath:                 wrapPointer("/foo/bar/"),
				FrontendResourceBasePath: wrapPointer("/foo/bar/"),
				FrontendAssetFolder:      wrapPointer("./web"),
			},
		},
		{
			before: func() {
				os.Args = []string{os.Args[0], "--base-path", "/foo/bar", "--frontend-resource-base-path", "/foo"}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
			name: "FrontendResourceBasePath should complement the last /",
			want: &ServerParameters{
				ViewerMode:               wrapPointer(false),
				Port:                     wrapPointer(8080),
				Host:                     wrapPointer("localhost"),
				BasePath:                 wrapPointer("/foo/bar/"),
				FrontendResourceBasePath: wrapPointer("/foo/"),
				FrontendAssetFolder:      wrapPointer("./web"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := &ServerParameters{}
			tc.before()
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
