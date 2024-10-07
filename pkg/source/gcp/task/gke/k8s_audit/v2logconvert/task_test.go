package v2logconvert

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/ioconfig"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	gcp_log "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/k8saudittask"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/testlog"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/testtask"
)

func TestLogFillerTask(t *testing.T) {
	builder := history.NewBuilder(&ioconfig.IOConfig{})
	baseLog := `protoPayload:
  authenticationInfo: 
    principalEmail: user@example.com
  methodName: io.k8s.core.v1.pods.create
  resourceName: core/v1/namespaces/default/pods/my-pod
  response:
    '@type': core.k8s.io/v1.Pod
    foo: bar
  status:
    code: 0
timestamp: "2024-01-01T00:00:00+09:00"`
	logOpts := [][]testlog.TestLogOpt{
		{
			testlog.StringField("insertId", "insertid-1"),
		},
		{
			testlog.StringField("insertId", "insertid-2"),
			testlog.StringField("timestamp", "2024-01-01T00:01:00+09:00"),
		},
		{
			testlog.StringField("insertId", "insertid-3"),
			testlog.StringField("timestamp", "2024-01-01T00:02:00+09:00"),
		},
	}
	logs := []*log.LogEntity{}
	for _, opt := range logOpts {
		logs = append(logs, testlog.New(testlog.BaseYaml(baseLog)).With(opt...).MustBuildLogEntity(gcp_log.GCPCommonFieldExtractor{}))
	}

	_, err := testtask.RunSingleTask[struct{}](Task, inspection_task.TaskModeRun,
		testtask.PriorTaskResultFromID(inspection_task.BuilderGeneratorTaskId, builder),
		testtask.PriorTaskResultFromID(k8saudittask.K8sAuditQueryTaskId, logs),
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	for i := 0; i < len(logs); i++ {
		logId := logs[i].ID()
		_, err := builder.GetLog(logId)
		if err != nil {
			t.Errorf("failed to get log %s", logId)
		}
	}
}
