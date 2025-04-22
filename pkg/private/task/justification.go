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
	inspection_task_contextkey "github.com/GoogleCloudPlatform/khi/pkg/inspection/contextkey"
	inspection_task_interface "github.com/GoogleCloudPlatform/khi/pkg/inspection/interface"
	form_metadata "github.com/GoogleCloudPlatform/khi/pkg/inspection/metadata/form"
	"github.com/GoogleCloudPlatform/khi/pkg/inspection/metadata/header"
	inspection_task "github.com/GoogleCloudPlatform/khi/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	private_taskid "github.com/GoogleCloudPlatform/khi/pkg/private/taskid"
	gcp_task "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task"
	baremetal "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gdcv-for-baremetal"
	vmware "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gdcv-for-vmware"
	"github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gke"
	aws "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gke-on-aws"
	azure "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task/gke-on-azure"
	"github.com/GoogleCloudPlatform/khi/pkg/task"
	"github.com/GoogleCloudPlatform/khi/pkg/task/taskid"
)

var availableForAllGCPInspectionTypes = inspection_task.InspectionTypeLabel(gke.InspectionTypeId, aws.InspectionTypeId, azure.InspectionTypeId, baremetal.InspectionTypeId, vmware.InspectionTypeId)

var JustificationFormTask = inspection_task.NewInspectionTask(private_taskid.JustificationFormTaskID, []taskid.UntypedTaskReference{}, func(ctx context.Context, taskMode inspection_task_interface.InspectionTaskMode) (string, error) {
	gaLabelsMap := map[string]string{}
	if parameters.Private.GALabels != nil {
		gaLabelsMap = parameters.Private.GetMapOfGALabels()
	}
	if justification, found := gaLabelsMap["justification"]; found {
		metadataSet := khictx.MustGetValue(ctx, inspection_task_contextkey.InspectionRunMetadata)
		formFields, found := typedmap.Get(metadataSet, form_metadata.FormFieldSetMetadataKey)
		if !found {
			return "", fmt.Errorf("failed to get form fields")
		}
		formFields.SetField(form_metadata.TextParameterFormField{
			ParameterFormFieldBase: form_metadata.ParameterFormFieldBase{
				ID:       private_taskid.JustificationFormTaskID.ReferenceIDString(),
				Priority: math.MaxInt32,
				Type:     "Text",
				Label:    "Justification",
			},
			Readonly: true,
			Default:  justification,
		})
		return justification, nil
	}
	return "", nil
},
	availableForAllGCPInspectionTypes,
	inspection_task.NewRequiredTaskLabel())

var FilenameHeaderMetadataGeneratorTask = inspection_task.NewInspectionTask(private_taskid.FileNameHeaderMetadataGeneratorTask, []taskid.UntypedTaskReference{
	private_taskid.JustificationFormTaskID,
	gcp_task.InputClusterNameTaskID,
	gcp_task.InputEndTimeTaskID,
	gcp_task.InputStartTimeTaskID,
}, func(ctx context.Context, taskMode inspection_task_interface.InspectionTaskMode) (struct{}, error) {
	metadataSet := khictx.MustGetValue(ctx, inspection_task_contextkey.InspectionRunMetadata)
	header := typedmap.GetOrDefault(metadataSet, header.HeaderMetadataKey, &header.Header{})

	clusterName := task.GetTaskResult(ctx, gcp_task.InputClusterNameTaskID.GetTaskReference())
	endTime := task.GetTaskResult(ctx, gcp_task.InputEndTimeTaskID.GetTaskReference())
	startTime := task.GetTaskResult(ctx, gcp_task.InputStartTimeTaskID.GetTaskReference())
	justification := task.GetTaskResult(ctx, private_taskid.JustificationFormTaskID.GetTaskReference())

	header.SuggestedFileName = fmt.Sprintf("%s-%s-%s-%s.khi", justification, clusterName, endTime.Format("01021504"), startTime.Format("01021504"))
	return struct{}{}, nil
},
	availableForAllGCPInspectionTypes,
	inspection_task.NewRequiredTaskLabel())
