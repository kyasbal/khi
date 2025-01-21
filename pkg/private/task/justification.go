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

	form_metadata "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/form"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/header"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/parameters"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	baremetal "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gdcv-for-baremetal"
	vmware "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gdcv-for-vmware"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke"
	aws "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke-on-aws"
	azure "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke-on-azure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var availableForAllGCPInspectionTypes = inspection_task.InspectionTypeLabel(gke.InspectionTypeId, aws.InspectionTypeId, azure.InspectionTypeId, baremetal.InspectionTypeId, vmware.InspectionTypeId)

const justificationFormTaskID = gcp_task.GCPPrefix + "private/justification"

var JustificationFormTask = inspection_task.NewInspectionProcessor(justificationFormTaskID, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet, progress *progress.TaskProgress) (any, error) {
	gaLabelsMap := map[string]string{}
	if parameters.Private.GALabels != nil {
		gaLabelsMap = parameters.Private.GetMapOfGALabels()
	}
	if justification, found := gaLabelsMap["justification"]; found {
		m, err := inspection_task.GetMetadataSetFromVariable(v)
		if err != nil {
			return nil, err
		}
		formFields := m.LoadOrStore(form_metadata.FormFieldSetMetadataKey, &form_metadata.FormFieldSetMetadataFactory{}).(*form_metadata.FormFieldSet)
		formFields.SetField(&form_metadata.FormField{
			Id:        justificationFormTaskID,
			Priority:  math.MaxInt32,
			Type:      "Text",
			Label:     "Justification",
			AllowEdit: false,
			Default:   justification,
		})
		return justification, nil
	}
	return "", nil
},
	availableForAllGCPInspectionTypes,
	inspection_task.NewRequiredTaskLabel())

const filenameHeaderMetadataGeneratorTaskID = gcp_task.GCPPrefix + "private/header-metadata-filename"

var FilenameHeaderMetadataGeneratorTask = inspection_task.NewInspectionProcessor(filenameHeaderMetadataGeneratorTaskID, []string{
	justificationFormTaskID,
	gcp_task.InputClusterNameTaskID,
	gcp_task.InputEndTimeTaskID,
	gcp_task.InputStartTimeTaskID,
}, func(ctx context.Context, taskMode int, v *task.VariableSet, progress *progress.TaskProgress) (any, error) {
	m, err := inspection_task.GetMetadataSetFromVariable(v)
	if err != nil {
		return nil, err
	}
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	endTime, err := gcp_task.GetInputEndTimeFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	startTime, err := gcp_task.GetInputStartTimeFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	header := m.LoadOrStore(header.HeaderMetadataKey, &header.HeaderMetadataFactory{}).(*header.Header)
	justification, err := task.GetTypedVariableFromTaskVariable[string](v, justificationFormTaskID, "")
	if err != nil {
		return nil, err
	}
	header.SuggestedFileName = fmt.Sprintf("%s-%s-%s-%s.khi", justification, clusterName, endTime.Format("01021504"), startTime.Format("01021504"))
	return struct{}{}, nil
},
	availableForAllGCPInspectionTypes,
	inspection_task.NewRequiredTaskLabel())
