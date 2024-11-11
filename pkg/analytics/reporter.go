// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/analytics/types"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
)

var analyticsEndpoint = "https://khi-analytics-5dxktvcd7q-uc.a.run.app"

type AnalyticsReporter struct {
	debug          bool
	GlobalMetadata map[string]string
}

func NewAnalyticsReporter() *AnalyticsReporter {
	debug := false
	if parameters.Debug.AnalyticsDebug != nil {
		debug = *parameters.Debug.AnalyticsDebug
	}
	gaMetadata := map[string]string{}
	if parameters.Private.GALabels != nil {
		gaMetadata = parameters.Private.GetMapOfGALabels()
	}
	return &AnalyticsReporter{
		debug:          debug,
		GlobalMetadata: gaMetadata,
	}
}

func (r *AnalyticsReporter) ReportEvent(event types.AnalyticsEvent, metadata map[string]any) {
	if parameters.Debug.DisableAnalytics == nil || *parameters.Debug.DisableAnalytics {
		return
	}
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
