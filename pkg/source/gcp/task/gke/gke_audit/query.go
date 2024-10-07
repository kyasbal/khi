package gke_audit

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func GenerateGKEAuditQuery(projectName string, clusterName string) string {
	return fmt.Sprintf(`resource.type=("gke_cluster" OR "gke_nodepool")
logName="projects/%s/logs/cloudaudit.googleapis.com%%2Factivity"
resource.labels.cluster_name="%s"`, projectName, clusterName)
}

var GKEAuditLogQueryTaskId = query.GKEQueryPrefix + "gke-audit"
var GKEAuditQueryTask = query.NewQueryGeneratorTask(GKEAuditLogQueryTaskId, "GKE Audit logs", enum.LogTypeGkeAudit, []string{
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

	return []string{GenerateGKEAuditQuery(projectId, clusterName)}, nil
})
