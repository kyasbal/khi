package serialport

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	parser_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/parser"
)

func TestSerialPortLogParser_ParseBasicSerialPortLog(t *testing.T) {
	wantLogSummary := "[ OK ] Stopped getty@tty1.service."

	cs, err := parser_test.ParseFromYamlLogFile("test/logs/serialport/basic-serialport-log.yaml", &SerialPortLogParser{}, nil, nil)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

	event := cs.GetEvents(resourcepath.NodeSerialport("gke-sample-cluster-default-abcdefgh-abcd"))
	if len(event) != 1 {
		t.Errorf("got %d events, want 1", len(event))
	}

	gotLogSummary := cs.GetLogSummary()
	if gotLogSummary != wantLogSummary {
		t.Errorf("got %q log summary, want %q", gotLogSummary, wantLogSummary)
	}
}
