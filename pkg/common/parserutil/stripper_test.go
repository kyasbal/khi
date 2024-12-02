package parserutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStripSpecialSequences(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "strip nothing",
			input:    "this is text",
			expected: "this is text",
		},
		{
			name:     "strip escape sequences",
			input:    "this is\\r\\n text\\r\\n",
			expected: "this is text",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			breaklineStripper := NewSequenceStripper("\\r", "\\n")
			actual := StripSpecialSequences(tc.input, breaklineStripper)
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("the result is not matching with the expected result\n%s", diff)
			}
		})
	}
}

func TestANSIEscapeSequenceStripper(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "strip nothing",
			input:    "this is text",
			expected: "this is text",
		},
		{
			name:     "strip ansi escape sequences",
			input:    "\\x1b[31mthis is red text\\x1b[0m",
			expected: "this is red text",
		},
		{
			name:     "strip ansi escape sequences with multiple begin sequences",
			input:    "\\u001B[31mthis is red text\\033[0m",
			expected: "this is red text",
		},
		{
			name:     "strip ansi escape sequences with incomplete sequence",
			input:    "\\x1b[31mthis is red text\\x1b[",
			expected: "this is red text\\x1b[",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stripper := NewANSIEscapeSequenceStripper()
			actual := stripper.Strip(tc.input)
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("the result is not matching with the expected result\n%s", diff)
			}
		})
	}
}
