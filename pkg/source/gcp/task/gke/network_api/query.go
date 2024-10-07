package network_api

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

const GCPNetworkLogQueryTaskId = query.GKEQueryPrefix + "network-api"

func GenerateGCPNetworkAPIQuery(taskMode int, negNames []string) []string {
	nodeNamesWithNetworkEndpointGroups := []string{}
	for _, nodeName := range negNames {
		nodeNamesWithNetworkEndpointGroups = append(nodeNamesWithNetworkEndpointGroups, fmt.Sprintf("networkEndpointGroups/%s", nodeName))
	}
	if taskMode == inspection_task.TaskModeDryRun {
		return []string{queryFromNegNameFilter("-- neg name filters to be determined after audit log query")}
	} else {
		result := []string{}
		groups := queryutil.SplitToChildGroups(nodeNamesWithNetworkEndpointGroups, 10)
		for _, group := range groups {
			negNameFilter := fmt.Sprintf("protoPayload.resourceName:(%s)", strings.Join(group, " OR "))
			result = append(result, queryFromNegNameFilter(negNameFilter))
		}
		return result
	}
}

func queryFromNegNameFilter(negNameFilter string) string {
	return fmt.Sprintf(`resource.type="gce_network"
-protoPayload.methodName:("list" OR "get" OR "watch")
%s
`, negNameFilter)
}

var GCPNetworkLogQueryTask = query.NewQueryGeneratorTask(GCPNetworkLogQueryTaskId, "GCP network log", enum.LogTypeNetworkAPI, []string{
	k8saudittask.K8sAuditParseTaskId,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	builder, err := inspection_task.GetHistoryBuilderFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return GenerateGCPNetworkAPIQuery(i, builder.ClusterResource.NEGs.GetAllIdentifiers()), nil
})
