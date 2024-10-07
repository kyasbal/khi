package network_api

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/grouper"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parser"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	composer_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/cloud-composer"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
	"gopkg.in/yaml.v3"
)

type gceNetworkParser struct{}

// Dependencies implements parser.Parser.
func (*gceNetworkParser) Dependencies() []string {
	return []string{}
}

// Description implements parser.Parser.
func (*gceNetworkParser) Description() string {
	return `GCE network API audit log including NEG related audit logs to identify when the associated NEG was attached/detached.`
}

// GetParserName implements parser.Parser.
func (*gceNetworkParser) GetParserName() string {
	return "GCE Network Logs"
}

// LogTask implements parser.Parser.
func (*gceNetworkParser) LogTask() string {
	return GCPNetworkLogQueryTaskId
}

func (*gceNetworkParser) Grouper() grouper.LogGrouper {
	return grouper.AllDependentLogGrouper
}

// Parse implements parser.Parser.
func (*gceNetworkParser) Parse(ctx context.Context, l *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, variables *task.VariableSet) error {
	isFirst := l.Has("operation.first")
	isLast := l.Has("operation.last")
	operationId := l.GetStringOrDefault("operation.id", "unknown")
	methodName := l.GetStringOrDefault("protoPayload.methodName", "unknown")
	methodNameSplitted := strings.Split(methodName, ".")
	resourceName := l.GetStringOrDefault("protoPayload.resourceName", "unknown")
	resourceNameSplitted := strings.Split(resourceName, "/")
	negName := resourceNameSplitted[len(resourceNameSplitted)-1]
	principal := l.GetStringOrDefault("protoPayload.authenticationInfo.principalEmail", "unknown")
	resourcePath := fmt.Sprintf("networking.gke.io/v1beta1#servicenetworkendpointgroup#unknown#%s", negName)
	lease, err := builder.ClusterResource.NEGs.GetResourceLeaseHolderAt(negName, l.Timestamp())
	if err == nil {
		resourcePath = fmt.Sprintf("networking.gke.io/v1beta1#servicenetworkendpointgroup#%s#%s", lease.Holder.Namespace, negName)
	}
	if !(isLast && isFirst) && (isLast || isFirst) {
		state := enum.RevisionStateOperationStarted
		if isLast {
			state = enum.RevisionStateOperationFinished
		}
		operationPath := fmt.Sprintf("%s#%s-%s", resourcePath, methodNameSplitted[len(methodNameSplitted)-1], operationId)
		cs.RecordRevision(operationPath, &history.StagingResourceRevision{
			Verb:       enum.RevisionVerbCreate,
			State:      state,
			Requestor:  principal,
			ChangeTime: l.Timestamp(),
			Partial:    false,
		}, history.RewriteRelationship(enum.RelationshipOperation))
	} else {
		cs.RecordEvent(resourcePath)
	}
	if isFirst {
		method := methodNameSplitted[len(methodNameSplitted)-1]
		if method == "detachNetworkEndpoints" || method == "attachNetworkEndpoints" {
			isDetach := strings.HasPrefix(method, "detach")
			requestBody, err := l.GetChildYamlOf("protoPayload.request")
			if err != nil {
				return err
			}
			var negRequest NegAttachOrDetachRequest
			err = yaml.Unmarshal([]byte(requestBody), &negRequest)
			if err != nil {
				return err
			}
			for _, endpoint := range negRequest.NetworkEndpoints {
				lease, err := builder.ClusterResource.IPs.GetResourceLeaseHolderAt(endpoint.IpAddress, l.Timestamp())
				if err != nil {
					slog.WarnContext(ctx, fmt.Sprintf("Failed to identify the holder of the IP %s.\n This might be because the IP holder resource wasn't updated during the log period ", endpoint.IpAddress))
					continue
				}
				holder := lease.Holder
				if holder.Kind == "pod" {
					podPath := fmt.Sprintf("core/v1#pod#%s#%s", holder.Namespace, holder.Name)
					negSubresourcePath := fmt.Sprintf("%s#%s", podPath, negName)
					state := enum.RevisionStateConditionTrue
					verb := enum.RevisionVerbReady
					if isDetach {
						state = enum.RevisionStateConditionFalse
						verb = enum.RevisionVerbNonReady
					}
					cs.RecordRevision(negSubresourcePath, &history.StagingResourceRevision{
						Verb:       verb,
						State:      state,
						Requestor:  principal,
						ChangeTime: l.Timestamp(),
						Partial:    false,
					}, history.RewriteRelationship(enum.RelationshipNetworkEndpointGroup))
				}
			}
		}
	}
	if isFirst && !isLast {
		cs.RecordLogSummary(fmt.Sprintf("%s Started", methodName))
	} else if !isFirst && isLast {
		cs.RecordLogSummary(fmt.Sprintf("%s Finished", methodName))
	} else {
		cs.RecordLogSummary(methodName)
	}

	return nil
}

var _ parser.Parser = (*gceNetworkParser)(nil)

var NetowrkAPIParserTask = parser.NewParserTaskFromParser(gcp_task.GCPPrefix+"feature/network-api-parser", &gceNetworkParser{}, true, inspection_task.InspectionTaskLabel(gke.InspectionTypeId, composer_task.InspectionTypeId))
