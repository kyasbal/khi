package v2logconvert

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/k8saudittask"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var Task = inspection_task.NewInspectionProcessor(k8saudittask.LogConvertTaskId, []string{
	inspection_task.BuilderGeneratorTask.ID().String(),
	k8saudittask.K8sAuditQueryTaskId,
}, func(ctx context.Context, taskMode int, v *task.VariableSet, tp *progress.TaskProgress) (any, error) {
	if taskMode == inspection_task.TaskModeDryRun {
		return struct{}{}, nil
	}
	builder, err := inspection_task.GetHistoryBuilderFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	logs, err := task.GetTypedVariableFromTaskVariable[[]*log.LogEntity](v, k8saudittask.K8sAuditQueryTaskId, nil)
	if err != nil {
		return nil, err
	}
	processedCount := atomic.Int32{}
	updator := progress.NewProgressUpdator(tp, time.Second, func(tp *progress.TaskProgress) {
		current := processedCount.Load()
		tp.Percentage = float32(current) / float32(len(logs))
		tp.Message = fmt.Sprintf("%d/%d", current, len(logs))
	})
	err = updator.Start(ctx)
	if err != nil {
		return nil, err
	}
	defer updator.Done()
	err = builder.PrepareParseLogs(ctx, logs, func() {
		processedCount.Add(1)
	})
	if err != nil {
		return nil, err
	}
	return struct{}{}, nil
})
