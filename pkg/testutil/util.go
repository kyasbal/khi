package testutil

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

/**
* KHI has many tests using configurations and these configurations are read by glob.
* glob result can contain non meaningful relative path and slash.
 */
func GlobResultToChildTestName(filePath string, basePath string) string {
	result := strings.TrimPrefix(filePath, basePath)
	result = strings.ReplaceAll(result, "/", "_")
	return result
}

// Removes timestamps from slog logs
func RemoveSlogTimestampFromLine(log string) string {
	logs := strings.Split(log, "\n")
	for i := 0; i < len(logs); i++ {
		if logs[i] != "" {
			if logs[i][28] == 'Z' {
				logs[i] = logs[i][30:] // the first 29 chars are used for timestamps for logs with Zt
				continue
			}
			logs[i] = logs[i][35:] // the first 35 chars are used for timestamps for logs with time offset
		}
	}
	return strings.Join(logs, "\n")
}

func ResponseFromString(code int, response string) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBufferString(response)),
		StatusCode: code,
	}
}
