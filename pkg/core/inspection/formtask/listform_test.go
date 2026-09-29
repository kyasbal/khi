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

package formtask

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/server/upload"
	"github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore"
	"github.com/google/go-cmp/cmp"
)

func TestListFormTaskBuilder_ItemKeysAndDefault(t *testing.T) {
	testCases := []struct {
		name         string
		defaultCount int
		requestValue any
		hasRequest   bool
		wantKeys     []string
		wantValues   []string
	}{
		{
			name:         "generates default keys when no request value given",
			defaultCount: 2,
			hasRequest:   false,
			wantKeys:     []string{"0", "1"},
			wantValues:   []string{"val-0", "val-1"},
		},
		{
			name:         "uses custom string slice keys from request",
			defaultCount: 2,
			requestValue: []string{"0", "2"},
			hasRequest:   true,
			wantKeys:     []string{"0", "2"},
			wantValues:   []string{"val-0", "val-2"},
		},
		{
			name:         "uses custom any slice keys from request",
			defaultCount: 3,
			requestValue: []any{"custom-a", "custom-b"},
			hasRequest:   true,
			wantKeys:     []string{"custom-a", "custom-b"},
			wantValues:   []string{"val-custom-a", "val-custom-b"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskId := taskid.NewDefaultImplementationID[[]string]("list-test")
			builder := NewListFormTaskBuilder(
				taskId,
				10,
				"Test List",
				func(ctx context.Context, index int, itemKey string, itemFieldID string) (inspectionmetadata.ParameterFormField, string, error) {
					field := inspectionmetadata.TextParameterFormField{
						ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
							ID:    itemFieldID,
							Label: fmt.Sprintf("Item %s", itemKey),
							Type:  inspectionmetadata.Text,
						},
					}
					return field, fmt.Sprintf("val-%s", itemKey), nil
				},
			).WithDefaultCount(tc.defaultCount).
				WithDescription("Test list description").
				WithAddButtonLabel("Add Custom Item").
				WithDependencies([]coretask.Dependency{})

			task := builder.Build()
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(context.Background())
			inputs := map[string]any{}
			if tc.hasRequest {
				inputs[taskId.ReferenceIDString()] = tc.requestValue
			}

			result, metadata, err := inspectiontest.RunInspectionTask(ctx, task, inspectioncore.TaskModeDryRun, inputs)
			if err != nil {
				t.Fatalf("dry run unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantValues, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}

			fields, found := typedmap.Get(metadata, inspectionmetadata.FormFieldSetMetadataKey)
			if !found {
				t.Fatal("form field set metadata not found")
			}
			rawField := fields.DangerouslyGetField(taskId.ReferenceIDString())
			listField, ok := rawField.(inspectionmetadata.ListParameterFormField)
			if !ok {
				t.Fatalf("field type is %T, want ListParameterFormField", rawField)
			}

			if listField.AddButtonLabel != "Add Custom Item" {
				t.Errorf("listField.AddButtonLabel = %q, want %q", listField.AddButtonLabel, "Add Custom Item")
			}

			gotKeys := make([]string, len(listField.Items))
			for i, item := range listField.Items {
				gotKeys[i] = item.Key
			}
			if diff := cmp.Diff(tc.wantKeys, gotKeys); diff != "" {
				t.Errorf("list item keys mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestListFormTaskBuilder_CountConstraintsAndValidation(t *testing.T) {
	testCases := []struct {
		name             string
		minCount         int
		maxCount         int
		requestKeys      []string
		customValidator  ListFormValidator[string]
		customHintFunc   ListFormHintGenerator[string]
		wantDryRunHint   string
		wantDryRunType   inspectionmetadata.ParameterHintType
		wantRunErrSubstr string
	}{
		{
			name:             "violates minCount constraint",
			minCount:         2,
			maxCount:         5,
			requestKeys:      []string{"0"},
			wantDryRunHint:   "at least 2 item(s) are required",
			wantDryRunType:   inspectionmetadata.Error,
			wantRunErrSubstr: "at least 2 item(s) are required",
		},
		{
			name:             "violates maxCount constraint",
			minCount:         1,
			maxCount:         2,
			requestKeys:      []string{"0", "1", "2"},
			wantDryRunHint:   "at most 2 item(s) are allowed",
			wantDryRunType:   inspectionmetadata.Error,
			wantRunErrSubstr: "at most 2 item(s) are allowed",
		},
		{
			name:        "violates custom validator",
			minCount:    1,
			maxCount:    5,
			requestKeys: []string{"invalid"},
			customValidator: func(ctx context.Context, values []string) (string, error) {
				for _, v := range values {
					if v == "val-invalid" {
						return "invalid item detected", nil
					}
				}
				return "", nil
			},
			wantDryRunHint:   "invalid item detected",
			wantDryRunType:   inspectionmetadata.Error,
			wantRunErrSubstr: "invalid item detected",
		},
		{
			name:        "passes validation and applies custom hint",
			minCount:    1,
			maxCount:    5,
			requestKeys: []string{"item-1"},
			customValidator: func(ctx context.Context, values []string) (string, error) {
				return "", nil
			},
			customHintFunc: func(ctx context.Context, values []string) (string, inspectionmetadata.ParameterHintType, error) {
				return fmt.Sprintf("configured %d items", len(values)), inspectionmetadata.Warning, nil
			},
			wantDryRunHint:   "configured 1 items",
			wantDryRunType:   inspectionmetadata.Warning,
			wantRunErrSubstr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskId := taskid.NewDefaultImplementationID[[]string]("list-validation-test")
			builder := NewListFormTaskBuilder(
				taskId,
				5,
				"Validation List",
				func(ctx context.Context, index int, itemKey string, itemFieldID string) (inspectionmetadata.ParameterFormField, string, error) {
					field := inspectionmetadata.TextParameterFormField{
						ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
							ID:    itemFieldID,
							Label: itemKey,
							Type:  inspectionmetadata.Text,
						},
					}
					return field, fmt.Sprintf("val-%s", itemKey), nil
				},
			).WithMinCount(tc.minCount).WithMaxCount(tc.maxCount)

			if tc.customValidator != nil {
				builder.WithValidator(tc.customValidator)
			}
			if tc.customHintFunc != nil {
				builder.WithHintFunc(tc.customHintFunc)
			}

			task := builder.Build()
			inputs := map[string]any{
				taskId.ReferenceIDString(): tc.requestKeys,
			}

			// DryRun mode
			dryRunCtx := inspectiontest.WithDefaultTestInspectionTaskContext(context.Background())
			_, metadata, dryRunErr := inspectiontest.RunInspectionTask(dryRunCtx, task, inspectioncore.TaskModeDryRun, inputs)
			if dryRunErr != nil {
				t.Fatalf("dry run unexpected error: %v", dryRunErr)
			}

			fields, found := typedmap.Get(metadata, inspectionmetadata.FormFieldSetMetadataKey)
			if !found {
				t.Fatal("form field set metadata not found")
			}
			rawField := fields.DangerouslyGetField(taskId.ReferenceIDString())
			listField, ok := rawField.(inspectionmetadata.ListParameterFormField)
			if !ok {
				t.Fatalf("field type is %T, want ListParameterFormField", rawField)
			}

			if listField.Hint != tc.wantDryRunHint {
				t.Errorf("listField.Hint = %q, want %q", listField.Hint, tc.wantDryRunHint)
			}
			if listField.HintType != tc.wantDryRunType {
				t.Errorf("listField.HintType = %v, want %v", listField.HintType, tc.wantDryRunType)
			}

			// Run mode
			runCtx := inspectiontest.WithDefaultTestInspectionTaskContext(context.Background())
			_, _, runErr := inspectiontest.RunInspectionTask(runCtx, task, inspectioncore.TaskModeRun, inputs)
			if tc.wantRunErrSubstr != "" {
				if runErr == nil {
					t.Fatalf("expected run mode error containing %q, got nil", tc.wantRunErrSubstr)
				}
				if !strings.Contains(runErr.Error(), tc.wantRunErrSubstr) {
					t.Errorf("run error %q does not contain expected substring %q", runErr.Error(), tc.wantRunErrSubstr)
				}
			} else if runErr != nil {
				t.Fatalf("run mode unexpected error: %v", runErr)
			}
		})
	}
}

func TestListFormTaskBuilder_InvalidRequestInput(t *testing.T) {
	testCases := []struct {
		name         string
		requestValue any
		wantErrSub   string
	}{
		{
			name:         "duplicate keys in slice",
			requestValue: []string{"0", "0"},
			wantErrSub:   "duplicate item key",
		},
		{
			name:         "empty key string",
			requestValue: []string{"0", ""},
			wantErrSub:   "empty item key",
		},
		{
			name:         "non-slice input type",
			requestValue: 12345,
			wantErrSub:   "expected []string or []any",
		},
		{
			name:         "non-string element in any slice",
			requestValue: []any{"0", 123},
			wantErrSub:   "is not a string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskId := taskid.NewDefaultImplementationID[[]string]("invalid-input-test")
			builder := NewListFormTaskBuilder(
				taskId,
				1,
				"Invalid Input List",
				func(ctx context.Context, index int, itemKey string, itemFieldID string) (inspectionmetadata.ParameterFormField, string, error) {
					return inspectionmetadata.TextParameterFormField{
						ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
							ID:   itemFieldID,
							Type: inspectionmetadata.Text,
						},
					}, itemKey, nil
				},
			)

			task := builder.Build()
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(context.Background())
			inputs := map[string]any{
				taskId.ReferenceIDString(): tc.requestValue,
			}

			_, _, err := inspectiontest.RunInspectionTask(ctx, task, inspectioncore.TaskModeDryRun, inputs)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErrSub)
			}
			if !strings.Contains(err.Error(), tc.wantErrSub) {
				t.Errorf("error %q does not contain expected substring %q", err.Error(), tc.wantErrSub)
			}
		})
	}
}

func TestListFormTaskBuilder_FileItems(t *testing.T) {
	mockToken := mockUploadToken{id: "list-file-token"}

	testCases := []struct {
		name             string
		mode             inspectioncore.InspectionTaskModeType
		storeStatus      upload.UploadStatus
		wantErr          bool
		wantErrSubstring string
		wantItemPending  bool
	}{
		{
			name:            "TaskModeDryRun with uploading status succeeds and sets pending",
			mode:            inspectioncore.TaskModeDryRun,
			storeStatus:     upload.UploadStatusUploading,
			wantErr:         false,
			wantItemPending: true,
		},
		{
			name:        "TaskModeRun with completed status returns upload results",
			mode:        inspectioncore.TaskModeRun,
			storeStatus: upload.UploadStatusCompleted,
			wantErr:     false,
		},
		{
			name:             "TaskModeRun with waiting status returns error",
			mode:             inspectioncore.TaskModeRun,
			storeStatus:      upload.UploadStatusWaiting,
			wantErr:          true,
			wantErrSubstring: "file upload is not completed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prevStore := upload.DefaultUploadFileStore
			defer func() { upload.DefaultUploadFileStore = prevStore }()

			storeResult := upload.UploadResult{
				Status: tc.storeStatus,
				Token:  mockToken,
			}
			upload.DefaultUploadFileStore = &mockFileFormStore{
				token:           mockToken,
				result:          storeResult,
				completedResult: storeResult,
			}

			taskId := taskid.NewDefaultImplementationID[[]upload.UploadResult]("files-list-task")
			builder := NewListFormTaskBuilder(
				taskId,
				10,
				"Log Files",
				func(ctx context.Context, index int, itemKey string, itemFieldID string) (inspectionmetadata.ParameterFormField, upload.UploadResult, error) {
					return EvaluateFileFormField(ctx, itemFieldID, index, fmt.Sprintf("File %s", itemKey), "Upload log", nil)
				},
			).WithDefaultCount(1)

			task := builder.Build()
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(context.Background())
			result, metadata, err := inspectiontest.RunInspectionTask(ctx, task, tc.mode, map[string]any{})

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tc.wantErrSubstring != "" && !strings.Contains(err.Error(), tc.wantErrSubstring) {
					t.Errorf("error %q does not contain expected substring %q", err.Error(), tc.wantErrSubstring)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != 1 {
				t.Fatalf("result length = %d, want 1", len(result))
			}

			fields, found := typedmap.Get(metadata, inspectionmetadata.FormFieldSetMetadataKey)
			if !found {
				t.Fatal("form field set metadata not found")
			}
			rawField := fields.DangerouslyGetField(taskId.ReferenceIDString())
			listField, ok := rawField.(inspectionmetadata.ListParameterFormField)
			if !ok {
				t.Fatalf("field type is %T, want ListParameterFormField", rawField)
			}
			if len(listField.Items) != 1 {
				t.Fatalf("list items count = %d, want 1", len(listField.Items))
			}

			fileField, ok := listField.Items[0].Field.(inspectionmetadata.FileParameterFormField)
			if !ok {
				t.Fatalf("item field type is %T, want FileParameterFormField", listField.Items[0].Field)
			}
			if fileField.Pending != tc.wantItemPending {
				t.Errorf("fileField.Pending = %v, want %v", fileField.Pending, tc.wantItemPending)
			}
		})
	}
}

type testNodeLogs struct {
	NodeName   string
	Kubelet    upload.UploadResult
	Containerd upload.UploadResult
}

func TestListFormTaskBuilder_GroupItems(t *testing.T) {
	testCases := []struct {
		name       string
		request    []string
		wantLength int
	}{
		{
			name:       "generates repeatable group items with multiple file children",
			request:    []string{"node-a", "node-b"},
			wantLength: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prevStore := upload.DefaultUploadFileStore
			defer func() { upload.DefaultUploadFileStore = prevStore }()

			mockToken := mockUploadToken{id: "group-file-token"}
			storeResult := upload.UploadResult{
				Status: upload.UploadStatusCompleted,
				Token:  mockToken,
			}
			upload.DefaultUploadFileStore = &mockFileFormStore{
				token:           mockToken,
				result:          storeResult,
				completedResult: storeResult,
			}

			taskId := taskid.NewDefaultImplementationID[[]testNodeLogs]("node-group-list-task")
			builder := NewListFormTaskBuilder(
				taskId,
				1,
				"Nodes",
				func(ctx context.Context, index int, itemKey string, itemFieldID string) (inspectionmetadata.ParameterFormField, testNodeLogs, error) {
					kubeletID := fmt.Sprintf("%s/kubelet", itemFieldID)
					kubeletField, kubeletRes, kubeletErr := EvaluateFileFormField(ctx, kubeletID, 1, "Kubelet Log", "", nil)
					if kubeletErr != nil {
						return nil, testNodeLogs{}, kubeletErr
					}

					containerdID := fmt.Sprintf("%s/containerd", itemFieldID)
					containerdField, containerdRes, containerdErr := EvaluateFileFormField(ctx, containerdID, 2, "Containerd Log", "", nil)
					if containerdErr != nil {
						return nil, testNodeLogs{}, containerdErr
					}

					groupField := inspectionmetadata.GroupParameterFormField{
						ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
							ID:          itemFieldID,
							Label:       fmt.Sprintf("Node: %s", itemKey),
							Type:        inspectionmetadata.Group,
							Description: "Per-node logs",
						},
						Children: []inspectionmetadata.ParameterFormField{
							kubeletField,
							containerdField,
						},
					}

					return groupField, testNodeLogs{
						NodeName:   itemKey,
						Kubelet:    kubeletRes,
						Containerd: containerdRes,
					}, nil
				},
			)

			task := builder.Build()
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(context.Background())
			inputs := map[string]any{
				taskId.ReferenceIDString(): tc.request,
			}

			result, metadata, err := inspectiontest.RunInspectionTask(ctx, task, inspectioncore.TaskModeRun, inputs)
			if err != nil {
				t.Fatalf("run mode unexpected error: %v", err)
			}

			if len(result) != tc.wantLength {
				t.Fatalf("result length = %d, want %d", len(result), tc.wantLength)
			}
			for i, key := range tc.request {
				if result[i].NodeName != key {
					t.Errorf("result[%d].NodeName = %q, want %q", i, result[i].NodeName, key)
				}
			}

			fields, found := typedmap.Get(metadata, inspectionmetadata.FormFieldSetMetadataKey)
			if !found {
				t.Fatal("form field set metadata not found")
			}
			rawField := fields.DangerouslyGetField(taskId.ReferenceIDString())
			listField, ok := rawField.(inspectionmetadata.ListParameterFormField)
			if !ok {
				t.Fatalf("field type is %T, want ListParameterFormField", rawField)
			}
			if len(listField.Items) != tc.wantLength {
				t.Fatalf("list items count = %d, want %d", len(listField.Items), tc.wantLength)
			}

			// Verify recursive file collection on the metadata
			fileIDs := fields.GetFileFieldIDs()
			wantFileIDs := []string{
				"node-group-list-task/node-a/kubelet",
				"node-group-list-task/node-a/containerd",
				"node-group-list-task/node-b/kubelet",
				"node-group-list-task/node-b/containerd",
			}
			if diff := cmp.Diff(wantFileIDs, fileIDs); diff != "" {
				t.Errorf("GetFileFieldIDs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
