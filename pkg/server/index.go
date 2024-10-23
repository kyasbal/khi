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
		result = append(result, fmt.Sprintf(`<meta id="ga-meta-%s" content="%s">`, key, value))
	}
	return result
}

func getServerBasePathMetaTag() string {
	serverBasePath := os.Getenv("KHI_SERVER_BASE_PATH")
	if serverBasePath == "" {
		serverBasePath = "/"
	}
	if serverBasePath[len(serverBasePath)-1] != '/' {
		serverBasePath += "/"
	}
	return fmt.Sprintf(`<meta id="server-base-path" content="%s">`, serverBasePath)
}

// getBaseTag returns the `<base>` tag to rewrite the base url of resources accessed with relative path on frontend.
func getBaseTag() string {
	frontendResourceBasePath := os.Getenv("KHI_FRONTEND_RESOURCE_BASE_PATH")
	if frontendResourceBasePath == "" && os.Getenv("KHI_SERVER_BASE_PATH") != "" {
		frontendResourceBasePath = os.Getenv("KHI_SERVER_BASE_PATH")
	}
	if frontendResourceBasePath == "" {
		frontendResourceBasePath = "/"
	}
	if frontendResourceBasePath[len(frontendResourceBasePath)-1] != '/' {
		frontendResourceBasePath += "/"
	}
	return fmt.Sprintf(`<base href="%s">`, frontendResourceBasePath)
}

func injectSnipptToHtml(originalHtml string, plactholder string, snippit string) string {
	return strings.ReplaceAll(originalHtml, plactholder, snippit)
}

// replaceLocalDevServerOnlyTag removed tags only used in the local dev environment.
func replaceLocalDevServerOnlyTag(source string) string {
	return strings.ReplaceAll(source, `<base href="/" />`, "") // The base tag must be supplied but it will be injected from backend usually. It needs to be added by default on the static index.html because KHI can't inject the tag to the static index.html served by angular dev server.
}

func generateIndexHtmlWithGALabels(originalPath string) string {
	replaceTarget := `<!--INJECT GENERATED CODE HERE FROM BACKEND-->`
	data, err := os.ReadFile(originalPath)
	if err != nil {
		panic(err)
	}

	if !strings.Contains(string(data), replaceTarget) {
		panic("Inject target string was not found!")
	}
	gaLabels := getCommaSeperatedKVPairEnv("KHI_GA_LABELS")

	var injectedTags []string
	injectedTags = append(injectedTags, getBaseTag())
	injectedTags = append(injectedTags, generateGaMetaTags(gaLabels)...)
	injectedTags = append(injectedTags, getServerBasePathMetaTag())

	indexSource := string(data)

	return injectSnipptToHtml(replaceLocalDevServerOnlyTag(indexSource), replaceTarget, strings.Join(injectedTags, "\n"))
}
