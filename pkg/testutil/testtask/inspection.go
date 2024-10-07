package testtask

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func RunSingleTask[T any](target task.Definition, mode int, opts ...TestRunTaskParameterOpt) (T, error) {
	return RunMultipleTask[T](target, []task.Definition{}, mode, opts...)
}

func RunMultipleTask[T any](target task.Definition, availableTasks []task.Definition, mode int, opts ...TestRunTaskParameterOpt) (T, error) {
	params := generateVariableSetFromOpts(opts...)
	sourceTaskSet, err := task.NewSet([]task.Definition{target})
	if err != nil {
		return *new(T), err
	}

	mockedParameterTasks := []task.Definition{}
	for key, value := range params {
		nextTaskValue := value
		mockedParameterTasks = append(mockedParameterTasks, task.NewProcessorTask(key, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
			return nextTaskValue, nil
		}))
	}

	availableTaskSet, err := task.NewSet(append(availableTasks, mockedParameterTasks...))
	if err != nil {
		return *new(T), err
	}

	resolved, err := sourceTaskSet.ResolveTask(availableTaskSet)
	if err != nil {
		return *new(T), err
	}

	localRunner, err := task.NewLocalRunner(resolved)
	if err != nil {
		return *new(T), err
	}

	localRunner = localRunner.WithCacheProvider(&task.LocalTaskVariableCache{})

	err = localRunner.Run(context.Background(), mode, map[string]any{
		inspection_task.MetadataVariableName: metadata.NewSet(),
	})
	if err != nil {
		return *new(T), err
	}

	<-localRunner.Wait()
	result, err := localRunner.Result()
	if err != nil {
		return *new(T), err
	}
	return task.GetTypedVariableFromTaskVariable(result, target.ID().String(), *new(T))
}

func generateVariableSetFromOpts(opts ...TestRunTaskParameterOpt) map[string]any {
	parameters := map[string]any{}
	parameters[inspection_task.MetadataVariableName] = metadata.NewSet()
	for _, opt := range opts {
		opt.AddParam(parameters)
	}
	return parameters
}
