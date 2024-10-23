package resourcepath

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"

func ControlplaneComponent(cluster string, component string) ResourcePath {
	if cluster == "" {
		cluster = nonSpecifiedPlaceholder
	}
	if component == "" {
		component = nonSpecifiedPlaceholder
	}
	controlPlaneComponentResourcePath := SubresourceLayerGeneralItem("@Cluster", "controlplane", "cluster-scope", cluster, component)
	controlPlaneComponentResourcePath.ParentRelationship = enum.RelationshipControlPlaneComponent
	return controlPlaneComponentResourcePath
}
