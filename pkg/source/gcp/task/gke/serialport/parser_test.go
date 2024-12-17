// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
