package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/ioconfig"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const BuilderGeneratorTaskId = InspectionTaskPrefix + "builder-generator"

var BuilderGeneratorTask = task.NewProcessorTask(BuilderGeneratorTaskId, []string{ioconfig.IOConfigTaskName}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
	ioConfig, err := ioconfig.GetIOConfigFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	return history.NewBuilder(ioConfig), nil
})
