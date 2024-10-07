package k8s_container

import (
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

var MainMessageSeverityParsers = []MainMessageSeverityParser{
	&MetricsContainerLogSeverityParser{},
}

// MainMessageSeverityParser is used for MainMessage extracted from a log entry.
// These parsers are only used when the structured log itself didn't contain severity info and KHI need to read severtiy from the content of log.
type MainMessageSeverityParser interface {
	TryParse(message string) enum.Severity
}

func ParseSeverity(message string) enum.Severity {
	for _, parser := range MainMessageSeverityParsers {
		severity := parser.TryParse(message)
		if severity != enum.SeverityUnknown {
			return severity
		}
	}
	return enum.SeverityUnknown
}

type MetricsContainerLogSeverityParser struct{}

// TryParse implements MainMessageSeverityParser.
func (m *MetricsContainerLogSeverityParser) TryParse(message string) enum.Severity {
	fragments := strings.Split(message, "\t")
	if len(fragments) < 2 {
		return enum.SeverityUnknown
	}
	_, err := time.Parse(time.RFC3339, fragments[0])
	if err != nil {
		return enum.SeverityUnknown
	}
	severityStr := fragments[1]
	if severityStr == "info" {
		return enum.SeverityInfo
	} else if severityStr == "warn" {
		return enum.SeverityWarning
	} else if severityStr == "error" {
		return enum.SeverityError
	} else {
		return enum.SeverityUnknown
	}
}

var _ (MainMessageSeverityParser) = (*MetricsContainerLogSeverityParser)(nil)
