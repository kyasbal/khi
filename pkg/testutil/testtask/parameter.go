package testtask

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"

// TestRunTaskParameterOpt is type used for the Functional Option Pattern on RunSingleTask
type TestRunTaskParameterOpt interface {
	AddParam(params map[string]any)
}

type priorTaskResultOpt struct {
	taskId    string
	parameter any
}

// AddParam implements RunSingleTaskParameterOpt.
func (p *priorTaskResultOpt) AddParam(params map[string]any) {
	params[p.taskId] = p.parameter
}

var _ TestRunTaskParameterOpt = (*priorTaskResultOpt)(nil)

// PriorTaskResult returns RunSingleTaskParameterOpt to fill a parameter of a task result with given value.
func PriorTaskResult(task task.Definition, parameter any) TestRunTaskParameterOpt {
	return &priorTaskResultOpt{
		taskId:    task.ID().String(),
		parameter: parameter,
	}
}

// PriorTaskResultFromID returns RunSingleTaskParameterOpt to fill a parameter of a task result with given value.
func PriorTaskResultFromID(id string, parameter any) TestRunTaskParameterOpt {
	return &priorTaskResultOpt{
		taskId:    id,
		parameter: parameter,
	}
}
