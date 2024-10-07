package schema

// Log severity used in KHI.
// There would be more various severity depending on the log type.
// But KHI only has these 4 different type and each parser should change the severity to them if the original severity was not in there.
type KHILogSeverity = string

const SeverityInfo = "INFO"
const SeverityWarn = "WARN"
const SeverityError = "ERROR"
const SeverityFatal = "FATAL"
const SeverityUnknown = "UNKNOWN"
