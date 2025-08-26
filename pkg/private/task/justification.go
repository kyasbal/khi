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

package task

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
	private_taskid "github.com/GoogleCloudPlatform/khi/pkg/private/taskid"
	gcp_task "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task"
	baremetal "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gdcv-for-baremetal"
	vmware "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gdcv-for-vmware"
	"github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gke"
	aws "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gke-on-aws"
	azure "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gke-on-azure"
	inspection_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/contract"
)

var availableForAllGCPInspectionTypes = inspection_contract.InspectionTypeLabel(gke.InspectionTypeId, aws.InspectionTypeId, azure.InspectionTypeId, baremetal.InspectionTypeId, vmware.InspectionTypeId)

var JustificationFormTask = inspectiontaskbase.NewInspectionTask(private_taskid.JustificationFormTaskID, []taskid.UntypedTaskReference{}, func(ctx context.Context, taskMode inspection_contract.InspectionTaskModeType) (string, error) {
	gaLabelsMap := map[string]string{}
	if parameters.Private.GALabels != nil {
		gaLabelsMap = parameters.Private.GetMapOfGALabels()
	}
	if justification, found := gaLabelsMap["justification"]; found {
		metadataSet := khictx.MustGetValue(ctx, inspection_contract.InspectionRunMetadata)
		formFields, found := typedmap.Get(metadataSet, inspectionmetadata.FormFieldSetMetadataKey)
		if !found {
			return "", fmt.Errorf("failed to get form fields")
		}
		formFields.SetField(inspectionmetadata.TextParameterFormField{
			ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
				ID:       private_taskid.JustificationFormTaskID.ReferenceIDString(),
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
},
	availableForAllGCPInspectionTypes,
	inspection_contract.NewRequiredTaskLabel())

var FilenameHeaderMetadataGeneratorTask = inspectiontaskbase.NewInspectionTask(private_taskid.FileNameHeaderMetadataGeneratorTask, []taskid.UntypedTaskReference{
	private_taskid.JustificationFormTaskID.Ref(),
	gcp_task.InputClusterNameTaskID.Ref(),
	gcp_task.InputEndTimeTaskID.Ref(),
	gcp_task.InputStartTimeTaskID.Ref(),
}, func(ctx context.Context, taskMode inspection_contract.InspectionTaskModeType) (struct{}, error) {
	metadataSet := khictx.MustGetValue(ctx, inspection_contract.InspectionRunMetadata)
	header := typedmap.GetOrDefault(metadataSet, inspectionmetadata.HeaderMetadataKey, &inspectionmetadata.HeaderMetadata{})

	clusterName := coretask.GetTaskResult(ctx, gcp_task.InputClusterNameTaskID.Ref())
	endTime := coretask.GetTaskResult(ctx, gcp_task.InputEndTimeTaskID.Ref())
	startTime := coretask.GetTaskResult(ctx, gcp_task.InputStartTimeTaskID.Ref())
	justification := coretask.GetTaskResult(ctx, private_taskid.JustificationFormTaskID.Ref())

	header.SuggestedFileName = fmt.Sprintf("%s-%s-%s-%s.khi", justification, clusterName, endTime.Format("01021504"), startTime.Format("01021504"))
	return struct{}{}, nil
},
	availableForAllGCPInspectionTypes,
	inspection_contract.NewRequiredTaskLabel())
