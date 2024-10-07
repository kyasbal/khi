package task_test

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

// Deprecated. Use testtask package instead.
func MockProcessorTaskFromTaskId(taskId string, value any) task.Definition {
	return task.NewProcessorTask(taskId, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		return value, nil
	})
}

// Deprecated. Use testtask package instead.
// RunTaskGraph executes the task graph just with provided dependency tasks
func RunTaskGraph(target task.Definition, mode int, initialParameters map[string]any, dependencies ...task.Definition) (*task.VariableSet, error) {
	sourceDs, err := task.NewSet([]task.Definition{target})
	if err != nil {
		return nil, err
	}
	availableDs, err := task.NewSet(dependencies)
	if err != nil {
		return nil, err
	}

	resolved, err := sourceDs.ResolveTask(availableDs)
	if err != nil {
		return nil, err
	}

	localRunner, err := task.NewLocalRunner(resolved)
	if err != nil {
		return nil, err
	}
	localRunner = localRunner.WithCacheProvider(&task.LocalTaskVariableCache{})

	err = localRunner.Run(context.Background(), mode, initialParameters)
	if err != nil {
		return nil, err
	}

	<-localRunner.Wait()
	return localRunner.Result()
}
