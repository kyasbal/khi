package resourcepath

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

var PlaceholderForEmptyField = "unknown"

func APIVersionLayerGeneralItem(apiVersion string) ResourcePath {
	if apiVersion == "" {
		apiVersion = PlaceholderForEmptyField
	}
	if apiVersion == "v1" {
		apiVersion = "core/v1"
	}
	return ResourcePath{
		Path:               apiVersion,
		ParentRelationship: enum.RelationshipChild,
	}
}

func KindLayerGeneralItem(apiVersion, kind string) ResourcePath {
	if kind == "" {
		kind = PlaceholderForEmptyField
	}
	return ResourcePath{
		Path:               fmt.Sprintf("%s#%s", APIVersionLayerGeneralItem(apiVersion).Path, kind),
		ParentRelationship: enum.RelationshipChild,
	}
}

func NamespaceLayerGeneralItem(apiVersion, kind, namespace string) ResourcePath {
	if namespace == "" {
		namespace = PlaceholderForEmptyField
	}
	return ResourcePath{
		Path:               fmt.Sprintf("%s#%s", KindLayerGeneralItem(apiVersion, kind).Path, namespace),
		ParentRelationship: enum.RelationshipChild,
	}
}

func NameLayerGeneralItem(apiVersion, kind, namespace, name string) ResourcePath {
	if name == "" {
		name = PlaceholderForEmptyField
	}
	return ResourcePath{
		Path:               fmt.Sprintf("%s#%s", NamespaceLayerGeneralItem(apiVersion, kind, namespace).Path, name),
		ParentRelationship: enum.RelationshipChild,
	}
}

func SubresourceLayerGeneralItem(apiVersion, kind, namespace, name, subresource string) ResourcePath {
	if subresource == "" {
		subresource = PlaceholderForEmptyField
	}
	return ResourcePath{
		Path:               fmt.Sprintf("%s#%s", NameLayerGeneralItem(apiVersion, kind, namespace, name).Path, subresource),
		ParentRelationship: enum.RelationshipChild,
	}
}

func FromK8sOperation(op model.KubernetesObjectOperation) ResourcePath {
	var path string
	if op.SubResourceName != "" {
		path = strings.ToLower(strings.Join([]string{
			op.APIVersion,
			op.GetSingularKindName(),
			op.Namespace,
			op.Name,
			op.SubResourceName,
		}, "#"))
	} else {
		path = strings.ToLower(strings.Join([]string{
			op.APIVersion,
			op.GetSingularKindName(),
			op.Namespace,
			op.Name,
		}, "#"))
	}
	return ResourcePath{
		Path:               path,
		ParentRelationship: enum.RelationshipChild,
	}
}
