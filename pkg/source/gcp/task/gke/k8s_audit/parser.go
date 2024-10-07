package k8s_audit

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/adapter"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/grouper"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourceinfo/resourcelease"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	model_k8s "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/k8s"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parser"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/k8s"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/k8saudittask"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
	"gopkg.in/yaml.v3"
	goyaml "gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
)

var GKEK8sAuditLogParseJob = parser.NewParserTaskFromParser(gcp_task.GCPPrefix+"feature/audit-parser", NewGKEAuditLogParser(), true)

type k8sAuditLogParser struct {
}

// gkeAuditLogParser implements parser.Parser
var _ parser.Parser = (*k8sAuditLogParser)(nil)

type K8SOperationTemplateInput struct {
	Log          *log.LogEntity
	K8SOperation *model.KubernetesObjectOperation
}

func (p *k8sAuditLogParser) GetParserName() string {
	return "Kubernetes Audit Logs"
}

func (p *k8sAuditLogParser) Description() string {
	return `Visualize Kubernetes audit logs in GKE. 
This parser reveals how these resources are created,updated or deleted`
}

func (p *k8sAuditLogParser) Dependencies() []string {
	return []string{
		gcp_task.InputStartTimeVariableName,
		gcp_task.K8sResourceMergeConfigTaskId,
	}
}

func (p *k8sAuditLogParser) LogTask() string {
	return k8saudittask.K8sAuditQueryTaskId
}

func (*k8sAuditLogParser) Grouper() grouper.LogGrouper {
	return grouper.AllDependentLogGrouper
}

func (p *k8sAuditLogParser) Parse(ctx context.Context, inputLog *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, v *task.VariableSet) error {
	mergerConfig, err := gcp_task.GetK8sResourceMergeConfigFromTaskVariable(v)
	if err != nil {
		return err
	}
	startTime, err := gcp_task.GetInputStartTimeFromTaskVariable(v)
	if err != nil {
		return err
	}

	resourceName, err := inputLog.GetString("protoPayload.resourceName")
	if err != nil {
		return fmt.Errorf("resourceName not found")
	}

	methodName, err := inputLog.GetString("protoPayload.methodName")
	if err != nil {
		return fmt.Errorf("methodName not found")
	}

	authenticatedPrincipal, err := inputLog.GetString("protoPayload.authenticationInfo.principalEmail")
	if err != nil {
		return fmt.Errorf("principalPayload not found")
	}

	readerFactory, err := inspection_task.GetReaderFactoryFromTaskVariable(v)
	if err != nil {
		return err
	}

	k8sOperation := k8s.ParseKubernetesOperation(resourceName, methodName)
	if k8sOperation.SubResourceName == "status" {
		k8sOperation.SubResourceName = ""
	}

	if k8sOperation.APIVersion == "core/v1" && k8sOperation.PluralKind == "node" {
		handleNodeResource(k8sOperation, builder)
	}

	code := inputLog.GetIntOrDefault("protoPayload.status.code", 0)
	if code != 0 {
		message := inputLog.GetStringOrDefault("protoPayload.status.message", "Unknown")
		cs.RecordEvent(k8sOperation.CovertToResourcePath())
		cs.RecordLogSeverity(enum.SeverityError)
		cs.RecordLogSummary(fmt.Sprintf("【%s】%s", message, methodName))
		return nil
	}

	logSummary := fmt.Sprintf("%s on %s.%s.%s(%s in %s)", enum.RevisionVerbs[k8sOperation.Verb].Label, k8sOperation.Namespace, k8sOperation.Name, k8sOperation.SubResourceName, k8sOperation.PluralKind, k8sOperation.APIVersion)
	cs.RecordLogSummary(logSummary)

	if k8sOperation.Verb == enum.RevisionVerbDeleteCollection {
		err := deleteResourcesUnderNamespace(cs, builder, k8sOperation, inputLog, authenticatedPrincipal)
		if err != nil {
			return err
		}
		return nil
	}

	// Parse from response field
	responseReader, err := inputLog.Fields.ReaderSingle("protoPayload.response")
	if err == nil {
		responseKind, err := responseReader.ReadString("kind")
		if err == nil && responseKind != "Status" { // when kube-apiserver didn't respond entire resource
			return revisionsFromResponseField(ctx, readerFactory, cs, inputLog, builder, k8sOperation, authenticatedPrincipal, startTime)
		}
	}

	// Parse from request field
	_, err = inputLog.Fields.ReaderSingle("protoPayload.request")
	if err == nil {
		return revisionsFromRequestField(ctx, readerFactory, mergerConfig, cs, inputLog, builder, k8sOperation, authenticatedPrincipal, startTime)
	}

	tb := builder.GetTimelineBuilder(k8sOperation.CovertToResourcePath())
	lastRevision := tb.GetLatestRevision()
	cs.RecordRevision(k8sOperation.CovertToResourcePath(), &history.StagingResourceRevision{
		Verb:       k8sOperation.Verb,
		Requestor:  authenticatedPrincipal,
		Body:       "# Resource data is not available. This request audit log is recorded as Metadata level.",
		ChangeTime: inputLog.Timestamp(),
		State:      verbToRevisionState(k8sOperation.Verb, lastRevision),
	})
	return nil
}

func NewGKEAuditLogParser() parser.Parser {
	return &k8sAuditLogParser{}
}

func revisionsFromRequestField(ctx context.Context, readerFactory *structure.ReaderFactory, mergeConfigRegistry *model_k8s.MergeConfigRegistry, cs *history.ChangeSet, inputLog *log.LogEntity, builder *history.Builder, k8sOp *model.KubernetesObjectOperation, authenticatedPrincipal string, startTime time.Time) error {
	requestYamlReader, err := inputLog.Fields.ReaderSingle("protoPayload.request")
	if err != nil {
		return err
	}
	requestYaml, err := requestYamlReader.ToYaml("")
	if err != nil {
		return err
	}
	isPartialRequest := true
	if k8sOp.Verb != enum.RevisionVerbPatch || !strings.Contains(requestYaml, "'@type': k8s.io/Patch") {
		isPartialRequest = false
	}
	tb := builder.GetTimelineBuilder(k8sOp.CovertToResourcePath())
	requestYaml = removeAtType(requestYaml)
	err = generateRevisionsForPodBindings(cs, inputLog, builder, k8sOp, requestYaml)
	if err != nil {
		return err
	}
	if !isPartialRequest {
		lastRevision := tb.GetLatestRevision()
		handleCreationTime(ctx, startTime, requestYamlReader, builder, cs, inputLog, k8sOp)
		cs.RecordRevision(k8sOp.CovertToResourcePath(), &history.StagingResourceRevision{
			Verb:       k8sOp.Verb,
			Body:       requestYaml,
			Requestor:  authenticatedPrincipal,
			ChangeTime: inputLog.Timestamp(),
			State:      verbToRevisionState(k8sOp.Verb, lastRevision),
		})
		return handleResourceManifests(ctx, k8sOp, builder, cs, inputLog, requestYamlReader)
	}
	// for partial requests(Needs merge.)
	if lastRevision := tb.GetLatestRevision(); lastRevision != nil {
		prev, err := tb.GetLatestRevisionBody()
		if err != nil {
			return err
		}
		resolver := mergeConfigRegistry.Get(k8sOp.APIVersion, k8sOp.PluralKind)
		mergedReader, err := readerFactory.NewReader(adapter.MergeYaml(prev, requestYaml, resolver)) //structure.StrategicMergeYaml(prev, requestYaml, resolver)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("Failed to merge patch request\n%v", err))
			slog.DebugContext(ctx, fmt.Sprintf("Merge\nOLD:\n%s\n\n=====\n\nNEW:\n%s", prev, requestYaml))
			mergedReader, err = readerFactory.NewReader(adapter.Yaml(requestYaml))
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Failed to fall back to use the patch resource directly. Failed to read the yaml. \n%s", requestYaml))
				return err
			}
		}
		mergedYaml, err := mergedReader.ToYaml("")
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed to serialize the merged patch request.\n%s", err))
			return err
		}
		isDeleted, err := isDeletedResource(mergedReader)
		if err != nil {
			// Failed to get if is deleted or not
			isDeleted = false
		}
		revisionState := verbToRevisionState(k8sOp.Verb, lastRevision)
		if isDeleted {
			revisionState = enum.RevisionStateDeleted
		}
		handleCreationTime(ctx, startTime, mergedReader, builder, cs, inputLog, k8sOp)
		cs.RecordRevision(k8sOp.CovertToResourcePath(), &history.StagingResourceRevision{
			Verb:       k8sOp.Verb,
			Body:       mergedYaml,
			Partial:    lastRevision.Partial,
			Requestor:  authenticatedPrincipal,
			ChangeTime: inputLog.Timestamp(),
			State:      revisionState,
		})
		return handleResourceManifests(ctx, k8sOp, builder, cs, inputLog, mergedReader)
	} else {
		resolver := mergeConfigRegistry.Get(k8sOp.APIVersion, k8sOp.PluralKind)
		mergedReader, err := readerFactory.NewReader(adapter.MergeYaml("{}", requestYaml, resolver))
		if err != nil {
			return err
		}
		handleCreationTime(ctx, startTime, mergedReader, builder, cs, inputLog, k8sOp)
		isDeleted, err := isDeletedResource(mergedReader)
		if err != nil {
			isDeleted = false
		}
		revisionState := verbToRevisionState(k8sOp.Verb, nil)
		if isDeleted {
			revisionState = enum.RevisionStateDeleted
		}
		cs.RecordRevision(k8sOp.CovertToResourcePath(), &history.StagingResourceRevision{
			Verb:       k8sOp.Verb,
			Body:       requestYaml,
			Partial:    true,
			Requestor:  authenticatedPrincipal,
			ChangeTime: inputLog.Timestamp(),
			State:      revisionState,
		})
		return handleResourceManifests(ctx, k8sOp, builder, cs, inputLog, mergedReader)
	}
}

func revisionsFromResponseField(ctx context.Context, readerFactory *structure.ReaderFactory, cs *history.ChangeSet, inputLog *log.LogEntity, builder *history.Builder, k8sOp *model.KubernetesObjectOperation, authenticatedPrincipal string, startTime time.Time) error {
	responseYamlReader, err := inputLog.Fields.ReaderSingle("protoPayload.response")
	if err != nil {
		return err
	}
	responseYaml, err := responseYamlReader.ToYaml("")
	if err != nil {
		return err
	}
	responseYaml = removeAtType(responseYaml)
	err = generateRevisionsForPodBindings(cs, inputLog, builder, k8sOp, responseYaml)
	if err != nil {
		return err
	}
	tb := builder.GetTimelineBuilder(k8sOp.CovertToResourcePath())
	lastRevision := tb.GetLatestRevision()
	handleCreationTime(ctx, startTime, responseYamlReader, builder, cs, inputLog, k8sOp)
	isDeleted, err := isDeletedResource(responseYamlReader)
	if err != nil {
		isDeleted = false
	}
	revisionState := verbToRevisionState(k8sOp.Verb, lastRevision)
	if isDeleted {
		revisionState = enum.RevisionStateDeleted
	}
	revision := &history.StagingResourceRevision{
		Verb:       k8sOp.Verb,
		Body:       responseYaml,
		Requestor:  authenticatedPrincipal,
		ChangeTime: inputLog.Timestamp(),
		State:      revisionState,
	}
	manifest, err := readerFactory.NewReader(adapter.Yaml(responseYaml))
	if err != nil {
		return err
	}
	if err == nil {
		handleOwnerReference(cs, manifest, k8sOp, revision)
	}
	cs.RecordRevision(k8sOp.CovertToResourcePath(), revision)
	err = handleResourceManifests(ctx, k8sOp, builder, cs, inputLog, responseYamlReader)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("Failed to process the response yaml\n%s", err))
	}
	if k8sOp.Verb == enum.RevisionVerbDelete {
		handleDeletion(k8sOp, builder, cs, inputLog, authenticatedPrincipal)
	}
	return nil
}

// handleCreationTime will parse `metadata.creationTime` of the resource manifest and modify the verb of the request and create an inffered revision.
func handleCreationTime(ctx context.Context, queryStartTime time.Time, resourceYaml *structure.Reader, builder *history.Builder, cs *history.ChangeSet, log *log.LogEntity, k8sOp *model.KubernetesObjectOperation) {
	if k8sOp.Verb == enum.RevisionVerbCreate {
		return // this handler DOES NOT accept "create" operation
	}

	tb := builder.GetTimelineBuilder(k8sOp.CovertToResourcePath())
	if tb.GetLatestRevision() != nil {
		return // if there were last revision in the history, no need to check creation time
	}

	creationTime, err := getCreationTimeFromManifest(resourceYaml)
	if err != nil {
		slog.DebugContext(ctx, fmt.Sprintf("Failed to get the creation time from log %s\n%v", log.ID(), err))
		return
	}

	// Rewrite the verb to `create` because this is the first request during the log period and creationTime is after the start time
	// Resources can be created with `update` request but we can't see if the updated resource was existed from before or not.
	// Thus, if the resource creationTime is after start time of log query, the update should be regarded as a creation request.
	if creationTime.Sub(queryStartTime) > 0 {
		k8sOp.Verb = enum.RevisionVerbCreate
		return
	}

	// The creation request was not found in the log period. Infer the creation request from creation time.
	cs.RecordRevision(k8sOp.CovertToResourcePath(), &history.StagingResourceRevision{
		Verb:       enum.RevisionVerbCreate,
		Body:       `# The actual creation time is not included in the query time. This creation request was inferred from "metadata.creationTime" in the later requests.`,
		Requestor:  "unknown",
		ChangeTime: creationTime,
		State:      enum.RevisionStateInferred,
		Inferred:   true,
	})
}

// Handles other type of shadowed resources.
func generateRevisionsForPodBindings(cs *history.ChangeSet, inputLog *log.LogEntity, builder *history.Builder, k8sOp *model.KubernetesObjectOperation, content string) error {
	methodName, err := inputLog.GetString("protoPayload.methodName")
	if err != nil {
		return err
	}
	if strings.HasPrefix(methodName, "io.k8s.core.v1.pods.binding.") {
		target := inputLog.GetStringOrDefault("protoPayload.request.target.name", inputLog.GetStringOrDefault("protoPayload.response.target.name", "unknown"))
		shadowedPath := fmt.Sprintf("core/v1#node#cluster-scope#%s#%s-%s", target, k8sOp.Namespace, k8sOp.Name)
		cs.RecordAliasRelationship(k8sOp.CovertToResourcePath(), shadowedPath, history.RewriteRelationship(enum.RelationshipPodBinding))
		podK8sOp := model.KubernetesObjectOperation{
			APIVersion: k8sOp.APIVersion,
			PluralKind: k8sOp.PluralKind,
			Namespace:  k8sOp.Namespace,
			Name:       k8sOp.Name,
		}
		podScheduledStatusPath := fmt.Sprintf("%s#PodScheduled", podK8sOp.CovertToResourcePath())
		cs.RecordRevision(podScheduledStatusPath, &history.StagingResourceRevision{
			Verb:       enum.RevisionVerbStatusTrue,
			Body:       "# PodScheduled status was inferred to be `True` from a binding resource",
			Partial:    false,
			Requestor:  "",
			ChangeTime: inputLog.Timestamp(),
			State:      enum.RevisionStateConditionTrue,
		}, history.RewriteRelationship(enum.RelationshipResourceStatus))
	}
	return nil
}

// deleteResourcesUnderNamespace add the deletion records for each resources under the namespace.
// This function is expected to be called by a k8s api request of `deletecollection` verb.
func deleteResourcesUnderNamespace(cs *history.ChangeSet, builder *history.Builder, target *model.KubernetesObjectOperation, l *log.LogEntity, principal string) error {
	resources := builder.GetChildResources(fmt.Sprintf("%s#%s#%s", target.APIVersion, target.PluralKind, target.Namespace))
	for _, resource := range resources {
		childOperation := &model.KubernetesObjectOperation{
			APIVersion: target.APIVersion,
			PluralKind: target.PluralKind,
			Namespace:  target.Namespace,
			Name:       resource.ResourceName,
			Verb:       enum.RevisionVerbDelete,
		}
		tb := builder.GetTimelineBuilder(resource.FullResourcePath)
		latest := tb.GetLatestRevision()
		// The child resource is already deleted. The deletecollection won't make any effect on the deleted resource in the past.
		if latest != nil && latest.State == enum.RevisionStateDeleted {
			continue
		}
		cs.RecordRevision(childOperation.CovertToResourcePath(), &history.StagingResourceRevision{
			Verb:       enum.RevisionVerbDelete,
			Body:       "",
			Partial:    false,
			Requestor:  principal,
			ChangeTime: l.Timestamp(),
			State:      enum.RevisionStateDeleted,
		})
		err := handleDeletion(childOperation, builder, cs, l, principal)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleOwnerReference(cs *history.ChangeSet, manifestBody *structure.Reader, sourceK8sResource *model.KubernetesObjectOperation, sourceRevision *history.StagingResourceRevision) {
	if !manifestBody.Has("metadata.ownerReferences") {
		return
	}
	ownerReferencesReaders, err := manifestBody.Reader("metadata.ownerReferences[]")
	if err != nil {
		return
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
		namespace := sourceK8sResource.Namespace
		if kind == "Node" {
			namespace = "cluster-scope"
		}
		shadowedResource := model.KubernetesObjectOperation{
			APIVersion:      apiVersion,
			PluralKind:      strings.ToLower(kind),
			Name:            name,
			Namespace:       namespace,
			SubResourceName: fmt.Sprintf("%s-%s-%s", sourceK8sResource.PluralKind, sourceK8sResource.Namespace, sourceK8sResource.Name),
		}
		cs.RecordAliasRelationship(sourceK8sResource.CovertToResourcePath(), shadowedResource.CovertToResourcePath(), history.RewriteRelationship(enum.RelationshipOwnerReference))
	}
}

func parseStatusConditions(k8sOp *model.KubernetesObjectOperation, builder *history.Builder, resourceYamlReader *structure.Reader, cs *history.ChangeSet, l *log.LogEntity) error {
	var resourceContainingStatus model.K8sResourceContainingStatus
	err := resourceYamlReader.ReadReflect("", &resourceContainingStatus)
	if err != nil {
		return err
	}
	if resourceContainingStatus.Status == nil || len(resourceContainingStatus.Status.Conditions) == 0 {
		// This resource has no status field or no conditions in status field
		return nil
	}
	for _, condition := range resourceContainingStatus.Status.Conditions {
		lastTransitionTime, err := time.Parse(time.RFC3339, condition.LastTransitionTime)
		if err != nil {
			continue
		}
		conditionTime := lastTransitionTime
		lastHeartbeatTime, err := time.Parse(time.RFC3339, condition.LastHeartbeatTime)
		if err == nil && lastHeartbeatTime.Sub(conditionTime) > 0 {
			conditionTime = lastHeartbeatTime
		}
		lastProbeTime, err := time.Parse(time.RFC3339, condition.LastProbeTime)
		if err == nil && lastProbeTime.Sub(lastHeartbeatTime) > 0 {
			conditionTime = lastProbeTime
		}
		// Ignore if the transition time was older than the last revision

		statusPath := fmt.Sprintf("%s#%s", k8sOp.CovertToResourcePath(), condition.Type)
		if k8sOp.SubResourceName != "" {
			parentOp := model.KubernetesObjectOperation{
				APIVersion: k8sOp.APIVersion,
				PluralKind: k8sOp.PluralKind,
				Namespace:  k8sOp.Namespace,
				Name:       k8sOp.Name,
				Verb:       k8sOp.Verb,
			}
			statusPath = fmt.Sprintf("%s#%s", parentOp.CovertToResourcePath(), condition.Type)
		}
		tb := builder.GetTimelineBuilder(statusPath)
		latest := tb.GetLatestRevision()
		latestTime := time.Time{}
		if latest != nil {
			latestTime = latest.ChangeTime
		} else {
			creationTime, err := getCreationTimeFromManifest(resourceYamlReader)
			if err == nil && conditionTime.Sub(creationTime) != 0 {
				cs.RecordRevision(statusPath, &history.StagingResourceRevision{
					Verb:       enum.RevisionVerbStatusUnknown,
					Body:       "# status is unknown but existence is inferred from the later log.",
					Partial:    false,
					Inferred:   true,
					Requestor:  "",
					ChangeTime: creationTime,
					State:      enum.RevisionStateConditionUnknown,
				})
			}
		}
		if conditionTime.Sub(latestTime) > 0 {
			conditionYaml, err := goyaml.Marshal(condition)
			if err != nil {
				continue
			}
			cs.RecordRevision(statusPath, &history.StagingResourceRevision{
				Verb:       conditionStateToRevisionVerb(condition.Status),
				Body:       string(conditionYaml),
				Partial:    false,
				Requestor:  "",
				ChangeTime: conditionTime,
				State:      conditionStateToRevisionState(condition.Status),
			}, history.RewriteRelationship(enum.RelationshipResourceStatus))
		}
	}

	return nil
}

func parsePodContainerStatuses(k8sOp *model.KubernetesObjectOperation, builder *history.Builder, resourceYamlReader *structure.Reader, cs *history.ChangeSet, l *log.LogEntity) error {
	const errorTimestampInUnix = -62135596800 // Unix time for 0001-01-01T00:00:00Z
	var pod corev1.Pod
	err := resourceYamlReader.ReadReflectK8sManifest("", &pod)
	if err != nil {
		return err
	}
	// Audit logs is not assured to be generated when a container becomes ready. And the ready field has no timestamp of the change
	// containers ready last transition time is used for the time of becoming ready in the containers.
	containersReadyTime := time.Unix(1<<63-1, 0)
	for _, podConditions := range pod.Status.Conditions {
		if podConditions.Type == "ContainersReady" {
			containersReadyTime = podConditions.LastTransitionTime.Time
			break
		}
	}
	statuses := []corev1.ContainerStatus{}
	statuses = append(statuses, pod.Status.ContainerStatuses...)
	statuses = append(statuses, pod.Status.InitContainerStatuses...)
	for i, status := range statuses {
		statusYaml, err := yaml.Marshal(status)
		if err != nil {
			return err
		}
		isInitContainer := i >= len(pod.Status.ContainerStatuses)
		cpath := resourcepath.Container(k8sOp.Namespace, k8sOp.Name, status.Name)
		changed := builder.ClusterResource.ContainerStatuses.IsNewChange(k8sOp.Namespace, k8sOp.Name, status.Name, status)
		tb := builder.GetTimelineBuilder(cpath)
		last := tb.GetLatestRevision()
		if changed {
			// Current container is running
			if status.State.Running != nil {
				running := status.State.Running
				time := running.StartedAt.Time
				if last != nil && time.Sub(last.ChangeTime) > 0 && l.Timestamp().Sub(time) > 0 && status.Ready {
					cs.RecordRevision(cpath, &history.StagingResourceRevision{
						Verb:       enum.RevisionVerbContainerNonReady,
						Body:       string(statusYaml),
						Requestor:  "",
						Partial:    false,
						ChangeTime: running.StartedAt.Time,
						State:      enum.RevisionStateContainerRunningNonReady,
					}, history.RewriteRelationship(enum.RelationshipContainer))
				}
				if status.Ready {
					readinessChangeTime := l.Timestamp()
					if !isInitContainer && last != nil && containersReadyTime.Sub(last.ChangeTime) > 0 {
						readinessChangeTime = containersReadyTime
					}
					cs.RecordRevision(cpath, &history.StagingResourceRevision{
						Verb:       enum.RevisionVerbContainerReady,
						Body:       string(statusYaml),
						Requestor:  "",
						Partial:    false,
						ChangeTime: readinessChangeTime,
						State:      enum.RevisionStateContainerRunningReady,
					}, history.RewriteRelationship(enum.RelationshipContainer))
				} else {
					cs.RecordRevision(cpath, &history.StagingResourceRevision{
						Verb:       enum.RevisionVerbContainerNonReady,
						Body:       string(statusYaml),
						Requestor:  "",
						Partial:    false,
						ChangeTime: l.Timestamp(),
						State:      enum.RevisionStateContainerRunningNonReady,
					}, history.RewriteRelationship(enum.RelationshipContainer))
				}
			} else if status.State.Terminated != nil { // Current container is terminated
				terminated := status.State.Terminated
				if terminated.FinishedAt.Time.Unix() == errorTimestampInUnix {
					// Pod termination status can have errornous timestamp when it can't be determined.
					// We still don't know the exact time that happening but it should be in between last change time and current log time.
					// Use timestamp log in the case.
					terminated.FinishedAt.Time = l.Timestamp()
				}
				verb := enum.RevisionVerbContainerSuccess
				state := enum.RevisionStateContainerTerminatedWithSuccess
				if terminated.ExitCode != 0 {
					verb = enum.RevisionVerbContainerError
					state = enum.RevisionStateContainerTerminatedWithError
				}
				cs.RecordRevision(cpath, &history.StagingResourceRevision{
					Verb:       verb,
					Body:       string(statusYaml),
					Requestor:  "",
					Partial:    false,
					ChangeTime: terminated.FinishedAt.Time,
					State:      state,
				}, history.RewriteRelationship(enum.RelationshipContainer))
			} else if status.State.Waiting != nil { // Current container is waiting
				cs.RecordRevision(cpath, &history.StagingResourceRevision{
					Verb:       enum.RevisionVerbContainerWaiting,
					Body:       string(statusYaml),
					Requestor:  "",
					Partial:    false,
					ChangeTime: l.Timestamp(),
					State:      enum.RevisionStateContainerWaiting,
				}, history.RewriteRelationship(enum.RelationshipContainer))
			}
		}
	}
	return nil
}

// handleResourceManifests receives the current manifest and call subroutines regarding the resource types
func handleResourceManifests(ctx context.Context, k8sOp *model.KubernetesObjectOperation, builder *history.Builder, cs *history.ChangeSet, l *log.LogEntity, manifestReader *structure.Reader) error {
	err := parseStatusConditions(k8sOp, builder, manifestReader, cs, l)
	if err != nil {
		return err
	}
	if k8sOp.APIVersion == "core/v1" && k8sOp.PluralKind == "pod" {
		err := parsePodContainerStatuses(k8sOp, builder, manifestReader, cs, l)
		if err != nil {
			return err
		}
		var podResource model.Pod
		err = manifestReader.ReadReflect("", &podResource)
		if err != nil {
			return err
		}
		if podResource.Status != nil && podResource.Status.PodIPs != nil {
			for _, podIP := range podResource.Status.PodIPs {
				if podResource.Status.HostIP == podIP.IP {
					// Ignore .spec.hostNetwork = true not to overwrite node's IP lease
					continue
				}
				builder.ClusterResource.IPs.TouchResourceLease(podIP.IP, l.Timestamp(), resourcelease.NewK8sResourceLeaseHolder(k8sOp.PluralKind, k8sOp.Namespace, k8sOp.Name))
			}
		}
	}
	// Add endpointslices status on pod subresource
	if k8sOp.APIVersion == "discovery.k8s.io/v1" && k8sOp.PluralKind == "endpointslice" {
		uid, err := builder.ClusterResource.EndpointSlices.ReadEndpointSlice(manifestReader, l.Timestamp())
		if err != nil {
			return err
		}
		diffs, err := builder.ClusterResource.EndpointSlices.GetLastDiffs(uid)
		if err != nil {
			return err
		}
		for _, diff := range diffs {
			target := diff.TargetRef
			if target == nil {
				// targetRef can be nil when this Endpoint is associated to a ExternalService resource
				continue
			}
			if target.Kind != "Pod" {
				slog.WarnContext(ctx, fmt.Sprintf("unsupported target ref of endpoint slice %s", target.Kind))
				continue
			}
			state := enum.RevisionStateConditionFalse
			verb := enum.RevisionVerbNonReady
			if diff.Current != nil {
				if diff.Current.Conditions.Ready {
					state = enum.RevisionStateConditionTrue
					verb = enum.RevisionVerbReady
				}
			}
			resourcePath := fmt.Sprintf("core/v1#pod#%s#%s#endpointslces-%s-%s", target.Namespace, target.Name, k8sOp.Namespace, k8sOp.Name)
			manifestYaml, err := manifestReader.ToYaml("")
			if err != nil {
				return err
			}
			cs.RecordRevision(resourcePath, &history.StagingResourceRevision{
				Body:       manifestYaml,
				State:      state,
				Verb:       verb,
				ChangeTime: l.Timestamp(),
			}, history.RewriteRelationship(enum.RelationshipEndpointSlice))
		}
	}

	if k8sOp.APIVersion == "networking.gke.io/v1beta1" && k8sOp.PluralKind == "servicenetworkendpointgroup" {
		builder.ClusterResource.NEGs.TouchResourceLease(k8sOp.Name, l.Timestamp(), resourcelease.NewK8sResourceLeaseHolder(k8sOp.PluralKind, k8sOp.Namespace, k8sOp.Name))
	}

	return nil
}

func handleDeletion(k8sOp *model.KubernetesObjectOperation, builder *history.Builder, cs *history.ChangeSet, l *log.LogEntity, principal string) error {
	if k8sOp.APIVersion == "core/v1" && k8sOp.PluralKind == "pod" {
		lastRevision := builder.GetTimelineBuilder(k8sOp.CovertToResourcePath()).GetLatestRevision()
		if lastRevision != nil && lastRevision.State != enum.RevisionStateDeleted {
			podBindingOp := *k8sOp
			podBindingOp.SubResourceName = "binding"
			cs.RecordRevision(podBindingOp.CovertToResourcePath(), &history.StagingResourceRevision{
				Verb:       enum.RevisionVerbDelete,
				Body:       "# binding resource was inferred to be deleted from the pod deletion request",
				Partial:    false,
				Requestor:  principal,
				ChangeTime: l.Timestamp(),
				State:      enum.RevisionStateDeleted,
			})
		}
	}
	return nil
}

func handleNodeResource(k8sOp *model.KubernetesObjectOperation, builder *history.Builder) {
	builder.ClusterResource.AddNode(k8sOp.Name)
}

// isDeletedResource reads `labels.deletionTimestamp` field in the manifest and return if the resource is already deleted or not.
// patch request itself can't say if the resource is already deleted or not, it needs to read the field to determine the resource status.
func isDeletedResource(reader *structure.Reader) (bool, error) {
	type resourceManifestMetadata struct {
		DeletionTimestamp *time.Time `json:"deletionTimestamp"`
	}
	type resourceManifest struct {
		Metadata *resourceManifestMetadata `json:"metadata"`
	}
	var parsedManifest resourceManifest
	err := reader.ReadReflect("", &parsedManifest)
	if err != nil {
		return false, err
	}
	if parsedManifest.Metadata == nil {
		return false, fmt.Errorf("failed to get metadata section in the manifest")
	}
	return parsedManifest.Metadata.DeletionTimestamp != nil, nil
}

func getCreationTimeFromManifest(manifest *structure.Reader) (time.Time, error) {
	type resourceManifestMetadata struct {
		CreationTimestamp *time.Time `json:"creationTimestamp"`
	}
	type resourceManifest struct {
		Metadata *resourceManifestMetadata `json:"metadata"`
	}
	var parsedManifest resourceManifest
	err := manifest.ReadReflect("", &parsedManifest)
	if err != nil {
		return time.Time{}, err
	}
	if parsedManifest.Metadata == nil {
		return time.Time{}, fmt.Errorf("failed to get metadata field in the manifest")
	}
	if parsedManifest.Metadata.CreationTimestamp == nil {
		return time.Time{}, fmt.Errorf("failed to get metadata.creationTimestamp field in the manifest")
	}
	return *parsedManifest.Metadata.CreationTimestamp, nil
}

func verbToRevisionState(verb enum.RevisionVerb, lastRevision *history.ResourceRevision) enum.RevisionState {
	if verb == enum.RevisionVerbDelete {
		return enum.RevisionStateDeleted
	} else if verb == enum.RevisionVerbPatch {
		// patch operation can happen on a deleted resource. Use the previous state for the patch request.
		if lastRevision == nil {
			return enum.RevisionStateExisting
		} else {
			return lastRevision.State
		}
	} else {
		return enum.RevisionStateExisting
	}
}

func conditionStateToRevisionVerb(conditionState string) enum.RevisionVerb {
	if conditionState == "True" {
		return enum.RevisionVerbStatusTrue
	} else if conditionState == "False" {
		return enum.RevisionVerbStatusFalse
	}
	return enum.RevisionVerbStatusUnknown
}

func conditionStateToRevisionState(conditionState string) enum.RevisionState {
	if conditionState == "True" {
		return enum.RevisionStateConditionTrue
	} else if conditionState == "False" {
		return enum.RevisionStateConditionFalse
	}
	return enum.RevisionStateConditionUnknown
}

// Remove @type in response or request payload
func removeAtType(yamlString string) string {
	if strings.Contains(yamlString, "'@type'") {
		index := strings.Index(yamlString, "\n")
		return yamlString[index+1:]
	}
	return yamlString
}
