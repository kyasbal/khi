package serialport

import (
	"context"
	"fmt"
	"strings"

	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query/queryutil"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/k8saudittask"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const SerialPortLogQueryTaskId = query.GKEQueryPrefix + "serialport"

const MaxNodesPerQuery = 30

func GenerateSerialPortQuery(taskMode int, nodeNames []string) []string {
	if taskMode == inspection_task.TaskModeDryRun {
		return []string{
			generateSerialPortQueryWithInstanceNameFilter("-- instance name filters to be determined after audit log query"),
		}
	} else {
		result := []string{}
		instanceNameGroups := queryutil.SplitToChildGroups(nodeNames, MaxNodesPerQuery)
		for _, group := range instanceNameGroups {
			instanceNameFilter := fmt.Sprintf(`labels."compute.googleapis.com/resource_name"=(%s)`, strings.Join(queryutil.WrapDoubleQuoteForStringArray(group), " OR "))
			result = append(result, generateSerialPortQueryWithInstanceNameFilter(instanceNameFilter))
		}
		return result
	}
}

func generateSerialPortQueryWithInstanceNameFilter(instanceNameFilter string) string {
	return fmt.Sprintf(`LOG_ID("serialconsole.googleapis.com%%2Fserial_port_1_output") OR
LOG_ID("serialconsole.googleapis.com%%2Fserial_port_2_output") OR
LOG_ID("serialconsole.googleapis.com%%2Fserial_port_3_output") OR
LOG_ID("serialconsole.googleapis.com%%2Fserial_port_debug_output")

%s`, instanceNameFilter)
}

var GKESerialPortLogQueryTask = query.NewQueryGeneratorTask(SerialPortLogQueryTaskId, "Serial port log", enum.LogTypeSerialPort, []string{
	k8saudittask.K8sAuditParseTaskId,
}, func(ctx context.Context, taskMode int, vs *task.VariableSet) ([]string, error) {
	builder, err := inspection_task.GetHistoryBuilderFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return GenerateSerialPortQuery(taskMode, builder.ClusterResource.GetNodes()), nil
})
