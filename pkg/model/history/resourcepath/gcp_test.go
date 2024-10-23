package resourcepath

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

func TestNetworkEndpointGroup(t *testing.T) {
	testCases := []struct {
		name         string
		negNamespace string
		negName      string
		expected     string
	}{
		{"NEG name specified", "my-namespace", "my-neg", "networking.gke.io/v1beta1#servicenetworkendpointgroup#my-namespace#my-neg"},
		{"Empty NEG namespace", "", "my-neg", "networking.gke.io/v1beta1#servicenetworkendpointgroup#unknown#my-neg"},
		{"Empty NEG name", "my-namespace", "", "networking.gke.io/v1beta1#servicenetworkendpointgroup#my-namespace#unknown"},
		{"Both empty", "", "", "networking.gke.io/v1beta1#servicenetworkendpointgroup#unknown#unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NetworkEndpointGroup(tc.negNamespace, tc.negName)
			if result.Path != tc.expected {
				t.Errorf("NetworkEndpointGroup(%s,%s).Path=%q, want %q", tc.negNamespace, tc.negName, result, tc.expected)
			}
			if result.ParentRelationship != enum.RelationshipChild {
				t.Errorf("NetworkEndpointGroup(%s,%s).ParentRelationshiop=%q, want %q", tc.negNamespace, tc.negName, result.ParentRelationship, enum.RelationshipChild)
			}
		})
	}
}
