// Copyright 2024 Google LLC
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

package privatecommon_impl

import (
	"context"
	"fmt"
	"math"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
)

// JustificationFormTask is a form task to show the current justification for any kind of inspection types.
var JustificationFormTask = inspectiontaskbase.NewInspectionTask(privatecommon_contract.JustificationFormTaskID, []taskid.UntypedTaskReference{}, func(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) (string, error) {
	gaLabelsMap := map[string]string{}
	if parameters.Private.GALabels != nil {
		gaLabelsMap = parameters.Private.GetMapOfGALabels()
	}
	if justification, found := gaLabelsMap["justification"]; found {
		metadataSet := khictx.MustGetValue(ctx, inspectioncore_contract.InspectionRunMetadata)
		formFields, found := typedmap.Get(metadataSet, inspectionmetadata.FormFieldSetMetadataKey)
		if !found {
			return "", fmt.Errorf("failed to get form fields")
		}
		formFields.SetField(inspectionmetadata.TextParameterFormField{
			ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
				ID:       privatecommon_contract.JustificationFormTaskID.ReferenceIDString(),
				Priority: math.MaxInt32,
				Type:     inspectionmetadata.Text,
				Label:    "Justification",
				HintType: inspectionmetadata.None,
			},
			Readonly: true,
			Default:  justification,
		})
		return justification, nil
	}
	return "", nil
}, coretask.NewRequiredTaskLabel())
