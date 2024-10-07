package testlog

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBaseYamlTestLogOpt(t *testing.T) {
	testCases := []struct {
		name        string
		inputYaml   string
		outputYaml  string
		expectError bool
	}{
		{
			name:      "basic valid yaml",
			inputYaml: `foo: bar`,
			outputYaml: `foo: bar
`,
			expectError: false,
		},
		{
			name:        "parses empty yaml as an empty map",
			inputYaml:   "",
			outputYaml:  "{}\n",
			expectError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tl := New(BaseYaml(tc.inputYaml))
			reader, err := tl.BuildReader()
			if tc.expectError {
				if err == nil {
					t.Errorf("Expecting an error but no error returned.")
				}
			} else {
				yamlStr, err := reader.ToYaml("")
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if diff := cmp.Diff(yamlStr, tc.outputYaml); diff != "" {
					t.Errorf("Yaml string mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
