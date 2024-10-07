package logger

import "log/slog"

var LogKindAttrKey = "log-kind"

// LogKind returns slog.Attr marks which log kind is associated to the log record.
func LogKind(lt string) slog.Attr {
	return slog.String(LogKindAttrKey, lt)
}
