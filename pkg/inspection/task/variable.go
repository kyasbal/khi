package task

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const (
	MetadataVariableName          = InspectionTaskPrefix + "metadata"
	InspectionResultVariableName  = InspectionTaskPrefix + "inspection-result"
	InspectionRequestVariableName = InspectionTaskPrefix + "request"
	InspectionIdVariableName      = InspectionTaskPrefix + "inspection-id"
)

func GetHistoryBuilderFromTaskVariable(v *task.VariableSet) (*history.Builder, error) {
	return task.GetTypedVariableFromTaskVariable[*history.Builder](v, BuilderGeneratorTaskId, nil)
}

func GetInspectionIdFromTaskVariable(v *task.VariableSet) (string, error) {
	return task.GetTypedVariableFromTaskVariable[string](v, InspectionIdVariableName, "<INVALID>")
}
