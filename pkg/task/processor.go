package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task/taskid"
)

type ProcessorFunc = func(ctx context.Context, taskMode int, v *VariableSet) (any, error)

// NewProcessor returns a task definition generates a variable named the task Id from one or more variables generated from the dependency.
// A processor task set the variable that has the same name of the task Id at the end.
func NewProcessorTask(taskImplementationIdInString string, dependenciesInString []string, processor ProcessorFunc, labelOpts ...LabelOpt) Definition {
	taskImplementationId := taskid.NewTaskImplementationId(taskImplementationIdInString)
	taskDependencyReferenceIds := []taskid.TaskReferenceId{}
	for _, dependency := range dependenciesInString {
		taskDependencyReferenceIds = append(taskDependencyReferenceIds, taskid.NewTaskReference(dependency))
	}
	return NewDefinitionFromFunc(taskImplementationId, taskDependencyReferenceIds, func(taskMode int) Runnable {
		return NewRunnableFunc(func(ctx context.Context, v *VariableSet) error {
			result, err := processor(ctx, taskMode, v)
			if err != nil {
				return err
			}
			return v.Set(taskImplementationId.ReferenceId().String(), result)
		})
	}, labelOpts...)
}
