package log

import (
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

// CommonLogFieldExtractor extracts information being available for all log entries(e.g timestamp)..etc
// These fields can be used outside of parsers specific to log types.
type CommonLogFieldExtractor interface {
	// Get unique ID from a log. Unique but shorter id is preferable.
	ID(l *LogEntity) string
	// Get the timestamp from a log. This is used for sorting.
	Timestamp(l *LogEntity) time.Time
	// Extract the main content of the structured logging.
	// textPayload or jsonPayload.MESSAGE in GCP
	MainMessage(l *LogEntity) (string, error)
	// Severity of this log.
	Severity(l *LogEntity) (enum.Severity, error)
	// ID visible to user. This value isn't necessary to be unique.
	DisplayID(l *LogEntity) string
	// Entire log body represented as a string.
	LogBody(l *LogEntity) string
}
