package compute_api

import (
	"context"
	"fmt"
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
)

type computeAPIParser struct {
}

// Dependencies implements parser.Parser.
func (*computeAPIParser) Dependencies() []string {
	return []string{}
}

// Description implements parser.Parser.
func (*computeAPIParser) Description() string {
	return `Compute API audit logs used for cluster related logs. This also visualize operations happened during the query time.`
}

// GetParserName implements parser.Parser.
func (*computeAPIParser) GetParserName() string {
	return `Compute API Logs`
}

// LogTask implements parser.Parser.
func (*computeAPIParser) LogTask() string {
	return ComputeAPIQueryTaskId
}
func (*computeAPIParser) Grouper() grouper.LogGrouper {
	return grouper.AllDependentLogGrouper
}

// Parse implements parser.Parser.
func (*computeAPIParser) Parse(ctx context.Context, l *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, variables *task.VariableSet) error {
	isFirst := l.Has("operation.first")
	isLast := l.Has("operation.last")
	operationId := l.GetStringOrDefault("operation.id", "unknown")
	methodName := l.GetStringOrDefault("protoPayload.methodName", "unknown")
	methodNameSplitted := strings.Split(methodName, ".")
	resourceName := l.GetStringOrDefault("protoPayload.resourceName", "unknown")
	resourceNameSplitted := strings.Split(resourceName, "/")
	instanceName := resourceNameSplitted[len(resourceNameSplitted)-1]
	principal := l.GetStringOrDefault("protoPayload.authenticationInfo.principalEmail", "unknown")
	resourcePath := fmt.Sprintf("core/v1#node#cluster-scope#%s", instanceName)
	// If this was an operation, it will be recorded as operation data
	if !(isLast && isFirst) && (isLast || isFirst) {
		state := enum.RevisionStateOperationStarted
		verb := enum.RevisionVerbOperationStart
		if isLast {
			state = enum.RevisionStateOperationFinished
			verb = enum.RevisionVerbOperationFinish
		}
		operationPath := fmt.Sprintf("%s#%s-%s", resourcePath, methodNameSplitted[len(methodNameSplitted)-1], operationId)
		cs.RecordRevision(operationPath, &history.StagingResourceRevision{
			Verb:       verb,
			State:      state,
			Requestor:  principal,
			ChangeTime: l.Timestamp(),
			Partial:    false,
		}, history.RewriteRelationship(enum.RelationshipOperation))
	}
	cs.RecordEvent(resourcePath)

	if isFirst && !isLast {
		cs.RecordLogSummary(fmt.Sprintf("%s Started", methodName))
	} else if !isFirst && isLast {
		cs.RecordLogSummary(fmt.Sprintf("%s Finished", methodName))
	} else {
		cs.RecordLogSummary(methodName)
	}

	return nil
}

var _ parser.Parser = (*computeAPIParser)(nil)

var ComputeAPIParserTask = parser.NewParserTaskFromParser(gcp_task.GCPPrefix+"feature/compute-api-parser", &computeAPIParser{}, true, inspection_task.InspectionTaskLabel(gke.InspectionTypeId, composer_task.InspectionTypeId))
