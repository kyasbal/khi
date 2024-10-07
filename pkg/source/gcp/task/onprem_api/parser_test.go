package onprem_api

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseResourceNameOfOnPremAPI(t *testing.T) {
	// Define a struct to hold test cases
	type testCase struct {
		resourceName string
		expected     *onpremResource
	}

	// Create test cases with various input scenarios
	var testCases = []testCase{
		// Valid cases
		{
			resourceName: "projects/12345/locations/asia-northeast1/baremetalClusters/my-cluster",
			expected: &onpremResource{
				ClusterName:  "my-cluster",
				NodepoolName: "",
				ClusterType:  "baremetal",
			},
		},
		{
			resourceName: "projects/67890/locations/us-central1/vmwareClusters/dev-cluster/vmwareNodePools/pool-1",
			expected: &onpremResource{
				ClusterName:  "dev-cluster",
				NodepoolName: "pool-1",
				ClusterType:  "vmware",
			},
		},
		{ // No cluster name
			resourceName: "projects/12345/locations/asia-northeast1",
			expected: &onpremResource{
				ClusterName:  "unknown",
				NodepoolName: "",
				ClusterType:  "unknown",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.resourceName, func(t *testing.T) {
			result := parseResourceNameOfOnPremAPI(tc.resourceName)

			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("Failed for resourceName: %s\nDifference:\n%s", tc.resourceName, diff)
			}
		})
	}

}
