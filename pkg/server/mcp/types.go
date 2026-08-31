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

package mcp

// InspectionSummaryItem represents high-level metadata for an inspection returned to MCP clients.
type InspectionSummaryItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	InspectTime string `json:"inspectTime,omitempty"`
	FileSize    int64  `json:"fileSize,omitempty"`
	ErrorCount  int    `json:"errorCount"`
}

// DetailedInspectionSummary provides full markdown summary and metadata context for an inspection.
type DetailedInspectionSummary struct {
	Inspection    InspectionSummaryItem `json:"inspection"`
	Markdown      string                `json:"markdown"`
	ErrorMessages []string              `json:"errorMessages,omitempty"`
	Variables     map[string]any        `json:"variables,omitempty"`
}

// ListInspectionsInput is the empty input parameter schema for list_inspections.
type ListInspectionsInput struct{}

// ListInspectionsOutput is the output parameter schema for list_inspections.
type ListInspectionsOutput struct {
	Inspections []*InspectionSummaryItem `json:"inspections"`
}

// GetInspectionSummaryInput is the input parameter schema for get_inspection_summary.
type GetInspectionSummaryInput struct {
	// The unique inspection ID to get summary for.
	InspectionID string `json:"inspectionId" jsonschema:"required"`
}

// GetInspectionSummaryOutput is the output parameter schema for get_inspection_summary.
type GetInspectionSummaryOutput struct {
	Summary *DetailedInspectionSummary `json:"summary"`
}

// OpenWorkbenchInput is the input parameter schema for open_workbench.
type OpenWorkbenchInput struct {
	// The unique inspection ID to open in the active workbench session.
	InspectionID string `json:"inspectionId" jsonschema:"required"`
}

// OpenWorkbenchOutput is the output parameter schema for open_workbench.
type OpenWorkbenchOutput struct {
	WorkbenchID  string `json:"workbenchId"`
	InspectionID string `json:"inspectionId"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

// WorkbenchStatusOutput is the output parameter schema for workbench status resource.
type WorkbenchStatusOutput struct {
	Active       bool   `json:"active"`
	WorkbenchID  string `json:"workbenchId,omitempty"`
	InspectionID string `json:"inspectionId,omitempty"`
	IndexState   string `json:"indexState,omitempty"`
}
