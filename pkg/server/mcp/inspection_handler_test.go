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
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	"github.com/google/go-cmp/cmp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupTestInspectionServer(t *testing.T) *coreinspection.InspectionTaskServer {
	t.Helper()
	taskServer, err := coreinspection.NewServer(nil)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}
	return taskServer
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
			errs[i] = &inspectionmetadata.ErrorMessage{ErrorId: i + 1, Message: "error message"}
		}
		typedmap.Set(md, inspectionmetadata.ErrorMessageSetMetadataKey, &inspectionmetadata.ErrorMessageSetMetadata{
			ErrorMessages: errs,
		})
	}
	ts.RegisterImportedInspection(id, nil, md.AsReadonly())
}

func TestInspectionHandler_ListInspections(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(ts *coreinspection.InspectionTaskServer)
		wantItems []*InspectionSummaryItem
	}{
		{
			name:      "empty server returns empty list",
			setup:     func(ts *coreinspection.InspectionTaskServer) {},
			wantItems: []*InspectionSummaryItem{},
		},
		{
			name: "returns sorted inspections",
			setup: func(ts *coreinspection.InspectionTaskServer) {
				registerMockInspection(ts, "inspection-001", "test-1", "gcp-gke", inspectionmetadata.TaskPhaseDone, 1700000000, 1700003600, 1700000000, 1024, 1)
				registerMockInspection(ts, "inspection-002", "test-2", "anthos", inspectionmetadata.TaskPhaseRunning, 1700005000, 1700008600, 1700005000, 2048, 0)
			},
			wantItems: []*InspectionSummaryItem{
				{
					ID:          "inspection-001",
					Name:        "test-1",
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
					Name:        "test-2",
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
			taskServer := setupTestInspectionServer(t)
			tc.setup(taskServer)

			handler := NewInspectionHandler(taskServer)
			got, err := handler.ListInspections()
			if err != nil {
				t.Fatalf("ListInspections() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantItems, got); diff != "" {
				t.Errorf("ListInspections() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestInspectionHandler_GetInspectionSummary(t *testing.T) {
	testCases := []struct {
		name         string
		setup        func(ts *coreinspection.InspectionTaskServer)
		inspectionID string
		wantErr      bool
		wantName     string
		wantErrors   int
	}{
		{
			name:         "non-existent inspection returns error",
			setup:        func(ts *coreinspection.InspectionTaskServer) {},
			inspectionID: "non-existent",
			wantErr:      true,
		},
		{
			name: "existing inspection returns valid summary",
			setup: func(ts *coreinspection.InspectionTaskServer) {
				registerMockInspection(ts, "inspection-123", "gke-cluster-logs", "gcp-gke", inspectionmetadata.TaskPhaseDone, 1700000000, 1700003600, 1700000000, 4096, 2)
			},
			inspectionID: "inspection-123",
			wantErr:      false,
			wantName:     "gke-cluster-logs",
			wantErrors:   2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskServer := setupTestInspectionServer(t)
			tc.setup(taskServer)

			handler := NewInspectionHandler(taskServer)
			got, err := handler.GetInspectionSummary(tc.inspectionID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetInspectionSummary() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			if got.Inspection.Name != tc.wantName {
				t.Errorf("Inspection.Name mismatch: want %q, got %q", tc.wantName, got.Inspection.Name)
			}
			if len(got.ErrorMessages) != tc.wantErrors {
				t.Errorf("ErrorMessages count mismatch: want %d, got %d", tc.wantErrors, len(got.ErrorMessages))
			}
			if !strings.Contains(got.Markdown, tc.wantName) {
				t.Errorf("Markdown expected to contain %q, got:\n%s", tc.wantName, got.Markdown)
			}
		})
	}
}

func TestInspectionHandler_MCPToolsAndResources(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "calls get_inspection_summary tool and reads resource template",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskServer := setupTestInspectionServer(t)
			registerMockInspection(taskServer, "inspection-abc", "sample-cluster", "gcp-gke", inspectionmetadata.TaskPhaseDone, 1700000000, 1700003600, 1700000000, 5000, 1)

			handler := NewInspectionHandler(taskServer)
			impl := &mcp.Implementation{Name: "test-mcp", Version: "1.0.0"}
			srv := mcp.NewServer(impl, nil)
			handler.RegisterTools(srv)
			handler.RegisterResources(srv)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			t1, t2 := mcp.NewInMemoryTransports()
			client := mcp.NewClient(impl, nil)

			go func() {
				_ = srv.Run(ctx, t1)
			}()

			session, err := client.Connect(ctx, t2, nil)
			if err != nil {
				t.Fatalf("client.Connect() failed: %v", err)
			}
			defer session.Close()

			// 1. Call get_inspection_summary tool
			inputJSON, _ := json.Marshal(GetInspectionSummaryInput{InspectionID: "inspection-abc"})
			var inputParams map[string]any
			_ = json.Unmarshal(inputJSON, &inputParams)

			toolRes, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name:      "get_inspection_summary",
				Arguments: inputParams,
			})
			if err != nil {
				t.Fatalf("CallTool(get_inspection_summary) failed: %v", err)
			}
			if toolRes.IsError {
				t.Fatalf("CallTool returned error result")
			}

			textContent, ok := toolRes.Content[0].(*mcp.TextContent)
			if !ok {
				t.Fatalf("expected *mcp.TextContent, got %T", toolRes.Content[0])
			}
			var summaryOut GetInspectionSummaryOutput
			if err := json.Unmarshal([]byte(textContent.Text), &summaryOut); err != nil {
				t.Fatalf("failed to unmarshal output: %v", err)
			}
			if summaryOut.Summary.Inspection.ID != "inspection-abc" {
				t.Errorf("summary inspection ID mismatch: want %q, got %q", "inspection-abc", summaryOut.Summary.Inspection.ID)
			}

			// 2. Read resource template khi://inspections/inspection-abc/summary
			resRes, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
				URI: "khi://inspections/inspection-abc/summary",
			})
			if err != nil {
				t.Fatalf("ReadResource(summary) failed: %v", err)
			}
			if len(resRes.Contents) == 0 {
				t.Fatalf("expected resource contents, got none")
			}
			if !strings.Contains(resRes.Contents[0].Text, "sample-cluster") {
				t.Errorf("resource text expected to contain sample-cluster, got:\n%s", resRes.Contents[0].Text)
			}
		})
	}
}
