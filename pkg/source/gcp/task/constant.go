package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/k8s"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const GCPPrefix = "cloud.google.com/"

// ClusterNamePrefixTaskId is the task ID for generating the cluster name prefix used in query.
// For GKE, it's just a task to return "" always.
// For Anthos on AWS, it should return "awsClusters/" because the `resource.labels.cluster_name` field would be `awsClusters/<cluster-name>`
// For Anthos on Azure, it will be "azureClusters/"
const ClusterNamePrefixTaskId = GCPPrefix + "cluster-name-prefix"

func GetClusterNamePrefixFromTaskVariable(v *task.VariableSet) (string, error) {
	return task.GetTypedVariableFromTaskVariable[string](v, ClusterNamePrefixTaskId, "")
}

const K8sResourceMergeConfigTaskId = GCPPrefix + "merge-config"

func GetK8sResourceMergeConfigFromTaskVariable(v *task.VariableSet) (*k8s.MergeConfigRegistry, error) {
	return task.GetTypedVariableFromTaskVariable[*k8s.MergeConfigRegistry](v, K8sResourceMergeConfigTaskId, nil)
}

var GCPDefaultK8sResourceMergeConfigTask = inspection_task.NewInspectionProducer(K8sResourceMergeConfigTaskId+"#gcp", func(ctx context.Context, taskMode int, progress *progress.TaskProgress) (any, error) {
	return k8s.GenerateDefaultMergeConfig()
})
