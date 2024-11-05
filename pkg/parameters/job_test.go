package parameters

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJobParameters(t *testing.T) {
	testCases := []struct {
		name   string
		want   *JobParameters
		before func()
	}{
		{
			name: "default",
			want: &JobParameters{
				JobMode:            wrapPointer(false),
				InspectionType:     wrapPointer(""),
				InspectionFeatures: wrapPointer(""),
				InspectionValues:   wrapPointer(""),
				ExportDestination:  wrapPointer(""),
			},
			before: func() {
				os.Args = []string{os.Args[0]}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			store := &JobParameters{}
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
