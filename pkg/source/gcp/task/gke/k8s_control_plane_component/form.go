package k8scontrolplanecomponent

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/form"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query/queryutil"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const priorityForControlPlaneGroup = gcp_task.FormBasePriority + 30000

const InputControlPlaneComponentNameFilterTaskId = gcp_task.GCPPrefix + "input/component-names"

var inputControlPlaneComponentNameAliasMap map[string][]string = map[string][]string{}

var InputControlPlaneComponentNameFilterTask = form.NewInputFormDefinitionBuilder(
	InputControlPlaneComponentNameFilterTaskId,
	priorityForControlPlaneGroup+1000,
	"Control plane component names",
).
	WithDefaultValueConstant("@any", true).
	WithSuggestionsConstant([]string{
		"apiserver",
		"controller-manager",
		"scheduler",
	}).
	WithDescription("Control plane component names to query(e.g. apiserver, controller-manager...etc)").
	WithValidator(func(ctx context.Context, value string, variables *task.VariableSet) (string, error) {
		result, err := queryutil.ParseSetFilter(value, inputControlPlaneComponentNameAliasMap, true, true, true)
		if err != nil {
			return "", err
		}
		return result.ValidationError, nil
	}).
	WithConverter(func(ctx context.Context, value string, variables *task.VariableSet) (any, error) {
		result, err := queryutil.ParseSetFilter(value, inputControlPlaneComponentNameAliasMap, true, true, true)
		if err != nil {
			return "", err
		}
		return result, nil
	}).
	Build()

func GetInputControlPlaneComponentNameFilterFromTaskVariable(tv *task.VariableSet) (*queryutil.SetFilterParseResult, error) {
	return task.GetTypedVariableFromTaskVariable[*queryutil.SetFilterParseResult](tv, InputControlPlaneComponentNameFilterTaskId, nil)
}
