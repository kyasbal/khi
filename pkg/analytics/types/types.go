package types

import "time"

type AnalyticsEvent = string

const (
	AnalyticsEventKHIStart        = "khi-start"
	AnalyticsEventKHITerminate    = "khi-terminate"
	AnalyticsEventInspectionStart = "inspection-start"
	AnalyticsEventInspectionEnd   = "inspection-end"
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
