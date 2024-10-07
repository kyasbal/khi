package resourcepath

import (
	"fmt"
)

const nonSpecifiedPlaceholder = "unknown"

func Container(namespace string, name string, containerName string) string {
	if namespace == "" {
		namespace = nonSpecifiedPlaceholder
	}
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	if containerName == "" {
		containerName = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("core/v1#pod#%s#%s#%s", namespace, name, containerName)
}

func Pod(namespace string, name string) string {
	if namespace == "" {
		namespace = nonSpecifiedPlaceholder
	}
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("core/v1#pod#%s#%s", namespace, name)
}

func Service(namespace string, name string) string {
	return NameLayerGeneralItem("core/v1", "service", namespace, name)
}

func Node(name string) string {
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("core/v1#node#cluster-scope#%s", name)
}
