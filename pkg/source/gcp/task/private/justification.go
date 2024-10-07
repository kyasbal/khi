package private

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	form_metadata "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/form"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/header"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	baremetal "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gdcv-for-baremetal"
	vmware "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gdcv-for-vmware"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke"
	aws "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke-on-aws"
	azure "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke-on-azure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var availableForAllGCPInspectionTypes = inspection_task.InspectionTaskLabel(gke.InspectionTypeId, aws.InspectionTypeId, azure.InspectionTypeId, baremetal.InspectionTypeId, vmware.InspectionTypeId)

// This will be deprecated in the future. KHI_GA_LABELS are used for frontend initially, we should define new environment variables for analytics.
// But for now, we will use the old environment variable in the transition time.
func getCommaSeperatedKVPairEnv() map[string]any {
	result := make(map[string]any)
	if env, hasEnv := os.LookupEnv("KHI_GA_LABELS"); hasEnv {
		keyValuePairs := strings.Split(env, ",")
		for _, pair := range keyValuePairs {
			keyValues := strings.Split(pair, "=")
			key := keyValues[0]
			value := "null"
			if len(keyValues) > 1 {
				value = keyValues[1]
			}
			result[key] = value
		}
	}

	return result
}

const justificationFormTaskId = gcp_task.GCPPrefix + "private/justification"

var JustificationFormTask = inspection_task.NewInspectionProcessor(justificationFormTaskId, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet, progress *progress.TaskProgress) (any, error) {
	gaLabelPairs := getCommaSeperatedKVPairEnv()
	if justification, found := gaLabelPairs["justification"]; found {
		m, err := inspection_task.GetMetadataSetFromVariable(v)
		if err != nil {
			return nil, err
		}
		formFields := m.LoadOrStore(form_metadata.FormFieldSetMetadataKey, &form_metadata.FormFieldSetMetadataFactory{}).(*form_metadata.FormFieldSet)
		formFields.SetField(&form_metadata.FormField{
			Id:        justificationFormTaskId,
			Priority:  math.MaxInt32,
			Type:      "Text",
			Label:     "Justification",
			AllowEdit: false,
			Default:   justification.(string),
		})
		return justification, nil
	}
	return "", nil
},
	availableForAllGCPInspectionTypes,
	inspection_task.NewRequiredTaskLabel())

const filenameHeaderMetadataGeneratorTaskId = gcp_task.GCPPrefix + "private/header-metadata-filename"

var FilenameHeaderMetadataGeneratorTask = inspection_task.NewInspectionProcessor(filenameHeaderMetadataGeneratorTaskId, []string{
	justificationFormTaskId,
	gcp_task.InputClusterName,
	gcp_task.InputEndTimeVariableName,
	gcp_task.InputStartTimeVariableName,
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
	justification, err := task.GetTypedVariableFromTaskVariable[string](v, justificationFormTaskId, "")
	if err != nil {
		return nil, err
	}
	header.SuggestedFileName = fmt.Sprintf("%s-%s-%s-%s.khi", justification, clusterName, endTime.Format("01021504"), startTime.Format("01021504"))
	return struct{}{}, nil
},
	availableForAllGCPInspectionTypes,
	inspection_task.NewRequiredTaskLabel())
