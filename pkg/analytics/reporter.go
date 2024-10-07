package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/analytics/types"
)

var analyticsEndpoint = "https://khi-analytics-5dxktvcd7q-uc.a.run.app"

type AnalyticsReporter struct {
	debug          bool
	GlobalMetadata map[string]any
}

// This will be deprecated in the future. KHI_GA_LABELS are used for frontend initially, we should define new environment variables for analytics.
// But for now, we will use the old environment variable in the transition time.
func getCommaSeperatedKVPairEnv() map[string]any {
	result := make(map[string]any)
	if env, hasEnv := os.LookupEnv("KHI_GA_LABELS"); hasEnv {
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

func NewAnalyticsReporter() *AnalyticsReporter {
	analyticsDebug, found := os.LookupEnv("KHI_ANALYTICS_DEBUG")
	isDebug := false
	if found && strings.ToLower(analyticsDebug) != "false" {
		isDebug = true
	}
	return &AnalyticsReporter{
		debug:          isDebug,
		GlobalMetadata: getCommaSeperatedKVPairEnv(),
	}
}

func (r *AnalyticsReporter) ReportEvent(event types.AnalyticsEvent, metadata map[string]any) {
	mergedMetadata := map[string]any{}
	for globalKey, globalValue := range r.GlobalMetadata {
		mergedMetadata[globalKey] = globalValue
	}
	for localKey, localValue := range metadata {
		mergedMetadata[localKey] = localValue
	}

	request := types.RecordAnalyticsDataRequest{
		Debug:    r.debug,
		Event:    event,
		Metadata: mergedMetadata,
	}

	marshalled, err := json.Marshal(request)
	if err != nil {
		slog.Warn(fmt.Sprintf("Failed to generate json string from the request\n%s", err))
		return
	}

	req, err := http.NewRequest("POST", analyticsEndpoint, bytes.NewReader(marshalled))
	if err != nil {
		slog.Warn(fmt.Sprintf("Failed to report the usage data\n%s", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn(fmt.Sprintf("Failed to report the usage data\n%d,%s\n%s", resp.StatusCode, resp.Status, err))
		return
	}
}
