package task

import (
	"context"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	common_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

type InspectionRequest struct {
	Values map[string]any
}

var InspectionTimeTaskId = InspectionTaskPrefix + "task/time"

// InspectionTimeProducer is a provider of inspection time.
// Tasks shouldn't use time.Now() directly to make test easier.
var InspectionTimeProducer common_task.Definition = common_task.NewProcessorTask(InspectionTimeTaskId, []string{}, func(ctx context.Context, taskMode int, v *common_task.VariableSet) (any, error) {
	return time.Now(), nil
})

// TestInspectionTimeTaskProducer is a function to generate a fake InspectionTimeProducer task with the given time string.
var TestInspectionTimeTaskProducer func(timeStr string) common_task.Definition = func(timeStr string) common_task.Definition {
	return common_task.NewProcessorTask(InspectionTimeTaskId, []string{}, func(ctx context.Context, taskMode int, v *common_task.VariableSet) (any, error) {
		time, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return nil, err
		}
		return time, nil
	})
}

func GetMetadataSetFromVariable(v *common_task.VariableSet) (*metadata.MetadataSet, error) {
	return common_task.GetTypedVariableFromTaskVariable[*metadata.MetadataSet](v, MetadataVariableName, nil)
}

func GetInspectionRequestFromVariable(v *common_task.VariableSet) (*InspectionRequest, error) {
	return common_task.GetTypedVariableFromTaskVariable[*InspectionRequest](v, InspectionRequestVariableName, nil)
}

func GetInspectionTimeFromTaskVariable(v *common_task.VariableSet) (time.Time, error) {
	return common_task.GetTypedVariableFromTaskVariable[time.Time](v, InspectionTimeTaskId, time.Time{})
}
