// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package privategkemaster_contract

import (
	"context"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/logutil"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudlogk8scontrolplane_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8scontrolplane/contract"
	googlecloudlogk8snode_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8snode/contract"
)

var (
	pathProjectID          = structured.CompileFieldPath("resource.labels.project_id")
	pathLabels             = structured.CompileFieldPath("labels")
	pathResourceType       = structured.CompileFieldPath("resource.type")
	pathLogName            = structured.CompileFieldPath("logName")
	pathPodID              = structured.CompileFieldPath("resource.labels.pod_id")
	pathNamespaceID        = structured.CompileFieldPath("resource.labels.namespace_id")
	pathContainerName      = structured.CompileFieldPath("resource.labels.container_name")
	pathTextPayload        = structured.CompileFieldPath("textPayload")
	pathJSONPayloadMessage = structured.CompileFieldPath("jsonPayload.message")
	pathJSONPayloadMESSAGE = structured.CompileFieldPath("jsonPayload.MESSAGE")
)

var (
	defaultGKEMasterParsersMap = map[string]logutil.StructuredLogParser{
		"kubelet": logutil.NewMultiTextLogParser(
			logutil.NewKLogTextParser(true),
			&logutil.FallbackRawTextLogParser{},
		),
		"kube-controller-manager": logutil.NewMultiTextLogParser(
			logutil.NewKLogTextParser(false),
			&logutil.FallbackRawTextLogParser{},
		),
		"containerd": logutil.NewMultiTextLogParser(
			logutil.NewLogfmtTextParser(),
			&logutil.FallbackRawTextLogParser{},
		),
	}
	defaultGKEMasterParser = logutil.NewMultiTextLogParser(
		logutil.NewJsonlTextParser(),
		logutil.NewKLogTextParser(false),
		&logutil.FallbackRawTextLogParser{},
	)
)

// PrivateGKEMasterParserType is the type of the private GKE master parser.
type PrivateGKEMasterParserType string

var (
	// PrivateGKEMasterParserTypeScheduler is the parser type for scheduler logs.
	PrivateGKEMasterParserTypeScheduler PrivateGKEMasterParserType = "scheduler"
	// PrivateGKEMasterParserTypeControllerManager is the parser type for controller manager logs.
	PrivateGKEMasterParserTypeControllerManager PrivateGKEMasterParserType = "controller-manager"
	// PrivateGKEMasterParserTypeKubelet is the parser type for kubelet logs.
	PrivateGKEMasterParserTypeKubelet PrivateGKEMasterParserType = "kubelet"
	// PrivateGKEMasterParserTypeContainerRuntime is the parser type for container runtime (containerd) logs.
	PrivateGKEMasterParserTypeContainerRuntime PrivateGKEMasterParserType = "container-runtime"
	// PrivateGKEMasterParserTypeOther is the parser type for other logs.
	PrivateGKEMasterParserTypeOther PrivateGKEMasterParserType = "other"
)

var componentNameToPrivateGKEMasterParserTypeMap = map[string]PrivateGKEMasterParserType{
	"kube-scheduler":          PrivateGKEMasterParserTypeScheduler,
	"kube-controller-manager": PrivateGKEMasterParserTypeControllerManager,
	"kubelet":                 PrivateGKEMasterParserTypeKubelet,
	"containerd":              PrivateGKEMasterParserTypeContainerRuntime,
}

// GKEMasterLogFieldSet is the field set for GKE Master logs.
type GKEMasterLogFieldSet struct {
	ProjectID      string
	HostName       string
	NamespaceID    string
	ComponentName  string
	PodID          string
	ContainerName  string
	StructuredBody *logutil.ParseStructuredLogResult
}

// ResourceTimelines returns the timeline paths for the GKE Master log.
func (g *GKEMasterLogFieldSet) ResourceTimelines(ctx context.Context, clusterName string) []*khifilev6.TimelinePath {
	projectTimeline := googlecloudcommon_contract.MustGCPProjectTimeline(ctx, g.ProjectID)
	gkeTimeline := googlecloudcommon_contract.MustGKEClusterTimeline(ctx, projectTimeline, clusterName)
	compTimeline := googlecloudlogk8scontrolplane_contract.MustControlPlaneComponentTimeline(ctx, gkeTimeline, g.ComponentName)

	if g.NamespaceID == "" {
		clusterTimeline := commonlogk8saudit_contract.MustK8sClusterTimeline(ctx, clusterName)
		apiVersionTimeline := commonlogk8saudit_contract.MustK8sAPIVersionTimeline(ctx, clusterTimeline, "core/v1")
		kindTimeline := commonlogk8saudit_contract.MustK8sKindTimeline(ctx, apiVersionTimeline, "node")
		nodeTimeline := commonlogk8saudit_contract.MustK8sClusterScopeResourceTimeline(ctx, kindTimeline, g.HostName)
		nodeCompTimeline := googlecloudlogk8snode_contract.MustNodeComponentTimeline(ctx, nodeTimeline, g.ComponentName)
		return []*khifilev6.TimelinePath{compTimeline, nodeCompTimeline}
	} else {
		clusterTimeline := commonlogk8saudit_contract.MustK8sClusterTimeline(ctx, clusterName)
		apiVersionTimeline := commonlogk8saudit_contract.MustK8sAPIVersionTimeline(ctx, clusterTimeline, "core/v1")
		kindTimeline := commonlogk8saudit_contract.MustK8sKindTimeline(ctx, apiVersionTimeline, "pod")
		nsTimeline := commonlogk8saudit_contract.MustK8sNamespaceTimeline(ctx, kindTimeline, g.NamespaceID)
		podTimeline := commonlogk8saudit_contract.MustK8sNamespacedResourceTimeline(ctx, nsTimeline, g.PodID)
		containerTimeline := commonlogk8saudit_contract.MustK8sContainerTimeline(ctx, podTimeline, g.ContainerName)
		return []*khifilev6.TimelinePath{compTimeline, containerTimeline}
	}
}

// PrivateGKEMasterParserType returns the parser type for the GKE Master log.
func (g *GKEMasterLogFieldSet) PrivateGKEMasterParserType() PrivateGKEMasterParserType {
	if parserType, found := componentNameToPrivateGKEMasterParserTypeMap[g.ComponentName]; found {
		return parserType
	}
	return PrivateGKEMasterParserTypeOther
}

// ExtractGKEMasterLog extracts GKEMasterLogFieldSet from a NodeReader.
func ExtractGKEMasterLog(reader *structured.NodeReader) (GKEMasterLogFieldSet, error) {
	if mock, ok := structured.GetMock[GKEMasterLogFieldSet](reader); ok {
		return mock, nil
	}
	var result GKEMasterLogFieldSet
	result.ProjectID = reader.ReadStringOrDefault(pathProjectID, "unknown")
	labelsReader, err := reader.GetReader(pathLabels)
	nodeName := "unknown-master-node"
	if err == nil {
		for c, v := range labelsReader.Children() {
			if c.Key == "compute.googleapis.com/resource_name" {
				nodeName = v.ReadStringOrDefault(structured.EmptyFieldPath, "")
				break
			}
		}
	}
	result.HostName = nodeName
	resourceType := reader.ReadStringOrDefault(pathResourceType, "")
	logName := reader.ReadStringOrDefault(pathLogName, "")
	componentNameBeginIndex := strings.LastIndex(logName, "/")
	if componentNameBeginIndex != -1 {
		result.ComponentName = logName[componentNameBeginIndex+1:]
	}
	if resourceType == "container" {
		result.PodID = reader.ReadStringOrDefault(pathPodID, "")
		result.NamespaceID = reader.ReadStringOrDefault(pathNamespaceID, "")
		result.ContainerName = reader.ReadStringOrDefault(pathContainerName, "")
	}

	// Parse structured log body
	messageSource := reader.ReadStringOrDefault(pathTextPayload, "")
	if messageSource == "" {
		messageSource = reader.ReadStringOrDefault(pathJSONPayloadMessage, "")
	}
	if messageSource == "" {
		messageSource = reader.ReadStringOrDefault(pathJSONPayloadMESSAGE, "")
	}
	parser, ok := defaultGKEMasterParsersMap[result.ComponentName]
	if !ok {
		parser = defaultGKEMasterParser
	}
	result.StructuredBody = parser.TryParse(messageSource)

	return result, nil
}

// ExtractGKEMasterCommonMessage extracts the common message string from a NodeReader.
func ExtractGKEMasterCommonMessage(reader *structured.NodeReader) (string, error) {
	if mock, ok := structured.GetMock[googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet](reader); ok {
		return mock.Message, nil
	}
	message := reader.ReadStringOrDefault(pathTextPayload, "")
	if message == "" {
		message = reader.ReadStringOrDefault(pathJSONPayloadMessage, "")
	}
	if message == "" {
		message = reader.ReadStringOrDefault(pathJSONPayloadMESSAGE, "")
	}
	return message, nil
}
