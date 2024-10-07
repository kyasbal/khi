package componentparser

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

// ControlPlaneComponentParser is an abstraction type to define the customized process for each components
type ControlPlaneComponentParser interface {
	// ShouldProcess return if the component must be processed by this parser or not
	ShouldProcess(component_name string) bool
	// Process handle the given logs to ingest to the ChangeSet. This method return false if the logs shouldn't be processed in the later parsers.
	Process(ctx context.Context, l *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, v *task.VariableSet) (bool, error)
}

var ComponentParsers []ControlPlaneComponentParser = []ControlPlaneComponentParser{
	&ControllerManagerComponentParser{},
	&SchedulerComponentParser{},
	&DefaultK8sControlPlaneComponentParser{},
}
