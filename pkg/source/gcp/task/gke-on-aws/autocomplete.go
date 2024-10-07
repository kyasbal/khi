package aws

import (
	"context"
	"fmt"
	"log/slog"

	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var AutocompleteClusterNames = task.NewCachedProcessor(gcp_task.AutocompleteClusterNamesTaskId+"#anthos-on-aws", []string{
	gcp_task.GCPApiClientTaskId,
	gcp_task.InputProjectIdVariableName,
}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
	client, err := gcp_task.GetGCPApiClientFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	projectId, err := gcp_task.GetInputProjectIdFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	if projectId != "" {
		clusterNames, err := client.GetAnthosAWSClusterNames(ctx, projectId)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("Failed to read the cluster names in the project %s\n%s", projectId, err))
			return &gcp_task.AutocompleteClusterNameList{
				ClusterNames: []string{},
				Error:        "Failed to get the list from API",
			}, nil
		}
		return &gcp_task.AutocompleteClusterNameList{
			ClusterNames: clusterNames,
			Error:        "",
		}, nil
	}
	return &gcp_task.AutocompleteClusterNameList{
		ClusterNames: []string{},
		Error:        "Project ID is empty",
	}, nil
}, inspection_task.InspectionTaskLabel(InspectionTypeId))
