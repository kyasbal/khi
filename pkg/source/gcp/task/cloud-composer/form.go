package composer_task

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/form"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const InputProjectIdVariableName = gcp_task.GCPPrefix + "input/location"

var AutocompleteComposerEnvironmentNamesTaskId = gcp_task.GCPPrefix + "autocomplete/composer-environment-names"

var AutocompleteComposerEnvironmentNames = task.NewCachedProcessor(AutocompleteComposerEnvironmentNamesTaskId, []string{
	gcp_task.GCPApiClientTaskId,
	gcp_task.InputLocationsVariableName,
	gcp_task.InputProjectIdVariableName,
}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
	projectId, err := gcp_task.GetInputProjectIdFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	location, err := gcp_task.GetInputLocationsFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	client, err := gcp_task.GetGCPApiClientFromTaskVariable(v)
	if err != nil {
		return nil, err
	}

	if projectId != "" && location != "" {
		clusterNames, err := client.GetComposerEnvironmentNames(ctx, projectId, location)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("Failed to read the composer environments in the (project,location) (%s, %s) \n%s", projectId, location, err))
			return []string{}, nil
		}
		return clusterNames, nil
	}
	return []string{}, nil
})

func GetAutocompleteComposerEnvironmentNamesTaskVariable(v *task.VariableSet) ([]string, error) {
	return task.GetTypedVariableFromTaskVariable[[]string](v, AutocompleteComposerEnvironmentNamesTaskId, nil)
}

const InputComposerEnvironmentVariableName = gcp_task.GCPPrefix + "input/composer/environment_name"

var InputComposerEnvironmentNameTask = form.NewInputFormDefinitionBuilder(InputComposerEnvironmentVariableName, gcp_task.PriorityForResourceIdentifierGroup+5000, "Composer Environment Name").WithDependencies(
	[]string{AutocompleteComposerEnvironmentNamesTaskId},
).WithSuggestionsFunc(func(ctx context.Context, value string, variables *task.VariableSet, previousValues []string) ([]string, error) {
	environments, err := GetAutocompleteComposerEnvironmentNamesTaskVariable(variables)
	if err != nil {
		return nil, err
	}
	return common.SortForAutocomplete(value, environments), nil
}).Build()

func GetInputComposerEnvironmentVariable(tv *task.VariableSet) (string, error) {
	return task.GetTypedVariableFromTaskVariable[string](tv, InputComposerEnvironmentNameTask.ID().String(), "<INVALID>")
}
