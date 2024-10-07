package autoscaler

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var AutoscalerQueryTaskId = query.GKEQueryPrefix + "autoscaler"

func GenerateAutoscalerQuery(projectId string, clusterName string, excludeStatus bool) string {
	excludeStatusQueryFragment := "-- include query for status log"
	if excludeStatus {
		excludeStatusQueryFragment = `-jsonPayload.status: ""`
	}
	return fmt.Sprintf(`resource.type="k8s_cluster"
resource.labels.project_id="%s"
resource.labels.cluster_name="%s"
%s
logName="projects/%s/logs/container.googleapis.com%%2Fcluster-autoscaler-visibility"`, projectId, clusterName, excludeStatusQueryFragment, projectId)
}

var AutoscalerQueryTask = query.NewQueryGeneratorTask(AutoscalerQueryTaskId, "Autoscaler logs", enum.LogTypeAutoscaler, []string{
	gcp_task.InputProjectIdVariableName,
	gcp_task.InputClusterName,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	projectId, err := gcp_task.GetInputProjectIdFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return []string{GenerateAutoscalerQuery(projectId, clusterName, true)}, nil
})
