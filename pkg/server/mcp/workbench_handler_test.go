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
	"testing"
	"time"

	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/logger"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/server/workbench"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	"github.com/google/go-cmp/cmp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func createTestWorkbenchManager(t *testing.T) (*coreinspection.InspectionTaskServer, *workbench.WorkbenchManager, string) {
	t.Helper()
	logger.InitGlobalKHILogger()
	ioConfig, err := inspectioncore_contract.NewIOConfigForTest()
	if err != nil {
		t.Fatalf("failed to create test IOConfig: %v", err)
	}
	tempDir := t.TempDir()
	ioConfig.DataDestination = tempDir
	ioConfig.TemporaryFolder = tempDir
	server, err := coreinspection.NewServer(ioConfig)
	if err != nil {
		t.Fatalf("failed to create inspection server: %v", err)
	}

	inspectionType := coreinspection.InspectionType{
		Id:   "test-type",
		Name: "Test Type",
	}
	if err := server.AddInspectionType(inspectionType); err != nil {
		t.Fatalf("failed to add inspection type: %v", err)
	}

	dummyTaskID := taskid.NewDefaultImplementationID[any]("dummy-task")
	dummyTask := coretask.NewTask(
		dummyTaskID,
		nil,
		func(ctx context.Context) (any, error) {
			return "success", nil
		},
		coretask.WithLabelValue(inspectioncore_contract.LabelKeyInspectionTypes, []string{inspectionType.Id}),
		coretask.WithLabelValue(inspectioncore_contract.LabelKeyInspectionDefaultFeatureFlag, true),
		coretask.WithLabelValue(inspectioncore_contract.LabelKeyInspectionFeatureFlag, true),
		coretask.NewSubsequentTaskRefsTaskLabel(inspectioncore_contract.SerializerTaskID.Ref()),
	)
	if err := server.AddTask(dummyTask); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	inspectionID, err := server.CreateInspection(inspectionType.Id)
	if err != nil {
		t.Fatalf("failed to create inspection: %v", err)
	}

	runner := server.GetInspection(inspectionID)
	if err := runner.Run(context.Background(), &inspectioncore_contract.InspectionRequest{Values: map[string]any{}}); err != nil {
		t.Fatalf("failed to run inspection: %v", err)
	}
	<-runner.Wait()

	indexMgr := workbench.NewInspectionIndexManager(server, t.TempDir())
	wbMgr := workbench.NewWorkbenchManager(server, indexMgr, 5*time.Minute, 0)

	return server, wbMgr, inspectionID
}

func TestWorkbenchHandler_OpenWorkbench(t *testing.T) {
	inspectionServer, wbMgr, validInspectionID := createTestWorkbenchManager(t)
	defer wbMgr.Stop()

	testCases := []struct {
		name         string
		inspectionID string
		wantErr      bool
		wantOutput   *OpenWorkbenchOutput
	}{
		{
			name:         "fails when inspectionID is empty",
			inspectionID: "",
			wantErr:      true,
		},
		{
			name:         "fails when inspectionID not found",
			inspectionID: "unknown-id",
			wantErr:      true,
		},
		{
			name:         "opens valid inspection successfully",
			inspectionID: validInspectionID,
			wantErr:      false,
			wantOutput: &OpenWorkbenchOutput{
				WorkbenchID:  ActiveMCPSessionID,
				InspectionID: validInspectionID,
				Status:       "READY",
				Message:      "Workbench session opened successfully.",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewWorkbenchHandler(wbMgr)
			got, err := handler.OpenWorkbench(context.Background(), tc.inspectionID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("OpenWorkbench() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			if diff := cmp.Diff(tc.wantOutput, got); diff != "" {
				t.Errorf("OpenWorkbench() mismatch (-want +got):\n%s", diff)
			}
		})
	}
	_ = inspectionServer
}

func TestWorkbenchHandler_GetActiveWorkbenchStatus(t *testing.T) {
	_, wbMgr, validInspectionID := createTestWorkbenchManager(t)
	defer wbMgr.Stop()

	testCases := []struct {
		name       string
		preOpen    bool
		wantActive bool
	}{
		{
			name:       "returns inactive when nothing opened",
			preOpen:    false,
			wantActive: false,
		},
		{
			name:       "returns active after open",
			preOpen:    true,
			wantActive: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewWorkbenchHandler(wbMgr)
			if tc.preOpen {
				_, err := handler.OpenWorkbench(context.Background(), validInspectionID)
				if err != nil {
					t.Fatalf("OpenWorkbench() setup failed: %v", err)
				}
			}

			status, err := handler.GetActiveWorkbenchStatus()
			if err != nil {
				t.Fatalf("GetActiveWorkbenchStatus() unexpected error: %v", err)
			}

			if status.Active != tc.wantActive {
				t.Errorf("status.Active mismatch: want %v, got %v", tc.wantActive, status.Active)
			}
		})
	}
}

func TestWorkbenchHandler_MCPToolsAndResources(t *testing.T) {
	_, wbMgr, validInspectionID := createTestWorkbenchManager(t)
	defer wbMgr.Stop()

	handler := NewWorkbenchHandler(wbMgr)
	impl := &mcp.Implementation{Name: "test-workbench-mcp", Version: "1.0.0"}
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

	// 1. Check initial inactive status
	res, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: "khi://workbench/status",
	})
	if err != nil {
		t.Fatalf("ReadResource(status) failed: %v", err)
	}
	var statusOut WorkbenchStatusOutput
	if err := json.Unmarshal([]byte(res.Contents[0].Text), &statusOut); err != nil {
		t.Fatalf("unmarshal status failed: %v", err)
	}
	if statusOut.Active {
		t.Errorf("expected inactive session initially, got active")
	}

	// 2. Call open_workbench tool
	inputJSON, _ := json.Marshal(OpenWorkbenchInput{InspectionID: validInspectionID})
	var inputParams map[string]any
	_ = json.Unmarshal(inputJSON, &inputParams)

	toolRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "open_workbench",
		Arguments: inputParams,
	})
	if err != nil {
		t.Fatalf("CallTool(open_workbench) failed: %v", err)
	}
	if toolRes.IsError {
		t.Fatalf("CallTool(open_workbench) returned error")
	}

	// 3. Check status is now active
	res2, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: "khi://workbench/status",
	})
	if err != nil {
		t.Fatalf("ReadResource(status) after open failed: %v", err)
	}
	var statusOut2 WorkbenchStatusOutput
	if err := json.Unmarshal([]byte(res2.Contents[0].Text), &statusOut2); err != nil {
		t.Fatalf("unmarshal status failed: %v", err)
	}
	if !statusOut2.Active || statusOut2.WorkbenchID != ActiveMCPSessionID {
		t.Errorf("expected active session with ID %s, got: %+v", ActiveMCPSessionID, statusOut2)
	}
}
