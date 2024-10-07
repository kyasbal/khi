package resourcepath

import (
	"testing"
)

func TestContainer(t *testing.T) {
	testCases := []struct {
		name          string
		namespace     string
		podName       string
		containerName string
		expected      string
	}{
		{"All specified", "my-namespace", "my-pod", "my-container", "core/v1#pod#my-namespace#my-pod#my-container"},
		{"Empty namespace", "", "my-pod", "my-container", "core/v1#pod#unknown#my-pod#my-container"},
		{"Empty pod name", "my-namespace", "", "my-container", "core/v1#pod#my-namespace#unknown#my-container"},
		{"Empty container name", "my-namespace", "my-pod", "", "core/v1#pod#my-namespace#my-pod#unknown"},
		{"Two empty", "", "", "my-container", "core/v1#pod#unknown#unknown#my-container"},
		{"Two empty #2", "my-namespace", "", "", "core/v1#pod#my-namespace#unknown#unknown"},
		{"All empty", "", "", "", "core/v1#pod#unknown#unknown#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Container(tc.namespace, tc.podName, tc.containerName)
			if result != tc.expected {
				t.Errorf("Container function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestPod(t *testing.T) {
	testCases := []struct {
		name      string
		namespace string
		podName   string
		expected  string
	}{
		{"All specified", "my-namespace", "my-pod", "core/v1#pod#my-namespace#my-pod"},
		{"Empty namespace", "", "my-pod", "core/v1#pod#unknown#my-pod"},
		{"Empty pod name", "my-namespace", "", "core/v1#pod#my-namespace#unknown"},
		{"Both empty", "", "", "core/v1#pod#unknown#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Pod(tc.namespace, tc.podName)
			if result != tc.expected {
				t.Errorf("Pod function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestNode(t *testing.T) {
	testCases := []struct {
		name     string
		nodeName string
		expected string
	}{
		{"Node name specified", "my-node", "core/v1#node#cluster-scope#my-node"},
		{"Empty node name", "", "core/v1#node#cluster-scope#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Node(tc.nodeName)
			if result != tc.expected {
				t.Errorf("Node function failed. Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
