package task

import (
	"context"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var TimeZoneShiftInputTaskId = GCPPrefix + "input/timezone-shift"

var TimeZoneShiftInputTask = inspection_task.NewInspectionProcessor(TimeZoneShiftInputTaskId, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet, progress *progress.TaskProgress) (any, error) {
	req, err := inspection_task.GetInspectionRequestFromVariable(v)
	if err != nil {
		return nil, err
	}
	if tzShiftAny, found := req.Values["timezoneShift"]; found {
		if tzShiftFloat, convertible := tzShiftAny.(float64); convertible {
			return time.FixedZone("Unknown", int(tzShiftFloat*3600)), nil
		} else {
			return time.UTC, nil
		}
	} else {
		return time.UTC, nil
	}
})

func GetTimezoneShiftInput(tv *task.VariableSet) (*time.Location, error) {
	return task.GetTypedVariableFromTaskVariable[*time.Location](tv, TimeZoneShiftInputTaskId, time.UTC)
}
