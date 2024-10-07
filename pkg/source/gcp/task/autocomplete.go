package task

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var AutocompleteClusterNamesTaskId = GCPPrefix + "autocomplete/cluster-names"

type AutocompleteClusterNameList struct {
	ClusterNames []string
	Error        string
}

func GetAutocompleteClusterNamesFromTaskVariable(v *task.VariableSet) (*AutocompleteClusterNameList, error) {
	return task.GetTypedVariableFromTaskVariable[*AutocompleteClusterNameList](v, AutocompleteClusterNamesTaskId, nil)
}
