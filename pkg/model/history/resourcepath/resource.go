package resourcepath

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

const nonSpecifiedPlaceholder = "unknown"

func Container(namespace string, name string, containerName string) ResourcePath {
	if namespace == "" {
		namespace = nonSpecifiedPlaceholder
	}
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	if containerName == "" {
		containerName = nonSpecifiedPlaceholder
	}
	containerResourcePath := SubresourceLayerGeneralItem("core/v1", "pod", namespace, name, containerName)
	containerResourcePath.ParentRelationship = enum.RelationshipContainer
	return containerResourcePath
}

func Pod(namespace string, name string) ResourcePath {
	if namespace == "" {
		namespace = nonSpecifiedPlaceholder
	}
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	return NameLayerGeneralItem("core/v1", "pod", namespace, name)
}

func Service(namespace string, name string) ResourcePath {
	if namespace == "" {
		namespace = nonSpecifiedPlaceholder
	}
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	return NameLayerGeneralItem("core/v1", "service", namespace, name)
}

func Node(name string) ResourcePath {
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	return NameLayerGeneralItem("core/v1", "node", "cluster-scope", name)
}
