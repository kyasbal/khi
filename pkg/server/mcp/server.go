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
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

// ListInspectionsInput is the empty input parameter schema for list_inspections.
type ListInspectionsInput struct{}

// ListInspectionsOutput is the output parameter schema for list_inspections.
type ListInspectionsOutput struct {
	Inspections []*InspectionSummaryItem `json:"inspections"`
}

// Server encapsulates the MCP server and its transport bindings.
type Server struct {
	inspectionServer  *coreinspection.InspectionTaskServer
	mcpServer         *mcp.Server
	streamableHandler *mcp.StreamableHTTPHandler
}

// NewServer creates a new KHI MCP server and registers tools and resources using the official MCP Go SDK.
func NewServer(inspectionServer *coreinspection.InspectionTaskServer) *Server {
	impl := &mcp.Implementation{
		Name:    "khi-mcp-server",
		Version: "1.0.0",
	}
	mcpSrv := mcp.NewServer(impl, nil)

	s := &Server{
		inspectionServer: inspectionServer,
		mcpServer:        mcpSrv,
	}

	s.registerTools()
	s.registerResources()

	s.streamableHandler = mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return s.mcpServer
	}, nil)

	return s
}

// MCPServer returns the underlying MCP server instance.
func (s *Server) MCPServer() *mcp.Server {
	return s.mcpServer
}

// HTTPHandler returns an http.Handler that handles Streamable HTTP MCP sessions.
func (s *Server) HTTPHandler() http.Handler {
	return s.streamableHandler
}

// ListInspections collects and returns all active and completed inspections as InspectionSummaryItems.
func (s *Server) ListInspections() ([]*InspectionSummaryItem, error) {
	if s.inspectionServer == nil {
		return []*InspectionSummaryItem{}, nil
	}

	runners := s.inspectionServer.GetAllRunners()
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

func (s *Server) registerTools() {
	listInspectionsTool := &mcp.Tool{
		Name:        "list_inspections",
		Description: "Lists all active and completed Kubernetes inspection datasets currently available in KHI.",
	}

	mcp.AddTool(s.mcpServer, listInspectionsTool, func(ctx context.Context, req *mcp.CallToolRequest, input ListInspectionsInput) (*mcp.CallToolResult, *ListInspectionsOutput, error) {
		items, err := s.ListInspections()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list inspections: %w", err)
		}
		return nil, &ListInspectionsOutput{Inspections: items}, nil
	})
}

func (s *Server) registerResources() {
	inspectionsResource := &mcp.Resource{
		URI:         "khi://inspections",
		Name:        "Inspection List",
		Description: "List of available inspection datasets on the KHI server.",
		MIMEType:    "application/json",
	}

	s.mcpServer.AddResource(inspectionsResource, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		items, err := s.ListInspections()
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
