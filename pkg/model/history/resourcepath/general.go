package resourcepath

import "fmt"

var PlaceholderForEmptyField = "unknown"

func APIVersionLayerGeneralItem(apiVersion string) string {
	if apiVersion == "" {
		apiVersion = PlaceholderForEmptyField
	}
	if apiVersion == "v1" {
		apiVersion = "core/v1"
	}
	return apiVersion
}

func KindLayerGeneralItem(apiVersion string, kind string) string {
	if kind == "" {
		kind = PlaceholderForEmptyField
	}
	return fmt.Sprintf("%s#%s", APIVersionLayerGeneralItem(apiVersion), kind)
}

func NamespaceLayerGeneralItem(apiVersion string, kind string, namespace string) string {
	if namespace == "" {
		namespace = PlaceholderForEmptyField
	}
	return fmt.Sprintf("%s#%s", KindLayerGeneralItem(apiVersion, kind), namespace)
}

func NameLayerGeneralItem(apiVersion string, kind string, namespace string, name string) string {
	if name == "" {
		name = PlaceholderForEmptyField
	}
	return fmt.Sprintf("%s#%s", NamespaceLayerGeneralItem(apiVersion, kind, namespace), name)
}

func SubresourceLayerGeneralItem(apiVersion string, kind string, namespace string, name string, subresource string) string {
	if subresource == "" {
		subresource = PlaceholderForEmptyField
	}
	return fmt.Sprintf("%s#%s", NameLayerGeneralItem(apiVersion, kind, namespace, name), subresource)
}
