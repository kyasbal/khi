package azure

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
)

var AnthosOnAzureClusterNamePrefixTask = inspection_task.NewInspectionProducer(task.ClusterNamePrefixTaskId+"#gke-on-azure", func(ctx context.Context, taskMode int, progress *progress.TaskProgress) (any, error) {
	return "azureClusters/", nil
}, inspection_task.InspectionTaskLabel(InspectionTypeId))
