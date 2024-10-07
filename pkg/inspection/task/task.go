package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

type InspectionProcessorFunc = func(ctx context.Context, taskMode int, v *task.VariableSet, progress *progress.TaskProgress) (any, error)

// NewInspectionProcessor generates a processor task.Definition with progress reporting feature
func NewInspectionProcessor(taskId string, dependencies []string, processor InspectionProcessorFunc, labelOpts ...task.LabelOpt) task.Definition {
	return task.NewProcessorTask(taskId, dependencies, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		md, err := GetMetadataSetFromVariable(v)
		if err != nil {
			return nil, err
		}
		progressRaw := md.LoadOrStore(progress.ProgressMetadataKey, &progress.ProgressMetadataFactory{})
		p := progressRaw.(*progress.Progress)
		defer p.ResolveTask(taskId)
		taskProgress, err := p.GetTaskProgress(taskId)
		if err != nil {
			return nil, err
		}
		return processor(ctx, taskMode, v, taskProgress)

	}, append([]task.LabelOpt{&ProgressReportableTaskLabelOptImpl{}}, labelOpts...)...)
}

// NewInspectionCachedProcessor generates a cached processor task.Definition with progress reporting feature
func NewInspectionCachedProcessor(taskId string, dependencies []string, processor InspectionProcessorFunc, labelOpts ...task.LabelOpt) task.Definition {
	return task.NewCachedProcessor(taskId, dependencies, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		md, err := GetMetadataSetFromVariable(v)
		if err != nil {
			return nil, err
		}
		progressRaw := md.LoadOrStore(progress.ProgressMetadataKey, &progress.ProgressMetadataFactory{})
		p := progressRaw.(*progress.Progress)
		defer p.ResolveTask(taskId)
		taskProgress, err := p.GetTaskProgress(taskId)
		if err != nil {
			return nil, err
		}
		return processor(ctx, taskMode, v, taskProgress)

	}, append([]task.LabelOpt{&ProgressReportableTaskLabelOptImpl{}}, labelOpts...)...)
}

type InspectionProducerFunc = func(ctx context.Context, taskMode int, progress *progress.TaskProgress) (any, error)

// NewInspectionProducer generates a producer task.Definition with progress reporting feature
func NewInspectionProducer(taskId string, producer InspectionProducerFunc, labelOpts ...task.LabelOpt) task.Definition {
	return task.NewProcessorTask(taskId, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		md, err := GetMetadataSetFromVariable(v)
		if err != nil {
			return nil, err
		}
		progressRaw := md.LoadOrStore(progress.ProgressMetadataKey, &progress.ProgressMetadataFactory{})
		p := progressRaw.(*progress.Progress)
		defer p.ResolveTask(taskId)
		taskProgress, err := p.GetTaskProgress(taskId)
		if err != nil {
			return nil, err
		}
		return producer(ctx, taskMode, taskProgress)

	}, append([]task.LabelOpt{&ProgressReportableTaskLabelOptImpl{}}, labelOpts...)...)
}
