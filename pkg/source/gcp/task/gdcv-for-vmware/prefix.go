package vmware

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
)

var AnthosOnVMWareClusterNamePrefixTask = inspection_task.NewInspectionProducer(task.ClusterNamePrefixTaskId+"#gdcv-for-vmware", func(ctx context.Context, taskMode int, progress *progress.TaskProgress) (any, error) {
	return "", nil
}, inspection_task.InspectionTaskLabel(InspectionTypeId))
