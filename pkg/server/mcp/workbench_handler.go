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
	"sync"

	"github.com/GoogleCloudPlatform/khi/pkg/server/workbench"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ActiveMCPSessionID is the fixed workbench session identifier used by the MCP server.
const ActiveMCPSessionID = "mcp-active-session"

// WorkbenchHandler handles MCP tools and resources related to Workbench sessions.
type WorkbenchHandler struct {
	workbenchManager *workbench.WorkbenchManager
	mu               sync.RWMutex
}

var _ DomainHandler = (*WorkbenchHandler)(nil)

// NewWorkbenchHandler creates a new WorkbenchHandler instance.
func NewWorkbenchHandler(workbenchManager *workbench.WorkbenchManager) *WorkbenchHandler {
	return &WorkbenchHandler{
		workbenchManager: workbenchManager,
	}
}

// Register registers all tools and resources provided by the workbench handler onto the MCP server.
func (h *WorkbenchHandler) Register(srv *mcp.Server) {
	h.RegisterTools(srv)
	h.RegisterResources(srv)
}

// OpenWorkbench opens or replaces the active MCP workbench session for the specified inspection ID.
func (h *WorkbenchHandler) OpenWorkbench(ctx context.Context, inspectionID string) (*OpenWorkbenchOutput, error) {
	if h.workbenchManager == nil {
		return nil, fmt.Errorf("workbench manager is not available")
	}
	if inspectionID == "" {
		return nil, fmt.Errorf("inspectionId is required")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	wb, err := h.workbenchManager.GetOrOpen(ctx, ActiveMCPSessionID, inspectionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open workbench: %w", err)
	}

	return &OpenWorkbenchOutput{
		WorkbenchID:  wb.ID(),
		InspectionID: wb.InspectionID(),
		Status:       "READY",
		Message:      "Workbench session opened successfully.",
	}, nil
}

// GetActiveWorkbenchStatus retrieves the current state and indexing progress of the active MCP workbench session.
func (h *WorkbenchHandler) GetActiveWorkbenchStatus() (*WorkbenchStatusOutput, error) {
	if h.workbenchManager == nil {
		return &WorkbenchStatusOutput{Active: false}, nil
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	wb, err := h.workbenchManager.Get(ActiveMCPSessionID)
	if err != nil || wb == nil || wb.IsClosed() {
		return &WorkbenchStatusOutput{Active: false}, nil
	}

	indexState, _, _, _ := wb.IndexStatus()
	var stateStr string
	switch indexState {
	case workbench.IndexStateBuilding:
		stateStr = "BUILDING"
	case workbench.IndexStateReady:
		stateStr = "READY"
	case workbench.IndexStateFailed:
		stateStr = "FAILED"
	default:
		stateStr = "NOT_STARTED"
	}

	return &WorkbenchStatusOutput{
		Active:       true,
		WorkbenchID:  wb.ID(),
		InspectionID: wb.InspectionID(),
		IndexState:   stateStr,
	}, nil
}

// RegisterTools registers workbench-related tools onto the MCP server.
func (h *WorkbenchHandler) RegisterTools(srv *mcp.Server) {
	openWorkbenchTool := &mcp.Tool{
		Name:        "open_workbench",
		Description: "Opens an inspection dataset in the active KHI Workbench session, initializing in-memory caches and search indexes.",
	}
	mcp.AddTool(srv, openWorkbenchTool, func(ctx context.Context, req *mcp.CallToolRequest, input OpenWorkbenchInput) (*mcp.CallToolResult, *OpenWorkbenchOutput, error) {
		if input.InspectionID == "" {
			return nil, nil, fmt.Errorf("inspectionId is required")
		}
		res, err := h.OpenWorkbench(ctx, input.InspectionID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open workbench: %w", err)
		}
		return nil, res, nil
	})
}

// RegisterResources registers workbench-related resources onto the MCP server.
func (h *WorkbenchHandler) RegisterResources(srv *mcp.Server) {
	workbenchStatusResource := &mcp.Resource{
		URI:         "khi://workbench/status",
		Name:        "Workbench Session Status",
		Description: "Current status and active inspection ID of the KHI Workbench MCP session.",
		MIMEType:    "application/json",
	}
	srv.AddResource(workbenchStatusResource, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		status, err := h.GetActiveWorkbenchStatus()
		if err != nil {
			return nil, fmt.Errorf("failed to get workbench status: %w", err)
		}
		data, err := json.Marshal(status)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal workbench status: %w", err)
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "khi://workbench/status",
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})
}
