package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
	common_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

// ReaderFactoryGeneratorTask generates the instance of Reader factory to be used in later task.
const ReaderFactoryGeneratorTaskId = InspectionTaskPrefix + "reader-factory-generator"

var ReaderFactoryGeneratorTask = task.NewProcessorTask(ReaderFactoryGeneratorTaskId, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
	return structure.NewReaderFactory(structuredatastore.NewLRUStructureDataStoreFactory()), nil
})

func GetReaderFactoryFromTaskVariable(v *task.VariableSet) (*structure.ReaderFactory, error) {
	return common_task.GetTypedVariableFromTaskVariable[*structure.ReaderFactory](v, ReaderFactoryGeneratorTaskId, nil)
}
