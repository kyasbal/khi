package env

import (
	"context"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const EnvironmentVariableKeyLabel = task.KHISystemPrefix + "/labels/environment-variable-key"
const EnvironmentVariableDefaultValueLabel = task.KHISystemPrefix + "/labels/environment-variable-default-value"

type EnvironmentVariable struct {
	Value  string
	Exists bool
}

// Digest implements task.CachableDependency.
func (ev *EnvironmentVariable) Digest() string {
	return fmt.Sprintf("%t-%s", ev.Exists, ev.Value)
}

var _ task.CachableDependency = (*EnvironmentVariable)(nil)

func GetEnvironmentVariableFromTaskVariables(ctx context.Context, environmentVariableName string, v *task.VariableSet) (*EnvironmentVariable, error) {
	envAny, err := v.Get(environmentVariableName)
	if err != nil {
		return nil, err
	}
	return envAny.(*EnvironmentVariable), nil
}

// EnvironmentVariableProducer creates a producer task from given arguments that is resolved with environment variable values
func EnvironmentVariableProducer(taskId string, environmentKey string, defaultValue string) task.Definition {
	return task.NewProcessorTask(taskId, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		value, exists := os.LookupEnv(environmentKey)
		if !exists {
			value = defaultValue
		}
		return &EnvironmentVariable{
			Value:  value,
			Exists: exists,
		}, nil
	}, task.WithLabel(EnvironmentVariableKeyLabel, environmentKey), task.WithLabel(EnvironmentVariableDefaultValueLabel, defaultValue))
}
