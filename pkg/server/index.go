package server

// Codes to generate index.html dynamatically

import (
	"fmt"
	"os"
	"strings"
)

// getCommaSeperatedKVPairEnv generates key value pairs from specified environment name
func getCommaSeperatedKVPairEnv(envName string) map[string]string {
	result := make(map[string]string)

	if env, hasEnv := os.LookupEnv(envName); hasEnv {
		keyValuePairs := strings.Split(env, ",")
		for _, pair := range keyValuePairs {
			keyValues := strings.Split(pair, "=")
			key := keyValues[0]
			value := "null"
			if len(keyValues) > 1 {
				value = keyValues[1]
			}
			result[key] = value
		}
	}

	return result
}

func generateGaMetaTags(gaLabels map[string]string) []string {
	result := make([]string, 0)
	for key, value := range gaLabels {
		result = append(result, fmt.Sprintf("<meta id=\"ga-meta-%s\" content=\"%s\">", key, value))
	}

	return result
}

func injectSnipptToHtml(originalHtml string, plactholder string, snippit string) string {
	return strings.ReplaceAll(originalHtml, plactholder, snippit)
}

func generateIndexHtmlWithGALabels(originalPath string) string {
	replaceTarget := "<!--INJECT GENERATED CODE HERE FROM BACKEND-->"
	data, err := os.ReadFile(originalPath)
	if err != nil {
		panic(err)
	}

	if !strings.Contains(string(data), replaceTarget) {
		panic("Inject target string was not found!")
	}
	gaLabels := getCommaSeperatedKVPairEnv("KHI_GA_LABELS")
	gaMetaTags := generateGaMetaTags(gaLabels)

	return injectSnipptToHtml(string(data), replaceTarget, strings.Join(gaMetaTags, "\n"))
}
