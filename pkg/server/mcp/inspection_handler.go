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

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// InspectionHandler handles MCP tools and resources related to inspections.
type InspectionHandler struct {
	inspectionServer *coreinspection.InspectionTaskServer
}

var _ DomainHandler = (*InspectionHandler)(nil)

// NewInspectionHandler creates a new InspectionHandler instance.
func NewInspectionHandler(inspectionServer *coreinspection.InspectionTaskServer) *InspectionHandler {
	return &InspectionHandler{
		inspectionServer: inspectionServer,
	}
}

// Register registers all tools and resources provided by the inspection handler onto the MCP server.
func (h *InspectionHandler) Register(srv *mcp.Server) {
	h.RegisterTools(srv)
	h.RegisterResources(srv)
}

// ListInspections collects and returns all active and completed inspections as InspectionSummaryItems.
func (h *InspectionHandler) ListInspections() ([]*InspectionSummaryItem, error) {
	if h.inspectionServer == nil {
		return []*InspectionSummaryItem{}, nil
	}

	runners := h.inspectionServer.GetAllRunners()
	items := make([]*InspectionSummaryItem, 0, len(runners))
	for _, runner := range runners {
		if !runner.Started() {
			continue
		}
		item, err := runnerToInspectionSummary(runner)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	slices.SortFunc(items, func(a, b *InspectionSummaryItem) int {
		return strings.Compare(a.ID, b.ID)
	})

	return items, nil
}

// GetInspectionSummary generates and returns a DetailedInspectionSummary for the given inspection ID.
func (h *InspectionHandler) GetInspectionSummary(inspectionID string) (*DetailedInspectionSummary, error) {
	if h.inspectionServer == nil {
		return nil, fmt.Errorf("inspection server is not available")
	}

	runner := h.inspectionServer.GetInspection(inspectionID)
	if runner == nil {
		for _, r := range h.inspectionServer.GetAllRunners() {
			if r.ID == inspectionID {
				runner = r
				break
			}
		}
	}

	if runner == nil {
		return nil, fmt.Errorf("inspection %q was not found", inspectionID)
	}

	summaryItem, err := runnerToInspectionSummary(runner)
	if err != nil {
		return nil, fmt.Errorf("failed to get inspection summary: %w", err)
	}

	md, err := runner.GetCurrentMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for inspection %q: %w", inspectionID, err)
	}

	var errorMessages []string
	if errorSet, found := typedmap.Get(md, inspectionmetadata.ErrorMessageSetMetadataKey); found && errorSet != nil {
		for _, errMsg := range errorSet.ErrorMessages {
			if errMsg != nil && errMsg.Message != "" {
				errorMessages = append(errorMessages, errMsg.Message)
			}
		}
	}

	variables := make(map[string]any)
	variables["InspectionName"] = summaryItem.Name
	variables["InspectionType"] = summaryItem.Type
	variables["Status"] = summaryItem.Status
	variables["StartTime"] = summaryItem.StartTime
	variables["EndTime"] = summaryItem.EndTime
	variables["FileSize"] = summaryItem.FileSize
	variables["ErrorCount"] = summaryItem.ErrorCount
	if len(errorMessages) > 0 {
		variables["Errors"] = errorMessages
	}

	markdown := buildDefaultMarkdownSummary(summaryItem, errorMessages)

	return &DetailedInspectionSummary{
		Inspection:    *summaryItem,
		Markdown:      markdown,
		ErrorMessages: errorMessages,
		Variables:     variables,
	}, nil
}

// RegisterTools registers inspection-related tools onto the MCP server.
func (h *InspectionHandler) RegisterTools(srv *mcp.Server) {
	listInspectionsTool := &mcp.Tool{
		Name:        "list_inspections",
		Description: "Lists all active and completed Kubernetes inspection datasets currently available in KHI.",
	}
	mcp.AddTool(srv, listInspectionsTool, func(ctx context.Context, req *mcp.CallToolRequest, input ListInspectionsInput) (*mcp.CallToolResult, *ListInspectionsOutput, error) {
		items, err := h.ListInspections()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list inspections: %w", err)
		}
		return nil, &ListInspectionsOutput{Inspections: items}, nil
	})

	getInspectionSummaryTool := &mcp.Tool{
		Name:        "get_inspection_summary",
		Description: "Retrieves a detailed markdown summary and metadata for a specific Kubernetes inspection dataset.",
	}
	mcp.AddTool(srv, getInspectionSummaryTool, func(ctx context.Context, req *mcp.CallToolRequest, input GetInspectionSummaryInput) (*mcp.CallToolResult, *GetInspectionSummaryOutput, error) {
		if input.InspectionID == "" {
			return nil, nil, fmt.Errorf("inspectionId is required")
		}
		summary, err := h.GetInspectionSummary(input.InspectionID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get inspection summary: %w", err)
		}
		return nil, &GetInspectionSummaryOutput{Summary: summary}, nil
	})
}

// RegisterResources registers inspection-related resources and templates onto the MCP server.
func (h *InspectionHandler) RegisterResources(srv *mcp.Server) {
	inspectionsResource := &mcp.Resource{
		URI:         "khi://inspections",
		Name:        "Inspection List",
		Description: "List of available inspection datasets on the KHI server.",
		MIMEType:    "application/json",
	}
	srv.AddResource(inspectionsResource, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		items, err := h.ListInspections()
		if err != nil {
			return nil, fmt.Errorf("failed to list inspections: %w", err)
		}
		data, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal inspection list: %w", err)
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "khi://inspections",
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})

	summaryResourceTemplate := &mcp.ResourceTemplate{
		URITemplate: "khi://inspections/{id}/summary",
		Name:        "Inspection Summary",
		Description: "Markdown summary of an inspection dataset on the KHI server.",
		MIMEType:    "text/markdown",
	}
	srv.AddResourceTemplate(summaryResourceTemplate, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		uri := req.Params.URI
		inspectionID := extractInspectionIDFromURI(uri)
		if inspectionID == "" {
			return nil, fmt.Errorf("invalid inspection summary URI: %s", uri)
		}

		summary, err := h.GetInspectionSummary(inspectionID)
		if err != nil {
			return nil, fmt.Errorf("failed to read inspection summary: %w", err)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      uri,
					MIMEType: "text/markdown",
					Text:     summary.Markdown,
				},
			},
		}, nil
	})
}

func extractInspectionIDFromURI(uri string) string {
	trimmed := strings.TrimPrefix(uri, "khi://inspections/")
	return strings.TrimSuffix(trimmed, "/summary")
}

func buildDefaultMarkdownSummary(item *InspectionSummaryItem, errorMessages []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Inspection Summary: %s\n\n", item.Name))
	sb.WriteString(fmt.Sprintf("- **ID**: `%s`\n", item.ID))
	sb.WriteString(fmt.Sprintf("- **Type**: %s\n", item.Type))
	sb.WriteString(fmt.Sprintf("- **Status**: %s\n", item.Status))
	if item.StartTime != "" && item.EndTime != "" {
		sb.WriteString(fmt.Sprintf("- **Time Range**: %s - %s\n", item.StartTime, item.EndTime))
	}
	if item.FileSize > 0 {
		sb.WriteString(fmt.Sprintf("- **File Size**: %d bytes\n", item.FileSize))
	}
	sb.WriteString(fmt.Sprintf("- **Errors Recorded**: %d\n\n", item.ErrorCount))

	sb.WriteString("## Errors & Inspection Warnings\n")
	if len(errorMessages) > 0 {
		for _, msg := range errorMessages {
			sb.WriteString(fmt.Sprintf("- %s\n", msg))
		}
	} else {
		sb.WriteString("- No inspection errors encountered.\n")
	}

	return sb.String()
}

func runnerToInspectionSummary(runner *coreinspection.InspectionTaskRunner) (*InspectionSummaryItem, error) {
	md, err := runner.GetCurrentMetadata()
	if err != nil {
		return nil, err
	}

	item := &InspectionSummaryItem{
		ID:     runner.ID,
		Status: "UNSPECIFIED",
	}

	if header, found := typedmap.Get(md, inspectionmetadata.HeaderMetadataKey); found && header != nil {
		item.Name = header.InspectionName
		item.Type = header.InspectionType
		item.FileSize = int64(header.FileSize)
		if header.StartTimeUnixSeconds > 0 {
			item.StartTime = time.Unix(header.StartTimeUnixSeconds, 0).UTC().Format(time.RFC3339)
		}
		if header.EndTimeUnixSeconds > 0 {
			item.EndTime = time.Unix(header.EndTimeUnixSeconds, 0).UTC().Format(time.RFC3339)
		}
		if header.InspectTimeUnixSeconds > 0 {
			item.InspectTime = time.Unix(header.InspectTimeUnixSeconds, 0).UTC().Format(time.RFC3339)
		}
	}

	if progress, found := typedmap.Get(md, inspectionmetadata.ProgressMetadataKey); found && progress != nil {
		item.Status = string(progress.Phase)
	}

	if errorSet, found := typedmap.Get(md, inspectionmetadata.ErrorMessageSetMetadataKey); found && errorSet != nil {
		item.ErrorCount = len(errorSet.ErrorMessages)
	}

	return item, nil
}
