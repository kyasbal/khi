package grouper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBasicGrouper(t *testing.T) {
	type inputStruct struct {
		Key   string
		Value int
	}
	testCases := []struct {
		name     string
		input    []inputStruct
		expected map[string][]inputStruct
	}{
		{
			name:     "empty input",
			input:    []inputStruct{},
			expected: map[string][]inputStruct{},
		},
		{
			name: "single group",
			input: []inputStruct{
				{
					Key:   "groupA",
					Value: 1,
				},
				{
					Key:   "groupA",
					Value: 2,
				},
			},
			expected: map[string][]inputStruct{
				"groupA": {
					{
						Key:   "groupA",
						Value: 1,
					},
					{
						Key:   "groupA",
						Value: 2,
					},
				},
			},
		},
		{
			name: "multiple groups",
			input: []inputStruct{
				{
					Key:   "groupA",
					Value: 1,
				},
				{
					Key:   "groupB",
					Value: 2,
				},
				{
					Key:   "groupA",
					Value: 3,
				},
			},
			expected: map[string][]inputStruct{
				"groupA": {
					{
						Key:   "groupA",
						Value: 1,
					},
					{
						Key:   "groupA",
						Value: 3,
					},
				},
				"groupB": {
					{
						Key:   "groupB",
						Value: 2,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grouper := NewBasicGrouper(func(input inputStruct) string {
				return input.Key
			})
			actual := grouper.Group(tc.input)
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Group() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
