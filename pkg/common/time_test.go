package common

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	JST := time.FixedZone("JST", 9*60*60)
	testCases := []struct {
		Input    string
		Expected time.Time
		Error    bool
	}{
		{
			Input:    "2023-01-02T03:04:05Z",
			Expected: time.Date(2023, time.January, 2, 3, 4, 5, 0, time.UTC),
		},
		{
			Input:    "2023-01-02T03:04:05+00:00",
			Expected: time.Date(2023, time.January, 2, 3, 4, 5, 0, time.UTC),
		},
		{
			Input:    "2023-01-02T03:04:05+09:00",
			Expected: time.Date(2023, time.January, 2, 3, 4, 5, 0, JST),
		},
		{
			Input: "2023-01-02T03:04:05+09:00-invalid-string",
			Error: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("testcase-%s", testCase.Input), func(t *testing.T) {
			result, err := ParseTime(testCase.Input)
			if testCase.Error {
				if err == nil {
					t.Errorf("expect the call ending with an error. But no error returned.")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error was returned\n%s", err)
				}
				if result.Unix() != testCase.Expected.Unix() {
					t.Errorf("the result is not matching\nexpect:\n%s\nactual:\n%s\n", testCase.Expected, result)
				}
			}
		})
	}
}
