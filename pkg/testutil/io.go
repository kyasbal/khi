package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func InitTestIO() {
	// Adjust current directory to be root of the repository
	ROOT_INDICATOR_FILE := ".root"
	for {
		cwd, err := os.Getwd()
		if err != nil {
			panic("Getwd failed to adjust working directory")
		}
		if _, err := os.Stat(filepath.Join(cwd, ROOT_INDICATOR_FILE)); err == nil {
			break
		}
		pathsSegments := strings.Split(cwd, "/")
		nextCwd := filepath.Join(pathsSegments[:len(pathsSegments)-1]...)
		err = os.Chdir("/" + nextCwd)
		if err != nil {
			panic("Failed to change current directory to " + nextCwd)
		}
	}
}

func MustReadText(filePath string) string {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func MustReadYaml(filePath string) map[string]any {
	result := map[string]any{}
	err := yaml.Unmarshal([]byte(MustReadText(filePath)), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func MustParseTimeRFC3339(timeStr string) time.Time {
	time, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}
	return time
}

// MustPlaceTemporalFile write a file at the specified path and returns a function to unlink it for cleaning up.
func MustPlaceTemporalFile(filePath, content string) (cleanup func()) {
	buf, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_, err = buf.WriteString(content)
	if err != nil {
		panic(err)
	}
	return func() {
		err := os.Remove(filePath)
		if err != nil {
			panic(err)
		}
	}
}

func GlobTestResources(fileGlob string, ignoredSuffixes []string) []string {
	matches, err := filepath.Glob(fileGlob)
	if err != nil {
		panic(err)
	}
	var result []string = make([]string, 0)
	for _, match := range matches {
		ignored := false
		for _, ignoreSuffix := range ignoredSuffixes {
			if strings.HasSuffix(match, ignoreSuffix) {
				ignored = true
				break
			}
		}
		if !ignored {
			result = append(result, match)
		}
	}
	return matches
}
