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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupTestServer(t *testing.T) (*Server, *coreinspection.InspectionTaskServer, *InspectionHandler) {
	t.Helper()
	taskServer, err := coreinspection.NewServer(nil)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}
	inspectionHandler := NewInspectionHandler(taskServer)
	return NewServer(inspectionHandler), taskServer, inspectionHandler
}

func TestServer_ClientSession_ListInspections(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(ts *coreinspection.InspectionTaskServer)
		wantCount int
	}{
		{
			name:      "executes list_inspections tool and returns empty list",
			setup:     func(ts *coreinspection.InspectionTaskServer) {},
			wantCount: 0,
		},
		{
			name: "executes list_inspections tool and returns populated items",
			setup: func(ts *coreinspection.InspectionTaskServer) {
				registerMockInspection(ts, "sample-id", "sample", "gcp-gke", inspectionmetadata.TaskPhaseDone, 0, 0, 0, 0, 0)
			},
			wantCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, taskServer, _ := setupTestServer(t)
			tc.setup(taskServer)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			t1, t2 := mcp.NewInMemoryTransports()
			client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)

			go func() {
				_ = server.MCPServer().Run(ctx, t1)
			}()

			session, err := client.Connect(ctx, t2, nil)
			if err != nil {
				t.Fatalf("client.Connect() failed: %v", err)
			}
			defer session.Close()

			res, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name: "list_inspections",
			})
			if err != nil {
				t.Fatalf("CallTool() failed: %v", err)
			}
			if res.IsError {
				t.Fatalf("CallTool() returned error result")
			}

			if len(res.Content) == 0 {
				t.Fatalf("expected result content, got none")
			}

			text, ok := res.Content[0].(*mcp.TextContent)
			if !ok {
				t.Fatalf("expected *mcp.TextContent, got %T", res.Content[0])
			}

			var out ListInspectionsOutput
			if err := json.Unmarshal([]byte(text.Text), &out); err != nil {
				t.Fatalf("failed to unmarshal output JSON: %v", err)
			}

			if len(out.Inspections) != tc.wantCount {
				t.Errorf("inspections count mismatch: want %d, got %d", tc.wantCount, len(out.Inspections))
			}
		})
	}
}

func TestServer_HTTPHandler_Discover(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "handles server/discover POST request successfully",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, _, _ := setupTestServer(t)

			discoverBody := `{"jsonrpc":"2.0","id":1,"method":"server/discover","params":{"_meta":{"io.modelcontextprotocol/clientCapabilities":{"elicitation":{"form":{},"url":{}},"roots":{"listChanged":true}},"io.modelcontextprotocol/clientInfo":{"name":"antigravity-client","version":"v1.0.0"},"io.modelcontextprotocol/protocolVersion":"2026-07-28"}}}`
			req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/mcp/sse", strings.NewReader(discoverBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json, text/event-stream")
			req.Header.Set("Mcp-Protocol-Version", "2026-07-28")
			req.Header.Set("Mcp-Method", "server/discover")
			w := httptest.NewRecorder()

			server.HTTPHandler().ServeHTTP(w, req)

			if w.Code >= http.StatusBadRequest {
				t.Fatalf("HTTPHandler() returned error code %d: %s", w.Code, w.Body.String())
			}
		})
	}
}
