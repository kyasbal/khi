package parameters

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDebugParameters(t *testing.T) {
	testCases := []struct {
		name   string
		want   *DebugParameters
		before func()
	}{
		{
			before: func() {
				os.Args = []string{os.Args[0]}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
			name: "default",
			want: &DebugParameters{
				Profiler:         wrapPointer(false),
				ProfilerService:  wrapPointer("khi"),
				ProfilerProject:  wrapPointer(""),
				DisableAnalytics: wrapPointer(false),
				AnalyticsDebug:   wrapPointer(false),
				Verbose:          wrapPointer(false),
				NoColor:          wrapPointer(false),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := &DebugParameters{}
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
