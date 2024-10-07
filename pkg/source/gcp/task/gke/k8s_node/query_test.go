package k8s_node

import (
	"testing"

	gcp_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp"
)

func TestGenerateK8sNodeQueryIsValid(t *testing.T) {

	testCases := []struct {
		Name        string
		ClusterName string
		ProjectName string
	}{
		{
			Name:        "Valid Query",
			ClusterName: "gcp-cluster-name",
			ProjectName: "gcp-project-id",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			query := GenerateK8sNodeLogQuery(tc.ProjectName, tc.ClusterName)
			err := gcp_test.IsValidLogQuery(query)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
