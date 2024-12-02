package serialport

import (
	"fmt"
	"math/rand"
	"testing"

	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"

	gcp_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateSerialPortQuery(t *testing.T) {
	testCases := []struct {
		name      string
		taskMode  int
		nodeNames []string
		wantQuery string
	}{
		{
			name:      "dry run mode",
			taskMode:  inspection_task.TaskModeDryRun,
			nodeNames: []string{},
			wantQuery: `LOG_ID("serialconsole.googleapis.com%2Fserial_port_1_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_2_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_3_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_debug_output")

-- instance name filters to be determined after audit log query`,
		},
		{
			name:      "with a single node name",
			taskMode:  inspection_task.TaskModeRun,
			nodeNames: []string{"node-1"},
			wantQuery: `LOG_ID("serialconsole.googleapis.com%2Fserial_port_1_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_2_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_3_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_debug_output")

labels."compute.googleapis.com/resource_name"=("node-1")`,
		},
		{
			name:      "with multiple node names",
			taskMode:  inspection_task.TaskModeRun,
			nodeNames: []string{"node-1", "node-2", "node-3"},
			wantQuery: `LOG_ID("serialconsole.googleapis.com%2Fserial_port_1_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_2_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_3_output") OR
LOG_ID("serialconsole.googleapis.com%2Fserial_port_debug_output")

labels."compute.googleapis.com/resource_name"=("node-1" OR "node-2" OR "node-3")`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := GenerateSerialPortQuery(tc.taskMode, tc.nodeNames)
			if diff := cmp.Diff(tc.wantQuery, query[0]); diff != "" {
				t.Errorf("the generated query is not matching with the expected query\n%s", diff)
			}
			err := gcp_test.IsValidLogQuery(query[0])
			if err != nil {
				t.Errorf("the generated query is invalid. error:%v", err)
			}
		})
	}
}

func TestMaximumNodeCountNotHittingQueryLengthLimit(t *testing.T) {
	nodeNames := []string{}
	for i := 0; i < MaxNodesPerQuery*2+1; i++ { // This query must be splitted with 3 sub groups.
		nodeNames = append(nodeNames, fmt.Sprintf(`gke-%s-%s-%s`, randomString(46), randomString(8), randomString(4)))
	}
	query := GenerateSerialPortQuery(inspection_task.TaskModeRun, nodeNames)
	if len(query) != 3 {
		t.Errorf("len(GenerateSerialPortQuery())=%d, want %d", len(query), 3)
	}
	for _, subquery := range query {
		err := gcp_test.IsValidLogQuery(subquery)
		if err != nil {
			t.Errorf("the generated query is invalid. error:%v", err)
		}
	}
}

func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	randomid := make([]rune, length)
	for i := range randomid {
		randomid[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomid)
}
