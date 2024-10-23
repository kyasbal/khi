package resourcepath

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
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
			if result.Path != tc.expected {
				t.Errorf("ControlplneComponent(%s,%s).Path=%q, want %q", tc.clusterName, tc.component, result, tc.expected)
			}
			if result.ParentRelationship != enum.RelationshipControlPlaneComponent {
				t.Errorf("ControlplaneComponent(%s,%s).ParentRelationshiop=%q, want %q", tc.clusterName, tc.component, result.ParentRelationship, enum.RelationshipControlPlaneComponent)
			}
		})
	}
}
