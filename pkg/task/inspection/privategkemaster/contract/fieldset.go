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
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/logutil"
	"github.com/GoogleCloudPlatform/khi/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudlogk8scontrolplane_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8scontrolplane/contract"
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
	HostName       string
	NamespaceID    string
	ComponentName  string
	PodID          string
	ContainerName  string
	StructuredBody *logutil.ParseStructuredLogResult
}

// ResourcePaths returns the resource paths for the GKE Master log.
func (g *GKEMasterLogFieldSet) ResourcePaths(clusterName string) []resourcepath.ResourcePath {
	if g.NamespaceID == "" {
		return []resourcepath.ResourcePath{
			resourcepath.ControlplaneComponent(clusterName, g.ComponentName),
			resourcepath.NodeComponent(g.HostName, g.ComponentName),
		}
	} else {
		return []resourcepath.ResourcePath{
			resourcepath.ControlplaneComponent(clusterName, g.ComponentName),
			resourcepath.Container(g.NamespaceID, g.PodID, g.ContainerName),
		}
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
}

// FieldSetKind implements [log.FieldSetReader].
func (g *GKEMasterLogFieldSetReader) FieldSetKind() string {
	return (&GKEMasterLogFieldSet{}).Kind()
}

// Read implements [log.FieldSetReader].
func (g *GKEMasterLogFieldSetReader) Read(reader *structured.NodeReader) (log.FieldSet, error) {
	result := &GKEMasterLogFieldSet{}
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
	parser := getStructuredLogParserForComponent(result.ComponentName)
	result.StructuredBody = parser.TryParse(messageSource)

	return result, nil
}

func getStructuredLogParserForComponent(componentName string) logutil.StructuredLogParser {
	switch componentName {
	case "kubelet":
		return logutil.NewMultiTextLogParser(
			logutil.NewKLogTextParser(true),
			&logutil.FallbackRawTextLogParser{},
		)
	case "kube-controller-manager":
		return logutil.NewMultiTextLogParser(
			logutil.NewKLogTextParser(false),
			&logutil.FallbackRawTextLogParser{},
		)
	case "containerd":
		return logutil.NewMultiTextLogParser(
			logutil.NewLogfmtTextParser(),
			&logutil.FallbackRawTextLogParser{},
		)
	default:
		return logutil.NewMultiTextLogParser(
			logutil.NewKLogTextParser(false),
			&logutil.FallbackRawTextLogParser{},
		)
	}
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
