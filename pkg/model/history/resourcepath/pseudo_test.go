package resourcepath

import (
	"testing"
)

func TestCluster(t *testing.T) {
	testCases := []struct {
		name        string
		clusterName string
		expected    string
	}{
		{"Cluster name specified", "my-cluster", "@Cluster#controlplane#cluster-scope#my-cluster"},
		{"Empty cluster name", "", "@Cluster#controlplane#cluster-scope#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Cluster(tc.clusterName)
			if result != tc.expected {
				t.Errorf("Cluster function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestAutoscaler(t *testing.T) {
	testCases := []struct {
		name        string
		clusterName string
		expected    string
	}{
		{"Cluster name specified", "my-cluster", "@Cluster#controlplane#cluster-scope#my-cluster#autoscaler"},
		{"Empty cluster name", "", "@Cluster#controlplane#cluster-scope#unknown#autoscaler"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Autoscaler(tc.clusterName)
			if result != tc.expected {
				t.Errorf("Autoscaler function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestNodepool(t *testing.T) {
	testCases := []struct {
		name         string
		clusterName  string
		nodepoolName string
		expected     string
	}{
		{"All specified", "my-cluster", "my-nodepool", "@Cluster#nodepool#my-cluster#my-nodepool"},
		{"Empty cluster name", "", "my-nodepool", "@Cluster#nodepool#unknown#my-nodepool"},
		{"Empty nodepool name", "my-cluster", "", "@Cluster#nodepool#my-cluster#unknown"},
		{"Both empty", "", "", "@Cluster#nodepool#unknown#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Nodepool(tc.clusterName, tc.nodepoolName)
			if result != tc.expected {
				t.Errorf("Nodepool function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestMig(t *testing.T) {
	testCases := []struct {
		name         string
		clusterName  string
		nodepoolName string
		migName      string
		expected     string
	}{
		{"All specified", "cluster", "nodepool", "mig", "@Cluster#nodepool#cluster#nodepool#mig"},
		{"Empty cluster name", "", "nodepool", "mig", "@Cluster#nodepool#unknown#nodepool#mig"},
		{"Empty nodepool name", "cluster", "", "mig", "@Cluster#nodepool#cluster#unknown#mig"},
		{"Empty mig name", "cluster", "nodepool", "", "@Cluster#nodepool#cluster#nodepool#unknown"},
		{"Two empty", "", "nodepool", "", "@Cluster#nodepool#unknown#nodepool#unknown"},
		{"Two empty #2", "cluster", "", "", "@Cluster#nodepool#cluster#unknown#unknown"},
		{"All empty", "", "", "", "@Cluster#nodepool#unknown#unknown#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Mig(tc.clusterName, tc.nodepoolName, tc.migName)
			if result != tc.expected {
				t.Errorf("Mig function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestNodeComponent(t *testing.T) {
	testCases := []struct {
		name             string
		nodeName         string
		syslogIdentifier string
		expected         string
	}{
		{"All specified", "my-node", "kubelet", "core/v1#node#cluster-scope#my-node#kubelet"},
		{"Empty node name", "", "kubelet", "core/v1#node#cluster-scope#unknown#kubelet"},
		{"Empty syslog identifier", "my-node", "", "core/v1#node#cluster-scope#my-node#unknown"},
		{"Both empty", "", "", "core/v1#node#cluster-scope#unknown#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NodeComponent(tc.nodeName, tc.syslogIdentifier)
			if result != tc.expected {
				t.Errorf("NodeComponent function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
