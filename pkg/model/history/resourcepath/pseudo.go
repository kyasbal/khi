package resourcepath

import (
	"fmt"
)

func Cluster(name string) string {
	if name == "" {
		name = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("@Cluster#controlplane#cluster-scope#%s", name)
}

func Autoscaler(clusterName string) string {
	return fmt.Sprintf("%s#autoscaler", Cluster(clusterName))
}

func Nodepool(clusterName string, nodepoolName string) string {
	if clusterName == "" {
		clusterName = nonSpecifiedPlaceholder
	}
	if nodepoolName == "" {
		nodepoolName = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("@Cluster#nodepool#%s#%s", clusterName, nodepoolName)
}

func Mig(clusterName string, nodepoolName string, migName string) string {
	if clusterName == "" {
		clusterName = nonSpecifiedPlaceholder
	}
	if nodepoolName == "" {
		nodepoolName = nonSpecifiedPlaceholder
	}
	if migName == "" {
		migName = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("@Cluster#nodepool#%s#%s#%s", clusterName, nodepoolName, migName)
}

func NodeComponent(nodeName string, syslogIdentifier string) string {
	if nodeName == "" {
		nodeName = nonSpecifiedPlaceholder
	}
	if syslogIdentifier == "" {
		syslogIdentifier = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("%s#%s", Node(nodeName), syslogIdentifier)
}

func NodeBinding(nodeName string, podNamespace string, podName string) string {
	if nodeName == "" {
		nodeName = nonSpecifiedPlaceholder
	}
	if podName == "" {
		podName = nonSpecifiedPlaceholder
	}
	if podNamespace == "" {
		podNamespace = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("%s#%s(%s)", Node(nodeName), podName, podNamespace)
}

func PodEndpointSlice(endpointSliceNamespace string, endpointSliceName string, podNamespace string, serviceName string) string {
	if podNamespace == "" {
		podNamespace = nonSpecifiedPlaceholder
	}
	if serviceName == "" {
		serviceName = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("%s#%s(%s)[EndpointSlice]", Pod(podNamespace, serviceName), endpointSliceName, endpointSliceNamespace)
}

func ServiceEndpointSlice(namespace string, endpointSliceName string, serviceName string) string {
	if namespace == "" {
		namespace = nonSpecifiedPlaceholder
	}
	if serviceName == "" {
		serviceName = nonSpecifiedPlaceholder
	}
	if endpointSliceName == "" {
		endpointSliceName = nonSpecifiedPlaceholder
	}
	return fmt.Sprintf("%s#%s(%s)[EndpointSlice]", Service(namespace, serviceName), endpointSliceName, namespace)
}
