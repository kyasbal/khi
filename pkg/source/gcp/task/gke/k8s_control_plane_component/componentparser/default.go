package componentparser

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

type DefaultK8sControlPlaneComponentParser struct {
}

// Process implements ControlPlaneComponentParser.
func (d *DefaultK8sControlPlaneComponentParser) Process(ctx context.Context, l *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, v *task.VariableSet) (bool, error) {
	component := l.GetStringOrDefault("resource.labels.component_name", "Unknown")
	clusterName := l.GetStringOrDefault("resource.labels.cluster_name", "Unknown")
	msg, err := l.MainMessage()
	if err == nil {
		cs.RecordLogSummary(msg)
	}
	cs.RecordEvent(resourcepath.ControlplaneComponent(clusterName, component))
	return false, nil
}

// ShouldProcess implements ControlPlaneComponentParser.
func (d *DefaultK8sControlPlaneComponentParser) ShouldProcess(component_name string) bool {
	return true
}

var _ ControlPlaneComponentParser = (*DefaultK8sControlPlaneComponentParser)(nil)
