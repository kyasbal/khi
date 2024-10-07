package serializer

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/inspectiondata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/ioconfig"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const SerializerTaskId = inspection_task.InspectionTaskPrefix + "serialize"

var SerializeTask = inspection_task.NewInspectionProcessor(SerializerTaskId, []string{inspection_task.InspectionMainSubgraphName + "-done", ioconfig.IOConfigTaskName, inspection_task.BuilderGeneratorTaskId}, func(ctx context.Context, taskMode int, v *task.VariableSet, progress *progress.TaskProgress) (any, error) {
	if taskMode == inspection_task.TaskModeDryRun {
		slog.DebugContext(ctx, "Skipping because this is in dryrun mode")
		return nil, nil
	}
	taskId, err := inspection_task.GetInspectionIdFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	ioConfig, err := ioconfig.GetIOConfigFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	builder, err := inspection_task.GetHistoryBuilderFromTaskVariable(v)
	if err != nil {
		return nil, err
	}
	store := inspectiondata.NewFileSystemInspectionResultRepository(filepath.Join(ioConfig.DataDestination, taskId+".khi"))
	writer, err := store.GetWriter()
	if err != nil {
		return nil, err
	}
	metadataSet, err := inspection_task.GetMetadataSetFromVariable(v)
	if err != nil {
		return nil, err
	}
	resultMetadata, err := metadataSet.ToMap(task.EqualLabelFilter(metadata.LabelKeyIncludedInResultBinaryFlag, true, false))
	if err != nil {
		return nil, err
	}
	err = builder.Finalize(ctx, resultMetadata, writer, progress)
	if err != nil {
		return nil, err
	}
	err = store.Close()
	if err != nil {
		return nil, err
	}
	return store, nil
})
