package bindingrecorder

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/types"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func Register(manager *recorder.RecorderTaskManager) error {
	manager.AddRecorder("binding", []string{}, func(ctx context.Context, resourcePath string, l *types.ResourceSpecificParserInput, prevState any, cs *history.ChangeSet, builder *history.Builder, vs *task.VariableSet) (any, error) {
		return nil, recordChangeSetForLog(ctx, resourcePath, l, cs)
	}, recorder.SubresourceLogGroupFilter("binding"), recorder.AnyLogFilter())
	return nil
}

func recordChangeSetForLog(ctx context.Context, resourcePath string, log *types.ResourceSpecificParserInput, cs *history.ChangeSet) error {
	if log.ResourceBodyReader == nil {
		return nil
	}
	target := log.ResourceBodyReader.ReadStringOrDefault("target.name", "unknown")

	podK8sOp := model.KubernetesObjectOperation{
		APIVersion: log.Operation.APIVersion,
		PluralKind: log.Operation.PluralKind,
		Namespace:  log.Operation.Namespace,
		Name:       log.Operation.Name,
	}
	podScheduledStatusPath := fmt.Sprintf("%s#PodScheduled", podK8sOp.CovertToResourcePath())
	if log.Operation.Verb == enum.RevisionVerbCreate {
		cs.RecordRevision(resourcepath.NodeBinding(target, log.Operation.Namespace, log.Operation.Name), &history.StagingResourceRevision{
			Verb:       enum.RevisionVerbCreate,
			Body:       log.ResourceBodyYaml,
			Partial:    false,
			Requestor:  log.PrincipalEmail,
			ChangeTime: log.Log.Timestamp(),
			State:      enum.RevisionStateExisting,
		}, history.RewriteRelationship(enum.RelationshipPodBinding))
		cs.RecordRevision(podScheduledStatusPath, &history.StagingResourceRevision{
			Verb:       enum.RevisionVerbStatusTrue,
			Body:       "# PodScheduled status was inferred to be `True` from a binding resource",
			Partial:    false,
			Requestor:  "",
			ChangeTime: log.Log.Timestamp(),
			State:      enum.RevisionStateConditionTrue,
		}, history.RewriteRelationship(enum.RelationshipResourceStatus))
	} else {
		cs.RecordRevision(resourcepath.NodeBinding(target, log.Operation.Namespace, log.Operation.Name), &history.StagingResourceRevision{
			Verb:       enum.RevisionVerbDelete,
			Body:       log.ResourceBodyYaml,
			Partial:    false,
			Requestor:  log.PrincipalEmail,
			ChangeTime: log.Log.Timestamp(),
			State:      enum.RevisionStateDeleted,
		}, history.RewriteRelationship(enum.RelationshipPodBinding))
		cs.RecordRevision(podScheduledStatusPath, &history.StagingResourceRevision{
			Verb:       enum.RevisionVerbStatusFalse,
			Body:       "# PodScheduled status was inferred to be `False` from a binding resource",
			Partial:    false,
			Requestor:  "",
			ChangeTime: log.Log.Timestamp(),
			State:      enum.RevisionStateConditionFalse,
		}, history.RewriteRelationship(enum.RelationshipResourceStatus))
	}
	return nil
}
