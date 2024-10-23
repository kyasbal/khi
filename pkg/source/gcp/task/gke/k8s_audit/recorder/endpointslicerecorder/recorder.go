package endpointslicerecorder

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourceinfo/resourcelease"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/recorder"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/types"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"

	goyaml "gopkg.in/yaml.v3"
)

func Register(manager *recorder.RecorderTaskManager) error {
	manager.AddRecorder("endpointslices", []string{}, func(ctx context.Context, resourcePath string, currentLog *types.ResourceSpecificParserInput, prevStateInGroup any, cs *history.ChangeSet, builder *history.Builder, vs *task.VariableSet) (any, error) {
		var prevEndpointSlice *model.EndpointSlice
		if prevStateInGroup != nil {
			prevEndpointSlice = prevStateInGroup.(*model.EndpointSlice)
		}
		return recordChangeSetForLog(ctx, resourcePath, currentLog, prevEndpointSlice, cs, builder)
	}, recorder.ResourceKindLogGroupFilter("endpointslice"), recorder.AndLogFilter(recorder.OnlySucceedLogs(), recorder.OnlyWithResourceBody()))
	return nil
}

func recordChangeSetForLog(ctx context.Context, resourcePath string, log *types.ResourceSpecificParserInput, prevEndpointSlices *model.EndpointSlice, cs *history.ChangeSet, builder *history.Builder) (*model.EndpointSlice, error) {
	var endpointSlice model.EndpointSlice
	err := log.ResourceBodyReader.ReadReflect("", &endpointSlice)
	if err != nil {
		return nil, err
	}
	relatedServiceName := ""
	if endpointSlice.Metadata != nil && endpointSlice.Metadata.OwnerReferences != nil {
		for _, owner := range endpointSlice.Metadata.OwnerReferences {
			if strings.ToLower(owner.Kind) == "service" {
				if relatedServiceName != "" {
					slog.WarnContext(ctx, fmt.Sprintf("multiple owners found for a single endpoint slice. ignoreing service %s", relatedServiceName))
				}
				relatedServiceName = owner.Name
			}
		}
	}
	shouldRecordForService := relatedServiceName != ""

	if endpointSlice.Endpoints == nil {
		endpointSlice.Endpoints = make([]*model.EndpointSliceEndpoint, 0)
	}

	for _, endpoint := range endpointSlice.Endpoints {
		var prev *model.EndpointSliceEndpoint
		if endpoint.TargetRef != nil {
			prev = lookupEndpointFromUid(prevEndpointSlices, endpoint.TargetRef.Uid)
		}
		shouldRecordForPod := endpoint.TargetRef != nil && endpoint.TargetRef.Kind == "Pod"
		shouldRecordCondition := endpoint.Conditions != nil && (prev == nil || !endpoint.Conditions.SameWith(prev.Conditions))
		// records Ips used in Pods. IP can be read from Pod manifest, but it can be ignored when users didn't turn on DATA_WRITE audit log, but endpoint slice update will be recorded always.
		if shouldRecordForPod {
			for _, address := range endpoint.Addresses {
				builder.ClusterResource.IPs.TouchResourceLease(address, log.Log.Timestamp(), resourcelease.NewK8sResourceLeaseHolder(endpoint.TargetRef.Kind, endpoint.TargetRef.Namespace, endpoint.TargetRef.Name))
			}
		}
		// record conditions as subresource of pod.
		if shouldRecordCondition {
			state := enum.RevisionStateEndpointUnready
			verb := enum.RevisionVerbNonReady
			if endpoint.Conditions.Ready {
				state = enum.RevisionStateEndpointReady
				verb = enum.RevisionVerbReady
			} else if endpoint.Conditions.Terminating {
				state = enum.RevisionStateEndpointTerminating
				verb = enum.RevisionVerbTerminating
			}
			endpointYaml, err := goyaml.Marshal(endpoint)
			if err != nil {
				slog.WarnContext(ctx, fmt.Sprintf("failed to marshal endpoint data to yaml\n%s", err.Error()))
				continue
			}
			if shouldRecordForPod {
				podEndpointSliceResourcePath := resourcepath.PodEndpointSlice(log.Operation.Namespace, log.Operation.Name, endpoint.TargetRef.Namespace, endpoint.TargetRef.Name)
				cs.RecordRevision(podEndpointSliceResourcePath, &history.StagingResourceRevision{
					Body:       string(endpointYaml),
					State:      state,
					Verb:       verb,
					ChangeTime: log.Log.Timestamp(),
				})
			}
			if shouldRecordForService {
				serviceEndpointSliceResourcePath := resourcepath.ServiceEndpointSlice(log.Operation.Namespace, log.Operation.Name, relatedServiceName)
				cs.RecordRevision(serviceEndpointSliceResourcePath, &history.StagingResourceRevision{
					Body:       string(endpointYaml),
					State:      state,
					Verb:       verb,
					ChangeTime: log.Log.Timestamp(),
				})

			}
		}
	}
	// find endpoints included in the previous revision but not in in the current revision
	if prevEndpointSlices != nil {
		for _, endpoint := range prevEndpointSlices.Endpoints {
			var current *model.EndpointSliceEndpoint
			if endpoint.TargetRef != nil {
				current = lookupEndpointFromUid(&endpointSlice, endpoint.TargetRef.Uid)
			}
			if current != nil {
				continue
			}
			if endpoint.TargetRef != nil {
				podEndpointSliceResourcePath := resourcepath.PodEndpointSlice(log.Operation.Namespace, log.Operation.Name, endpoint.TargetRef.Namespace, endpoint.TargetRef.Name)
				// Only process endpoints not included in current endpoint slices
				cs.RecordRevision(podEndpointSliceResourcePath, &history.StagingResourceRevision{
					Body:       "# This endpoint removed from endpoint list of the EndpointSlice",
					State:      enum.RevisionStateDeleted,
					Verb:       enum.RevisionVerbDelete,
					ChangeTime: log.Log.Timestamp(),
				})
			}

			if shouldRecordForService {
				serviceEndpointSliceResourcePath := resourcepath.ServiceEndpointSlice(log.Operation.Namespace, log.Operation.Name, relatedServiceName)
				cs.RecordRevision(serviceEndpointSliceResourcePath, &history.StagingResourceRevision{
					Body:       "# This endpoint removed from endpoint list of the EndpointSlice",
					State:      enum.RevisionStateDeleted,
					Verb:       enum.RevisionVerbDelete,
					ChangeTime: log.Log.Timestamp(),
				})
			}
		}
	}
	return &endpointSlice, nil
}

func lookupEndpointFromUid(endpointSlices *model.EndpointSlice, uid string) *model.EndpointSliceEndpoint {
	if endpointSlices == nil || endpointSlices.Endpoints == nil {
		return nil
	}
	for _, endpoint := range endpointSlices.Endpoints {
		if endpoint.TargetRef != nil && endpoint.TargetRef.Uid == uid {
			return endpoint
		}
	}
	return nil
}
