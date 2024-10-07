package k8s_node

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func GenerateK8sNodeLogQuery(projectId string, clusterId string) string {
	return fmt.Sprintf(`resource.type="k8s_node"
-logName="projects/%s/logs/events"
resource.labels.cluster_name="%s"
`, projectId, clusterId)
}

const GKENodeLogQueryTaskId = query.GKEQueryPrefix + "k8s-node"

var GKENodeQueryTask = query.NewQueryGeneratorTask(GKENodeLogQueryTaskId, "Kubernetes node log", enum.LogTypeNode, []string{
	gcp_task.InputProjectIdVariableName,
	gcp_task.InputClusterName,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	projectId, err := gcp_task.GetInputProjectIdFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return []string{GenerateK8sNodeLogQuery(projectId, clusterName)}, nil
})
