package multicloud_api

import (
	"testing"

	gcp_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateMultiCloudAPIQuery(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected string
	}{
		{
			Input: "awsClusters/cluster-foo",
			Expected: `resource.type="audited_resource"
resource.labels.service="gkemulticloud.googleapis.com"
resource.labels.method:("Update" OR "Create" OR "Delete")
protoPayload.resourceName:"awsClusters/cluster-foo"
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			actual := GenerateMultiCloudAPIQuery(testCase.Input)
			if diff := cmp.Diff(testCase.Expected, actual); diff != "" {
				t.Errorf("The generated result is not matching with the expected\n%s", diff)
			}
		})
	}
}

func TestGenerateMultiCloudAPIQueryIsValid(t *testing.T) {
	testCases := []struct {
		Name        string
		ClusterName string
	}{
		{
			Name:        "Valid Query",
			ClusterName: "awsClusters/cluster-foo",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			query := GenerateMultiCloudAPIQuery(tc.ClusterName)
			err := gcp_test.IsValidLogQuery(query)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
