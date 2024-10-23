package k8s_event

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	parser_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/parser"
)

func TestK8sEventParser_ParseSampleLog(t *testing.T) {
	wantLogSummary := "【NodeRegistrationCheckerDidNotRunChecks】Fri Sep 13 01:49:48 UTC 2024 - **     Node ready and registered. **"
	cs, err := parser_test.ParseFromYamlLogFile("test/logs/k8s_event/sample.yaml", &k8sEventParser{}, nil, nil)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	event := cs.GetEvents(resourcepath.Node("gke-gke-basic-1-default-5e5b794d-89xl"))
	if len(event) != 1 {
		t.Errorf("got %d events, want 1", len(event))
	}

	gotLogSummary := cs.GetLogSummary()
	if gotLogSummary != wantLogSummary {
		t.Errorf("got %q log summary, want %q", gotLogSummary, wantLogSummary)
	}
}

func TestK8sEventParser_ClusterScope(t *testing.T) {
	wantLogSummary := "Event exporter started watching. Some events may have been lost up to this point."
	cs, err := parser_test.ParseFromYamlLogFile("test/logs/k8s_event/cluster-scoped.yaml", &k8sEventParser{}, nil, nil)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	event := cs.GetEvents(resourcepath.Cluster("gke-basic-1"))
	if len(event) != 1 {
		t.Errorf("got %d events, want 1", len(event))
	}

	gotLogSummary := cs.GetLogSummary()
	if gotLogSummary != wantLogSummary {
		t.Errorf("got %q log summary, want %q", gotLogSummary, wantLogSummary)
	}
}
