// Copyright 2025 Google LLC
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

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/server/upload"
	core_contract "github.com/GoogleCloudPlatform/khi/pkg/task/core/contract"
	"github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore"
)

type FileFormTaskBuilder struct {
	FormTaskBuilderBase[upload.UploadResult]
	verifier upload.UploadFileVerifier
}

func NewFileFormTaskBuilder(id taskid.TaskImplementationID[upload.UploadResult], priority int, label string, verifier upload.UploadFileVerifier) *FileFormTaskBuilder {
	return &FileFormTaskBuilder{
		FormTaskBuilderBase: NewFormTaskBuilderBase(id, priority, label),
		verifier:            verifier,
	}
}

// WithDependencies sets the task dependencies
func (b *FileFormTaskBuilder) WithDependencies(dependencies []coretask.Dependency) *FileFormTaskBuilder {
	b.FormTaskBuilderBase.WithDependencies(dependencies)
	return b
}

// WithDescription sets the description for the form field
func (b *FileFormTaskBuilder) WithDescription(description string) *FileFormTaskBuilder {
	b.FormTaskBuilderBase.WithDescription(description)
	return b
}

// EvaluateFileFormField evaluates a file upload parameter for the given fieldID,
// issuing an upload token and checking upload and verification status according to the current task mode.
func EvaluateFileFormField(
	ctx context.Context,
	fieldID string,
	priority int,
	label string,
	description string,
	verifier upload.UploadFileVerifier,
) (inspectionmetadata.FileParameterFormField, upload.UploadResult, error) {
	req := khictx.MustGetValue(ctx, inspectioncore.InspectionTaskInput)
	taskMode := khictx.MustGetValue(ctx, inspectioncore.InspectionTaskMode)

	token := upload.DefaultUploadFileStore.GetUploadToken(GenerateUploadIDWithTaskContext(ctx, fieldID), verifier, fieldID)

	var uploadResult upload.UploadResult
	var err error
	if taskMode == inspectioncore.TaskModeRun {
		uploadResult, err = upload.DefaultUploadFileStore.GetCompletedResult(ctx, token, req)
	} else {
		uploadResult, err = upload.DefaultUploadFileStore.GetResult(token, req)
	}
	if err != nil {
		return inspectionmetadata.FileParameterFormField{}, upload.UploadResult{}, err
	}

	field := inspectionmetadata.FileParameterFormField{
		ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
			ID:          fieldID,
			Priority:    priority,
			Label:       label,
			Description: description,
			Type:        inspectionmetadata.File,
			HintType:    inspectionmetadata.None,
			Hint:        "",
		},
		Token:  token,
		Status: uploadResult.Status,
	}
	field = setFormHintsFromUploadResult(uploadResult, field)

	var runErr error
	if taskMode == inspectioncore.TaskModeRun {
		taskID := khictx.MustGetValue(ctx, core_contract.TaskImplementationIDContextKey)
		switch {
		case uploadResult.UploadError != nil:
			runErr = fmt.Errorf("file upload failed in task %s: %w", taskID, uploadResult.UploadError)
		case uploadResult.VerificationError != nil:
			runErr = fmt.Errorf("file verification failed in task %s: %w", taskID, uploadResult.VerificationError)
		case uploadResult.Status != upload.UploadStatusCompleted:
			runErr = fmt.Errorf("file upload is not completed in task %s (current status: %d)", taskID, uploadResult.Status)
		}
	}

	return field, uploadResult, runErr
}

// Build constructs the coretask.Task that manages file form evaluation and upload verification.
func (b *FileFormTaskBuilder) Build(labelOpts ...coretask.LabelOpt) coretask.Task[upload.UploadResult] {
	return coretask.NewTask(b.FormTaskBuilderBase.id, b.FormTaskBuilderBase.dependencies, func(ctx context.Context) (upload.UploadResult, error) {
		metadata := khictx.MustGetValue(ctx, inspectionmetadata.MapContextKey)

		fieldID := b.FormTaskBuilderBase.id.ReferenceIDString()
		field, uploadResult, evaluateErr := EvaluateFileFormField(
			ctx,
			fieldID,
			b.FormTaskBuilderBase.priority,
			b.FormTaskBuilderBase.label,
			b.FormTaskBuilderBase.description,
			b.verifier,
		)
		if field.ID != "" {
			formFields, found := typedmap.Get(metadata, inspectionmetadata.FormFieldSetMetadataKey)
			if !found {
				return upload.UploadResult{}, fmt.Errorf("failed to get form fields from metadata")
			}
			err := formFields.SetField(field)
			if err != nil {
				return upload.UploadResult{}, fmt.Errorf("failed to configure the form metadata in task `%s`\n%v", b.FormTaskBuilderBase.id, err)
			}
		}
		if evaluateErr != nil {
			return upload.UploadResult{}, evaluateErr
		}

		return uploadResult, nil
	}, append(labelOpts, inspectioncore.NewFormTaskLabelOpt(b.label, b.description))...)
}

// setFormHintsFromUploadResult sets the appropriate hint and hint type on a form field
// based on the upload result status and any errors encountered during the upload process.
func setFormHintsFromUploadResult(result upload.UploadResult, field inspectionmetadata.FileParameterFormField) inspectionmetadata.FileParameterFormField {
	switch {
	case result.UploadError != nil:
		field.Hint = result.UploadError.Error()
		field.HintType = inspectionmetadata.Error
	case result.VerificationError != nil:
		field.Hint = result.VerificationError.Error()
		field.HintType = inspectionmetadata.Error
	case result.Status == upload.UploadStatusWaiting:
		field.Hint = "Waiting a file to be uploaded."
		field.HintType = inspectionmetadata.Error
	case result.Status != upload.UploadStatusCompleted:
		field.Hint = "File is being processed. Please wait a moment."
		field.HintType = inspectionmetadata.Info
		field.Pending = true
	}
	return field
}

// GenerateUploadIDWithTaskContext generates the upload ID from form ID and task ID.
func GenerateUploadIDWithTaskContext(ctx context.Context, formId string) string {
	inspectionID := khictx.MustGetValue(ctx, inspectioncore.InspectionTaskInspectionID)
	taskID := khictx.MustGetValue(ctx, core_contract.TaskImplementationIDContextKey)
	return strings.ReplaceAll(fmt.Sprintf("%s_%s_%s", inspectionID, taskID.ReferenceIDString(), formId), "/", "_")
}
