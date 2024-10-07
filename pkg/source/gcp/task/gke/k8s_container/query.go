package k8s_container

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query/queryutil"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func GenerateK8sContainerQuery(clusterName string, namespacesFilter *queryutil.SetFilterParseResult, podNamesFilter *queryutil.SetFilterParseResult) string {
	return fmt.Sprintf(`resource.type="k8s_container"
resource.labels.cluster_name="%s"
%s
%s`, clusterName, generateNamespacesFilter(namespacesFilter), generatePodNamesFilter(podNamesFilter))
}

func generateNamespacesFilter(namespacesFilter *queryutil.SetFilterParseResult) string {
	if namespacesFilter.ValidationError != "" {
		return fmt.Sprintf("-- Failed to generate namespaces filter due to the validation error \"%s\"", namespacesFilter.ValidationError)
	}
	if namespacesFilter.SubtractMode {
		if len(namespacesFilter.Subtractives) == 0 {
			return "-- No namespace filter"
		}
		namespacesWithQuotes := []string{}
		for _, namespace := range namespacesFilter.Subtractives {
			namespacesWithQuotes = append(namespacesWithQuotes, fmt.Sprintf(`"%s"`, namespace))
		}
		return fmt.Sprintf(`-resource.labels.namespace_name=(%s)`, strings.Join(namespacesWithQuotes, " OR "))
	}

	if len(namespacesFilter.Additives) == 0 {
		return `-- Invalid: none of the resources will be selected. Ignoreing kind filter.`
	}
	namespacesWithQuotes := []string{}
	for _, namespace := range namespacesFilter.Additives {
		namespacesWithQuotes = append(namespacesWithQuotes, fmt.Sprintf(`"%s"`, namespace))
	}
	return fmt.Sprintf(`resource.labels.namespace_name=(%s)`, strings.Join(namespacesWithQuotes, " OR "))

}

func generatePodNamesFilter(podNamesFilter *queryutil.SetFilterParseResult) string {
	if podNamesFilter.ValidationError != "" {
		return fmt.Sprintf("-- Failed to generate pod name filter due to the validation error \"%s\"", podNamesFilter.ValidationError)
	}
	if podNamesFilter.SubtractMode {
		if len(podNamesFilter.Subtractives) == 0 {
			return "-- No pod name filter"
		}
		podNamesWithQuotes := []string{}
		for _, podName := range podNamesFilter.Subtractives {
			podNamesWithQuotes = append(podNamesWithQuotes, fmt.Sprintf(`"%s"`, podName))
		}
		return fmt.Sprintf(`-resource.labels.pod_name:(%s)`, strings.Join(podNamesWithQuotes, " OR "))
	}

	if len(podNamesFilter.Additives) == 0 {
		return `-- Invalid: none of the resources will be selected. Ignoreing kind filter.`
	}
	podNamesWithQuotes := []string{}
	for _, podName := range podNamesFilter.Additives {
		podNamesWithQuotes = append(podNamesWithQuotes, fmt.Sprintf(`"%s"`, podName))
	}
	return fmt.Sprintf(`resource.labels.pod_name:(%s)`, strings.Join(podNamesWithQuotes, " OR "))
}

var GKEContainerLogQueryTaskId = query.GKEQueryPrefix + "k8s-container"
var GKEContainerQueryTask = query.NewQueryGeneratorTask(GKEContainerLogQueryTaskId, "K8s container logs", enum.LogTypeContainer, []string{
	gcp_task.InputClusterName,
	InputContainerQueryNamespacesVariableName,
	InputContainerQueryPodNamesVariableName,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	namespacesFilter, err := GetInputContainerQueryNamespacesFilterFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	podNamesFilter, err := GetInputContainerQueryPodNamesFilterFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return []string{GenerateK8sContainerQuery(clusterName, namespacesFilter, podNamesFilter)}, nil
})
