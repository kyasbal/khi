package k8scontrolplanecomponent

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

func GenerateK8sControlPlaneQuery(clusterName string, projectId string, controlplaneComponentFilter *queryutil.SetFilterParseResult) string {
	return fmt.Sprintf(`resource.type="k8s_control_plane_component"
resource.labels.cluster_name="%s"
resource.labels.project_id="%s"
-sourceLocation.file="httplog.go"
%s`, clusterName, projectId, generateK8sControlPlaneComponentFilter(controlplaneComponentFilter))
}

const GKEK8sControlPlaneComponentQueryTaskId = query.GKEQueryPrefix + "k8s-controlplane"

var GKEK8sControlPlaneLogQueryTask = query.NewQueryGeneratorTask(GKEK8sControlPlaneComponentQueryTaskId, "K8s control plane logs", enum.LogTypeControlPlaneComponent, []string{
	gcp_task.InputProjectIdVariableName,
	gcp_task.InputClusterName,
	InputControlPlaneComponentNameFilterTaskId,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	projectId, err := gcp_task.GetInputProjectIdFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	controlPlaneComponentNameFilter, err := GetInputControlPlaneComponentNameFilterFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return []string{GenerateK8sControlPlaneQuery(clusterName, projectId, controlPlaneComponentNameFilter)}, nil
})

func generateK8sControlPlaneComponentFilter(filter *queryutil.SetFilterParseResult) string {
	if filter.ValidationError != "" {
		return fmt.Sprintf(`-- Failed to generate component name filter due to the validation error "%s"`, filter.ValidationError)
	}
	if filter.SubtractMode {
		if len(filter.Subtractives) == 0 {
			return "-- No component name filter"
		}
		return fmt.Sprintf(`-resource.labels.component_name:(%s)`, strings.Join(queryutil.WrapDoubleQuoteForStringArray(filter.Subtractives), " OR "))
	} else {
		if len(filter.Additives) == 0 {
			return `-- Invalid: none of the controlplane component will be selected. Ignoreing component name filter.`
		}
		return fmt.Sprintf(`resource.labels.component_name:(%s)`, strings.Join(queryutil.WrapDoubleQuoteForStringArray(filter.Additives), " OR "))
	}
}
