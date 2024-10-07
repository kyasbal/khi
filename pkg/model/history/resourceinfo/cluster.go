package resourceinfo

import (
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourceinfo/resourcelease"
	v1 "k8s.io/api/core/v1"
)

// Cluster stores resource information(node name, Pod IP,Host IP...etc) used from another parser.
// This struct must modify the own fields in thread safe.
type Cluster struct {
	lock           sync.Mutex
	nodeNames      map[string]struct{}
	EndpointSlices *EndpointSliceInfo
	IPs            *resourcelease.ResourceLeaseHistory[*resourcelease.K8sResourceLeaseHolder]
	// records lease history of NEG id to ServiceNetworkEndpointGroup
	NEGs          *resourcelease.ResourceLeaseHistory[*resourcelease.K8sResourceLeaseHolder]
	PodSandboxIds *resourcelease.ResourceLeaseHistory[*resourcelease.K8sResourceLeaseHolder]
	ContainerIds  *resourcelease.ResourceLeaseHistory[*resourcelease.ContainerLeaseHolder]
	// CRIResource       *CRIResourceBinder
	ContainerStatuses *ContainerStatuses
}

func NewClusterResourceInfo() *Cluster {
	ips := resourcelease.NewResourceLeaseHistory[*resourcelease.K8sResourceLeaseHolder]()
	return &Cluster{
		lock:           sync.Mutex{},
		nodeNames:      map[string]struct{}{},
		EndpointSlices: newEndpointSliceInfo(ips),
		IPs:            ips,
		NEGs:           resourcelease.NewResourceLeaseHistory[*resourcelease.K8sResourceLeaseHolder](),
		PodSandboxIds:  resourcelease.NewResourceLeaseHistory[*resourcelease.K8sResourceLeaseHolder](),
		ContainerIds:   resourcelease.NewResourceLeaseHistory[*resourcelease.ContainerLeaseHolder](),
		ContainerStatuses: &ContainerStatuses{
			lastObservedStatus: make(map[string]v1.ContainerStatus),
		},
	}
}

// AddNode registeres the node name
func (c *Cluster) AddNode(nodeName string) {
	defer c.lock.Unlock()
	c.lock.Lock()
	c.nodeNames[nodeName] = struct{}{}
}

// GetNodes returns copy of the list of node names
func (c *Cluster) GetNodes() []string {
	result := []string{}
	for key := range c.nodeNames {
		result = append(result, key)
	}
	return result
}
