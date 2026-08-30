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
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	"github.com/google/go-cmp/cmp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupTestServer(t *testing.T) (*Server, *coreinspection.InspectionTaskServer) {
	t.Helper()
	taskServer, err := coreinspection.NewServer(nil)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}
	return NewServer(taskServer), taskServer
}

func registerMockInspection(ts *coreinspection.InspectionTaskServer, id, name, typeName string, phase inspectionmetadata.TaskProgressPhase, startUnix, endUnix, inspectUnix int64, fileSize int, errCount int) {
	md := typedmap.NewTypedMap()
	typedmap.Set(md, inspectionmetadata.HeaderMetadataKey, &inspectionmetadata.HeaderMetadata{
		InspectionName:         name,
		InspectionType:         typeName,
		InspectTimeUnixSeconds: inspectUnix,
		StartTimeUnixSeconds:   startUnix,
		EndTimeUnixSeconds:     endUnix,
		FileSize:               fileSize,
	})
	typedmap.Set(md, inspectionmetadata.ProgressMetadataKey, &inspectionmetadata.Progress{
		Phase: phase,
	})
	if errCount > 0 {
		errs := make([]*inspectionmetadata.ErrorMessage, errCount)
		for i := 0; i < errCount; i++ {
			errs[i] = &inspectionmetadata.ErrorMessage{ErrorId: i + 1, Message: "error"}
		}
		typedmap.Set(md, inspectionmetadata.ErrorMessageSetMetadataKey, &inspectionmetadata.ErrorMessageSetMetadata{
			ErrorMessages: errs,
		})
	}
	ts.RegisterImportedInspection(id, nil, md.AsReadonly())
}

func TestServer_ListInspections(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(ts *coreinspection.InspectionTaskServer)
		wantItems []*InspectionSummaryItem
	}{
		{
			name:      "empty inspection server returns empty list",
			setup:     func(ts *coreinspection.InspectionTaskServer) {},
			wantItems: []*InspectionSummaryItem{},
		},
		{
			name: "returns sorted active and imported inspections",
			setup: func(ts *coreinspection.InspectionTaskServer) {
				registerMockInspection(ts, "inspection-001", "test-inspection-1", "gcp-gke", inspectionmetadata.TaskPhaseDone, 1700000000, 1700003600, 1700000000, 1024, 1)
				registerMockInspection(ts, "inspection-002", "test-inspection-2", "anthos", inspectionmetadata.TaskPhaseRunning, 1700005000, 1700008600, 1700005000, 2048, 0)
			},
			wantItems: []*InspectionSummaryItem{
				{
					ID:          "inspection-001",
					Name:        "test-inspection-1",
					Type:        "gcp-gke",
					Status:      "DONE",
					StartTime:   time.Unix(1700000000, 0).UTC().Format(time.RFC3339),
					EndTime:     time.Unix(1700003600, 0).UTC().Format(time.RFC3339),
					InspectTime: time.Unix(1700000000, 0).UTC().Format(time.RFC3339),
					FileSize:    1024,
					ErrorCount:  1,
				},
				{
					ID:          "inspection-002",
					Name:        "test-inspection-2",
					Type:        "anthos",
					Status:      "RUNNING",
					StartTime:   time.Unix(1700005000, 0).UTC().Format(time.RFC3339),
					EndTime:     time.Unix(1700008600, 0).UTC().Format(time.RFC3339),
					InspectTime: time.Unix(1700005000, 0).UTC().Format(time.RFC3339),
					FileSize:    2048,
					ErrorCount:  0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, taskServer := setupTestServer(t)
			tc.setup(taskServer)

			got, err := server.ListInspections()
			if err != nil {
				t.Fatalf("ListInspections() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantItems, got); diff != "" {
				t.Errorf("ListInspections() mismatch (-want +got):\n%s", diff)
			}
		})
	}
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
			server, taskServer := setupTestServer(t)
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
			server, _ := setupTestServer(t)

			discoverBody := `{"jsonrpc":"2.0","id":1,"method":"server/discover","params":{"_meta":{"io.modelcontextprotocol/clientCapabilities":{"elicitation":{"form":{},"url":{}},"roots":{"listChanged":true}},"io.modelcontextprotocol/clientInfo":{"name":"antigravity-client","version":"v1.0.0"},"io.modelcontextprotocol/protocolVersion":"2026-07-28"}}}`
			req := httptest.NewRequest(http.MethodPost, "/mcp/sse", strings.NewReader(discoverBody))
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
