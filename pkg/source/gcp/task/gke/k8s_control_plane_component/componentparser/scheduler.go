package componentparser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var ErrParserNoMatchingWithLog = errors.New("Parser didn't match with the given log")

type SchedulerComponentParser struct{}

// Process implements ControlPlaneComponentParser.
func (s *SchedulerComponentParser) Process(ctx context.Context, l *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, v *task.VariableSet) (bool, error) {
	path, err := s.podRelatedLogsToResourcePath(ctx, l)
	if err == nil {
		cs.RecordEvent(path)
	}
	return true, nil
}

// ShouldProcess implements ControlPlaneComponentParser.
func (s *SchedulerComponentParser) ShouldProcess(component_name string) bool {
	return component_name == "scheduler"
}

func (s *SchedulerComponentParser) podRelatedLogsToResourcePath(ctx context.Context, l *log.LogEntity) (string, error) {
	hasPodField := l.HasKLogField("pod")
	if hasPodField {
		pod, err := l.KLogField("pod")
		if err != nil {
			return "", ErrParserNoMatchingWithLog
		}
		splittedPodName := strings.Split(pod, "/")
		if len(splittedPodName) != 2 {
			slog.WarnContext(ctx, fmt.Sprintf("Unexpected pod klog format: %s", pod))
			return "", ErrParserNoMatchingWithLog
		}
		return resourcepath.Pod(splittedPodName[0], splittedPodName[1]), nil
	}
	return "", ErrParserNoMatchingWithLog
}

var _ ControlPlaneComponentParser = (*SchedulerComponentParser)(nil)
