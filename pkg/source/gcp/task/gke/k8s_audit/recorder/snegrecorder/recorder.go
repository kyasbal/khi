package snegrecorder

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourceinfo/resourcelease"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/types"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func Register(manager *recorder.RecorderTaskManager) error {
	manager.AddRecorder("sneg-fields", []string{}, func(ctx context.Context, resourcePath string, currentLog *types.ResourceSpecificParserInput, prevStateInGroup any, cs *history.ChangeSet, builder *history.Builder, vs *task.VariableSet) (any, error) {
		// record node name for querying compute engine api later.
		builder.ClusterResource.NEGs.TouchResourceLease(currentLog.Operation.Name, currentLog.Log.Timestamp(), resourcelease.NewK8sResourceLeaseHolder(
			currentLog.Operation.PluralKind,
			currentLog.Operation.Namespace,
			currentLog.Operation.Name,
		))
		return nil, nil
	}, recorder.ResourceKindLogGroupFilter("servicenetworkendpointgroup"), recorder.AndLogFilter(recorder.OnlySucceedLogs(), recorder.OnlyWithResourceBody()))
	return nil
}
