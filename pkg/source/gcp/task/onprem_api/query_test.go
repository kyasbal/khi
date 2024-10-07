package onprem_api

import (
	"testing"

	gcp_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateOnPremAPIQuery(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected string
	}{
		{
			Input: "baremetalClusters/my-cluster",
			Expected: `resource.type="audited_resource"
resource.labels.service="gkeonprem.googleapis.com"
resource.labels.method:("Update" OR "Create" OR "Delete" OR "Enroll" OR "Unenroll")
protoPayload.resourceName:"baremetalClusters/my-cluster"
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			actual := GenerateOnPremAPIQuery(testCase.Input)
			if diff := cmp.Diff(testCase.Expected, actual); diff != "" {
				t.Errorf("The generated result is not matching with the expected\n%s", diff)
			}
		})
	}
}

func TestGenerateOnPremAPIQueryIsValid(t *testing.T) {
	testCases := []struct {
		Name        string
		ClusterName string
	}{
		{
			Name:        "Valid Query",
			ClusterName: "baremetalClusters/my-cluster",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			query := GenerateOnPremAPIQuery(tc.ClusterName)
			err := gcp_test.IsValidLogQuery(query)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
