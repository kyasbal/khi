package inspection_test

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/inspectiondata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/form"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/header"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
	task_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/task"
)

// DryRunInspectionTaskGraph executes the task graph just with provided dependency tasks with several variables given in default local runner
func DryRunInspectionTaskGraph(target task.Definition, requestParams map[string]any, dependencies ...task.Definition) (*task.VariableSet, error) {
	ms := metadata.NewSet()
	ms.LoadOrStore(header.HeaderMetadataKey, &header.HeaderMetadataFactory{})
	ms.LoadOrStore(form.FormFieldSetMetadataKey, &form.FormFieldSetMetadataFactory{})
	return task_test.RunTaskGraph(target, inspection_task.TaskModeDryRun, map[string]any{
		inspection_task.MetadataVariableName: ms,
		inspection_task.InspectionRequestVariableName: &inspection_task.InspectionRequest{
			Values: requestParams,
		},
	}, dependencies...)
}

// RunInspectionTaskGraph executes the task graph just with provided dependency tasks with several variables given in the default local runner
func RunInspectionTaskGraph(target task.Definition, requestParams map[string]any, dependencies ...task.Definition) (*task.VariableSet, error) {
	ms := metadata.NewSet()
	ms.LoadOrStore(header.HeaderMetadataKey, &header.HeaderMetadataFactory{})
	ms.LoadOrStore(form.FormFieldSetMetadataKey, &form.FormFieldSetMetadataFactory{})

	return task_test.RunTaskGraph(target, inspection_task.TaskModeRun, map[string]any{
		inspection_task.MetadataVariableName: ms,
		inspection_task.InspectionRequestVariableName: &inspection_task.InspectionRequest{
			Values: requestParams,
		},
		inspection_task.InspectionResultVariableName: inspectiondata.NewFileSystemInspectionResultRepository("foo"),
	}, dependencies...)
}
