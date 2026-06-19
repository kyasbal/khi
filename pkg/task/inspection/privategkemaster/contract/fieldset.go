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
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	commonlogk8saudit_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/commonlogk8saudit/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudlogk8scontrolplane_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8scontrolplane/contract"
	googlecloudlogk8snode_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8snode/contract"
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

// Kind implements [log.FieldSet].
func (g *GKEMasterLogFieldSet) Kind() string {
	return "gke-master"
}

var _ log.FieldSet = (*GKEMasterLogFieldSet)(nil)

// GKEMasterLogFieldSetReader reads the GKE Master log field set.
type GKEMasterLogFieldSetReader struct {
	parsersMap    map[string]logutil.StructuredLogParser
	defaultParser logutil.StructuredLogParser
}

// NewGKEMasterLogFieldSetReader creates a new GKEMasterLogFieldSetReader.
func NewGKEMasterLogFieldSetReader() *GKEMasterLogFieldSetReader {
	return &GKEMasterLogFieldSetReader{
		parsersMap: map[string]logutil.StructuredLogParser{
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
		},
		defaultParser: logutil.NewMultiTextLogParser(
			logutil.NewJsonlTextParser(),
			logutil.NewKLogTextParser(false),
			&logutil.FallbackRawTextLogParser{},
		),
	}
}

// FieldSetKind implements [log.FieldSetReader].
func (g *GKEMasterLogFieldSetReader) FieldSetKind() string {
	return (&GKEMasterLogFieldSet{}).Kind()
}

// Read implements [log.FieldSetReader].
func (g *GKEMasterLogFieldSetReader) Read(reader *structured.NodeReader) (log.FieldSet, error) {
	result := &GKEMasterLogFieldSet{}
	result.ProjectID = reader.ReadStringOrDefault("resource.labels.project_id", "unknown")
	labelsReader, err := reader.GetReader("labels")
	nodeName := "unknown-master-node"
	if err == nil {
		for c, v := range labelsReader.Children() {
			if c.Key == "compute.googleapis.com/resource_name" {
				nodeName = v.ReadStringOrDefault("", "")
				break
			}
		}
	}
	result.HostName = nodeName
	resourceType := reader.ReadStringOrDefault("resource.type", "")
	logName := reader.ReadStringOrDefault("logName", "")
	componentNameBeginIndex := strings.LastIndex(logName, "/")
	if componentNameBeginIndex != -1 {
		result.ComponentName = logName[componentNameBeginIndex+1:]
	}
	if resourceType == "container" {
		result.PodID = reader.ReadStringOrDefault("resource.labels.pod_id", "")
		result.NamespaceID = reader.ReadStringOrDefault("resource.labels.namespace_id", "")
		result.ContainerName = reader.ReadStringOrDefault("resource.labels.container_name", "")
	}

	// Parse structured log body
	messageSource := reader.ReadStringOrDefault("textPayload", "")
	if messageSource == "" {
		messageSource = reader.ReadStringOrDefault("jsonPayload.message", "")
	}
	if messageSource == "" {
		messageSource = reader.ReadStringOrDefault("jsonPayload.MESSAGE", "")
	}
	parser, ok := g.parsersMap[result.ComponentName]
	if !ok {
		parser = g.defaultParser
	}
	result.StructuredBody = parser.TryParse(messageSource)

	return result, nil
}

var _ log.FieldSetReader = (*GKEMasterLogFieldSetReader)(nil)

// GKEMasterCommonFieldSetReader reads the common field set for GKE Master logs.
type GKEMasterCommonFieldSetReader struct {
}

// FieldSetKind implements [log.FieldSetReader].
func (g *GKEMasterCommonFieldSetReader) FieldSetKind() string {
	return (&googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{}).Kind()
}

// Read implements [log.FieldSetReader].
func (g *GKEMasterCommonFieldSetReader) Read(reader *structured.NodeReader) (log.FieldSet, error) {
	message := reader.ReadStringOrDefault("textPayload", "")
	if message == "" {
		message = reader.ReadStringOrDefault("jsonPayload.message", "")
	}
	if message == "" {
		message = reader.ReadStringOrDefault("jsonPayload.MESSAGE", "")
	}
	return &googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{
		Message: message,
	}, nil
}

var _ log.FieldSetReader = (*GKEMasterCommonFieldSetReader)(nil)
