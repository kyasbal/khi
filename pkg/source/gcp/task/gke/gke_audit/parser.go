package gke_audit

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
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parser"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	composer_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/cloud-composer"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

type gkeAuditLogParser struct {
}

// Dependencies implements parser.Parser.
func (*gkeAuditLogParser) Dependencies() []string {
	return []string{}
}

// Description implements parser.Parser.
func (*gkeAuditLogParser) Description() string {
	return `GKE audit log including cluster creation,deletion and upgrades.`
}

// GetParserName implements parser.Parser.
func (*gkeAuditLogParser) GetParserName() string {
	return `GKE Audit logs`
}

// LogTask implements parser.Parser.
func (*gkeAuditLogParser) LogTask() string {
	return GKEAuditLogQueryTaskId
}

func (*gkeAuditLogParser) Grouper() grouper.LogGrouper {
	return grouper.AllDependentLogGrouper
}

// Parse implements parser.Parser.
func (p *gkeAuditLogParser) Parse(ctx context.Context, l *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, variables *task.VariableSet) error {
	clusterName := l.GetStringOrDefault("resource.labels.cluster_name", "unknown")
	nodepoolName, err := getRelatedNodepool(l)
	isFirst := l.Has("operation.first")
	isLast := l.Has("operation.last")
	operationId := l.GetStringOrDefault("operation.id", "unknown")
	methodName := l.GetStringOrDefault("protoPayload.methodName", "unknown")
	principal := l.GetStringOrDefault("protoPayload.authenticationInfo.principalEmail", "unknown")
	operationResourcePath := ""
	if err != nil {
		// assume this is a cluster operation
		clusterResourcePath := resourcepath.Cluster(clusterName)
		if strings.HasSuffix(methodName, "CreateCluster") && isFirst {
			body, err := l.GetChildYamlOf("protoPayload.request.cluster")
			if err != nil {
				slog.WarnContext(ctx, fmt.Sprintf("Failed to get the cluster info from the log\n%v", err))
			}
			cs.RecordRevision(clusterResourcePath, &history.StagingResourceRevision{
				Verb:       enum.RevisionVerbCreate,
				State:      enum.RevisionStateExisting,
				Requestor:  principal,
				ChangeTime: l.Timestamp(),
				Partial:    false,
				Body:       body,
			})
		}
		if strings.HasSuffix(methodName, "DeleteCluster") && isFirst {
			cs.RecordRevision(clusterResourcePath, &history.StagingResourceRevision{
				Verb:       enum.RevisionVerbDelete,
				State:      enum.RevisionStateDeleted,
				Requestor:  principal,
				ChangeTime: l.Timestamp(),
				Partial:    false,
				Body:       "",
			})
		}
		methodNameSplitted := strings.Split(methodName, ".")
		methodVerb := methodNameSplitted[len(methodNameSplitted)-1]
		operationResourcePath = fmt.Sprintf("%s#%s-%s", clusterResourcePath, methodVerb, operationId)
		cs.RecordEvent(clusterResourcePath)
	} else {
		nodepoolResourcePath := fmt.Sprintf("@Cluster#nodepools#%s#%s", clusterName, nodepoolName)
		if strings.HasSuffix(methodName, "CreateNodePool") && isFirst {
			body, err := l.GetChildYamlOf("protoPayload.request.nodePool")
			if err != nil {
				slog.WarnContext(ctx, fmt.Sprintf("Failed to get the nodepool info from the log\n%v", err))
			}
			cs.RecordRevision(nodepoolResourcePath, &history.StagingResourceRevision{
				Verb:       enum.RevisionVerbCreate,
				State:      enum.RevisionStateExisting,
				Requestor:  principal,
				ChangeTime: l.Timestamp(),
				Partial:    false,
				Body:       body,
			})
		}
		if strings.HasSuffix(methodName, "DeleteNodePool") && isFirst {
			cs.RecordRevision(nodepoolResourcePath, &history.StagingResourceRevision{
				Verb:       enum.RevisionVerbDelete,
				State:      enum.RevisionStateDeleted,
				Requestor:  principal,
				ChangeTime: l.Timestamp(),
				Partial:    false,
				Body:       "",
			})
		}
		cs.RecordEvent(nodepoolResourcePath)
		methodNameSplitted := strings.Split(methodName, ".")
		methodVerb := methodNameSplitted[len(methodNameSplitted)-1]
		operationResourcePath = fmt.Sprintf("%s#%s-%s", nodepoolResourcePath, methodVerb, operationId)
	}

	// If this was an operation, it will be recorded as operation data
	if !(isLast && isFirst) && (isLast || isFirst) {
		state := enum.RevisionStateOperationStarted
		verb := enum.RevisionVerbOperationStart
		if isLast {
			state = enum.RevisionStateOperationFinished
			verb = enum.RevisionVerbOperationFinish
		}
		cs.RecordRevision(operationResourcePath, &history.StagingResourceRevision{
			Verb:       verb,
			State:      state,
			Requestor:  principal,
			ChangeTime: l.Timestamp(),
			Partial:    false,
		}, history.RewriteRelationship(enum.RelationshipOperation))
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

func getRelatedNodepool(l *log.LogEntity) (string, error) {
	nodepoolName, err := l.GetString("resource.labels.nodepool_name")
	if err == nil {
		return nodepoolName, nil
	}
	return l.GetString("protoPayload.request.update.desiredNodePoolId")
}

var _ parser.Parser = (*gkeAuditLogParser)(nil)

var GKEAuditLogParseJob = parser.NewParserTaskFromParser(gcp_task.GCPPrefix+"feature/gke-audit-parser", &gkeAuditLogParser{}, true, inspection_task.InspectionTaskLabel(gke.InspectionTypeId, composer_task.InspectionTypeId))
