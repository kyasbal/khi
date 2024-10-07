package testutil

import (
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveSlogTimestampFromLine(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "time=2024-02-08T10:48:10.459+09:00 level=INFO msg=\"inspection1 task1 info\"",
			expected: "level=INFO msg=\"inspection1 task1 info\"",
		},
		{
			input:    "time=2024-02-08T10:48:10.459Z level=INFO msg=\"inspection1 task1 info\"",
			expected: "level=INFO msg=\"inspection1 task1 info\"",
		},
		{
			input:    "", // Empty input
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			result := RemoveSlogTimestampFromLine(testCase.input)
			if result != testCase.expected {
				t.Errorf("RemoveSlogTimestampFromLine failed. Input: %s, Expected: %s, Got: %s", testCase.input, testCase.expected, result)
			}
		})
	}
}

func TestResponseFromString(t *testing.T) {
	t.Parallel()
	type args struct {
		code     int
		response string
	}
	tests := []struct {
		name     string
		args     args
		wantBody string
		wantCode int
	}{
		{
			name: "ResponseFromString should return a http.Response with the given code and response",
			args: args{
				code:     http.StatusOK,
				response: "test-response",
			},
			wantBody: "test-response",
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResponseFromString(tt.args.code, tt.args.response)
			if got.StatusCode != tt.wantCode {
				t.Errorf("ResponseFromString() = %v, want %v", got.StatusCode, tt.wantCode)
			}
			gotBody, _ := io.ReadAll(got.Body)
			if diff := cmp.Diff(string(gotBody), tt.wantBody); diff != "" {
				t.Errorf("ResponseFromString() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
