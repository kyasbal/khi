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

package types

import "time"

type AnalyticsEvent = string

const (
	AnalyticsEventKHIStart                   = "khi-start"
	AnalyticsEventKHITerminate               = "khi-terminate"
	AnalyticsEventInspectionStart            = "inspection-start"
	AnalyticsEventInspectionEnd              = "inspection-end"
	AnalyticsEventFrontendInit               = "INIT"
	AnalyticsEventFrontendInspect            = "INSPECT"
	AnalyticsEventFrontendOpenInspectionData = "OPEN_INSPECTION_DATA"
)

type RecordAnalyticsDataRequest struct {
	Event    AnalyticsEvent `json:"event"`
	Debug    bool           `json:"debug"`
	Metadata any            `json:"metadata"`
}

type RecordAnalyticsData struct {
	Event    AnalyticsEvent `bigquery:"event"`
	Time     time.Time      `bigquery:"time"`
	Debug    bool           `bigquery:"debug"`
	Metadata string         `bigquery:"metadata"`
}
