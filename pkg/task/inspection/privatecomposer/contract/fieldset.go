// Copyright 2026 Google LLC
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

package privatecomposer_contract

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
)

var postgresLogPrefixes = []string{
	"PANIC:",
	"FATAL:",
	"ERROR:",
	"WARNING:",
	"LOG:",
	"INFO:",
	"DETAIL:",
	"HINT:",
}

// CloudSQLFieldSet holds structured log data extracted from Cloud SQL entries.
type CloudSQLFieldSet struct {
	DatabaseID  string
	LogFileName string
	Summary     string
}

// Kind returns the unique identifier of this FieldSet.
func (f *CloudSQLFieldSet) Kind() string {
	return "privatecomposer.khi.google.com/CloudSQLFieldSet"
}

var _ log.FieldSet = (*CloudSQLFieldSet)(nil)

// CloudSQLFieldSetReader extracts CloudSQLFieldSet from a raw log node reader.
type CloudSQLFieldSetReader struct{}

// FieldSetKind returns the Kind string of CloudSQLFieldSet.
func (r *CloudSQLFieldSetReader) FieldSetKind() string {
	return "privatecomposer.khi.google.com/CloudSQLFieldSet"
}

// Read parses fields from structured Cloud SQL log data.
func (r *CloudSQLFieldSetReader) Read(reader *structured.NodeReader) (log.FieldSet, error) {
	rawDatabaseID := reader.ReadStringOrDefault("resource.labels.database_id", "")
	databaseID := extractDatabaseID(rawDatabaseID)

	logName := reader.ReadStringOrDefault("logName", "")
	logFileName := extractLogFileName(logName)

	var summary string
	if reader.Has("textPayload") {
		rawMessage := reader.ReadStringOrDefault("textPayload", "")
		summary = extractPostgresSummary(rawMessage)
	} else if reader.Has("protoPayload") {
		methodName := reader.ReadStringOrDefault("protoPayload.methodName", "")
		principalEmail := reader.ReadStringOrDefault("protoPayload.authenticationInfo.principalEmail", "")
		summary = extractAuditSummary(methodName, principalEmail)
	}

	return &CloudSQLFieldSet{
		DatabaseID:  databaseID,
		LogFileName: logFileName,
		Summary:     summary,
	}, nil
}

var _ log.FieldSetReader = (*CloudSQLFieldSetReader)(nil)

// extractLogFileName extracts the log file name from a Cloud Logging logName resource path.
func extractLogFileName(logName string) string {
	if logName == "" {
		return "unknown"
	}
	if idx := strings.LastIndex(logName, "%2F"); idx >= 0 {
		return logName[idx+3:]
	}
	if idx := strings.LastIndex(logName, "/"); idx >= 0 {
		return logName[idx+1:]
	}
	return logName
}

// extractPostgresSummary extracts the log summary after PostgreSQL severity prefixes.
func extractPostgresSummary(payload string) string {
	bestIdx := -1
	bestPrefixLen := 0
	for _, prefix := range postgresLogPrefixes {
		idx := strings.Index(payload, prefix)
		if idx >= 0 {
			if bestIdx == -1 || idx < bestIdx {
				bestIdx = idx
				bestPrefixLen = len(prefix)
			}
		}
	}
	if bestIdx >= 0 {
		return strings.TrimSpace(payload[bestIdx+bestPrefixLen:])
	}
	return strings.TrimSpace(payload)
}

// extractAuditSummary formats the log summary for Cloud SQL audit log entries.
func extractAuditSummary(methodName, principalEmail string) string {
	if principalEmail != "" && principalEmail != "unknown" {
		if methodName != "" {
			return fmt.Sprintf("%s from %s", methodName, principalEmail)
		}
		return fmt.Sprintf("from %s", principalEmail)
	}
	return methodName
}

// extractDatabaseID extracts the database instance ID from a Cloud SQL database_id label.
// Cloud SQL database_id typically takes the format "<project-id>:<instance-id>".
// Since the timeline path is already structured under the GCP project ID, the project prefix is stripped.
func extractDatabaseID(databaseID string) string {
	if databaseID == "" {
		return "unknown"
	}
	if _, instanceID, ok := strings.Cut(databaseID, ":"); ok {
		if instanceID == "" {
			return "unknown"
		}
		return instanceID
	}
	return databaseID
}
