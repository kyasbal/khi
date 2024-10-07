package resourceinfo

import (
	"fmt"
	"reflect"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type ContainerStatuses struct {
	lastObservedStatus map[string]corev1.ContainerStatus
	statusMapLock      sync.Mutex
}

func (c *ContainerStatuses) IsNewChange(namespace string, podname string, containerName string, status corev1.ContainerStatus) bool {
	c.statusMapLock.Lock()
	defer c.statusMapLock.Unlock()
	path := fmt.Sprintf("%s#%s#%s", namespace, podname, containerName)
	if last, found := c.lastObservedStatus[path]; !found {
		c.lastObservedStatus[path] = last
		return true
	} else {
		if !reflect.DeepEqual(last, status) {
			c.lastObservedStatus[path] = last
			return true
		} else {
			return false
		}
	}
}
