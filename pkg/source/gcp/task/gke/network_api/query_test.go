package network_api

import (
	"testing"

	gcp_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp"
)

func TestGenerateGenerateGCPNetworkAPIQueryIsValid(t *testing.T) {

	testCases := []struct {
		Name string
		NEGs []string
	}{
		{
			Name: "Valid Query",
			NEGs: []string{"neg-1", "neg-2"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			query := GenerateGCPNetworkAPIQuery(0, tc.NEGs)
			err := gcp_test.IsValidLogQuery(query[0])
			if err != nil {
				t.Errorf("Query is not valid: %v", err)
			}
		})
	}
}
