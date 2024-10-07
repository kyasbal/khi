package resourcepath

import (
	"testing"
)

func TestControlplaneComponent(t *testing.T) {
	testCases := []struct {
		name        string
		clusterName string
		component   string
		expected    string
	}{
		{"Component name specified", "cluster-name", "kube-apiserver", "@Cluster#controlplane#cluster-scope#cluster-name#kube-apiserver"},
		{"Empty component name", "cluster-name", "", "@Cluster#controlplane#cluster-scope#cluster-name#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ControlplaneComponent(tc.clusterName, tc.component)
			if result != tc.expected {
				t.Errorf("ControlplaneComponent function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
