package ownerreferencerecorder

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/types"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func Register(manager *recorder.RecorderTaskManager) error {
	manager.AddRecorder("owner-references", []string{}, func(ctx context.Context, resourcePath string, currentLog *types.ResourceSpecificParserInput, prevStateInGroup any, cs *history.ChangeSet, builder *history.Builder, vs *task.VariableSet) (any, error) {
		return nil, recordChangeSetForLog(ctx, resourcePath, currentLog, cs, builder)
	}, recorder.AnyLogGroupFilter(), recorder.AndLogFilter(recorder.OnlySucceedLogs(), recorder.OnlyWithResourceBody()))
	return nil
}

func recordChangeSetForLog(ctx context.Context, resourcePath string, log *types.ResourceSpecificParserInput, cs *history.ChangeSet, builder *history.Builder) error {
	if !log.ResourceBodyReader.Has("metadata.ownerReferences") {
		return nil
	}
	ownerReferencesReaders, err := log.ResourceBodyReader.Reader("metadata.ownerReferences[]")
	if err != nil {
		return nil
	}
	for _, referenceReader := range ownerReferencesReaders {
		kind, err := referenceReader.ReadString("kind")
		if err != nil {
			continue
		}
		apiVersion, err := referenceReader.ReadString("apiVersion")
		if err != nil {
			continue
		}
		name, err := referenceReader.ReadString("name")
		if err != nil {
			continue
		}
		if !strings.Contains(apiVersion, "/") {
			apiVersion = "core/" + apiVersion
		}
		namespace := log.Operation.Namespace
		// TODO: Usually ownerReference don't contain the namespace field but the owner should be in the same namespace.
		// But node is a cluster scopd resource. There should be better implementation rather than hard coding this rule here.
		if kind == "Node" {
			namespace = "cluster-scope"
		}
		cs.RecordAliasRelationship(log.Operation.CovertToResourcePath(), resourcepath.SubresourceLayerGeneralItem(
			apiVersion, strings.ToLower(kind), namespace, name, fmt.Sprintf("%s(%s)[%s]", log.Operation.Name, log.Operation.Namespace, log.Operation.GetSingularKindName()),
		), history.RewriteRelationship(enum.RelationshipOwnerReference))
	}
	return nil
}
