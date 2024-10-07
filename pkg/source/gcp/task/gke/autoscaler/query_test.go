package autoscaler

import (
	"testing"

	gcp_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp"
)

func TestGenerateAutoscalerQuery(t *testing.T) {
	testCases := []struct {
		projectId     string
		clusterName   string
		excludeStatus bool
		expected      string
	}{
		{
			projectId:     "my-project",
			clusterName:   "my-cluster",
			excludeStatus: false,
			expected: `resource.type="k8s_cluster"
resource.labels.project_id="my-project"
resource.labels.cluster_name="my-cluster"
-- include query for status log
logName="projects/my-project/logs/container.googleapis.com%2Fcluster-autoscaler-visibility"`,
		},
		{
			projectId:     "my-project",
			clusterName:   "my-cluster",
			excludeStatus: true,
			expected: `resource.type="k8s_cluster"
resource.labels.project_id="my-project"
resource.labels.cluster_name="my-cluster"
-jsonPayload.status: ""
logName="projects/my-project/logs/container.googleapis.com%2Fcluster-autoscaler-visibility"`,
		},
	}

	for _, tc := range testCases {
		result := GenerateAutoscalerQuery(tc.projectId, tc.clusterName, tc.excludeStatus)
		if result != tc.expected {
			t.Errorf("Expected query:\n%s\nGot:\n%s", tc.expected, result)
		}
	}
}

func TestGeneratedAutoscalerQueryIsValid(t *testing.T) {
	testCases := []struct {
		Name          string
		ProjectId     string
		ClusterName   string
		ExcludeStatus bool
	}{
		{
			Name:          "Valid Query",
			ProjectId:     "gcp-project-id",
			ClusterName:   "gcp-cluster-name",
			ExcludeStatus: false,
		},
		{
			Name:          "Valid Query with Exclude Status",
			ProjectId:     "gcp-project-id",
			ClusterName:   "gcp-cluster-name",
			ExcludeStatus: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			query := GenerateAutoscalerQuery(tc.ProjectId, tc.ClusterName, tc.ExcludeStatus)
			err := gcp_test.IsValidLogQuery(query)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
