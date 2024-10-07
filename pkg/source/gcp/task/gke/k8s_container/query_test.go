package k8s_container

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query/queryutil"
	gcp_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp"
)

func TestGenerateK8sContainerQueryIsValid(t *testing.T) {
	testCases := []struct {
		Name            string
		ClusterName     string
		PodNameFilter   *queryutil.SetFilterParseResult
		NamespaceFilter *queryutil.SetFilterParseResult
	}{
		{
			Name:            "with no set filters",
			ClusterName:     "foo-cluster",
			PodNameFilter:   &queryutil.SetFilterParseResult{Additives: []string{}},
			NamespaceFilter: &queryutil.SetFilterParseResult{Additives: []string{}},
		},
		{
			Name:            "with namespace filter",
			ClusterName:     "foo-cluster",
			PodNameFilter:   &queryutil.SetFilterParseResult{Additives: []string{}},
			NamespaceFilter: &queryutil.SetFilterParseResult{Additives: []string{"kube-system"}},
		},
		{
			Name:            "with pod name filter",
			ClusterName:     "foo-cluster",
			PodNameFilter:   &queryutil.SetFilterParseResult{Additives: []string{"nginx-pod"}},
			NamespaceFilter: &queryutil.SetFilterParseResult{Additives: []string{}},
		},
		{
			Name:            "with both filters",
			ClusterName:     "foo-cluster",
			PodNameFilter:   &queryutil.SetFilterParseResult{Additives: []string{"nginx-pod"}},
			NamespaceFilter: &queryutil.SetFilterParseResult{Additives: []string{"kube-system"}},
		},
		{
			Name:            "with complex filters",
			ClusterName:     "foo-cluster",
			PodNameFilter:   &queryutil.SetFilterParseResult{Additives: []string{"nginx-pod", "apache-pod"}},
			NamespaceFilter: &queryutil.SetFilterParseResult{Additives: []string{"kube-system", "istio-system"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			query := GenerateK8sContainerQuery(tc.ClusterName, tc.PodNameFilter, tc.NamespaceFilter)
			err := gcp_test.IsValidLogQuery(query)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
