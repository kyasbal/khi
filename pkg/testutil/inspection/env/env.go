package env_test

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/env"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
	task_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/task"
)

// MockedEnvironmentVariableProducer creates a mocked producer task from the EnvironmentVariableProducer.
// The environment variable is resolved with the default value when the value is empty string.
// This is only for testing purpose.
func MockedEnvironmentVariableProducer(parentVariableDefinition task.Definition, value string) task.Definition {
	dv, found := parentVariableDefinition.Labels().Get(env.EnvironmentVariableDefaultValueLabel)
	if !found {
		panic(fmt.Errorf("the given parent variable definition is not declared with EnvironmentVariableProducer"))
	}
	if value == "" {
		return task_test.MockProcessorTaskFromTaskId(parentVariableDefinition.ID().String(), &env.EnvironmentVariable{
			Value:  dv.(string),
			Exists: false,
		})
	} else {
		return task_test.MockProcessorTaskFromTaskId(parentVariableDefinition.ID().String(), &env.EnvironmentVariable{
			Value:  value,
			Exists: true,
		})
	}
}
