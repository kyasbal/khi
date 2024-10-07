package compute_api

import (
	"context"
	"fmt"
	"strings"

	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query/queryutil"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/k8saudittask"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var ComputeAPIQueryTaskId = query.GKEQueryPrefix + "compute-api"

func GenerateComputeAPIQuery(taskMode int, nodeNames []string) []string {
	if taskMode == inspection_task.TaskModeDryRun {
		return []string{
			generateComputeAPIQueryWithInstanceNameFilter("-- instance name filters to be determined after audit log query"),
		}
	} else {
		result := []string{}
		instanceNameGroups := queryutil.SplitToChildGroups(nodeNames, 30)
		for _, group := range instanceNameGroups {
			nodeNamesWithInstance := []string{}
			for _, nodeName := range group {
				nodeNamesWithInstance = append(nodeNamesWithInstance, fmt.Sprintf("instances/%s", nodeName))
			}
			instanceNameFilter := fmt.Sprintf("protoPayload.resourceName:(%s)", strings.Join(nodeNamesWithInstance, " OR "))
			result = append(result, generateComputeAPIQueryWithInstanceNameFilter(instanceNameFilter))
		}
		return result
	}
}

func generateComputeAPIQueryWithInstanceNameFilter(instanceNameFilter string) string {
	return fmt.Sprintf(`resource.type="gce_instance"
	-protoPayload.methodName:("list" OR "get" OR "watch")
	%s
	`, instanceNameFilter)
}

var ComputeAPIQueryTask = query.NewQueryGeneratorTask(ComputeAPIQueryTaskId, "Compute API Logs", enum.LogTypeComputeApi, []string{
	k8saudittask.K8sAuditParseTaskId,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	builder, err := inspection_task.GetHistoryBuilderFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return GenerateComputeAPIQuery(i, builder.ClusterResource.GetNodes()), nil
})
