package gcp

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/common"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	inspection_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/inspection"
	env_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/inspection/env"
)

func testPrepareInspectionServer(inspectionServer *inspection.InspectionTaskServer) error {
	err := commonPreparation(inspectionServer)
	if err != nil {
		return err
	}

	err = inspectionServer.AddTaskDefinition(env_test.MockedEnvironmentVariableProducer(task.EnvGcpIamTokenTask, ""))
	if err != nil {
		return err
	}

	err = inspectionServer.AddTaskDefinition(env_test.MockedEnvironmentVariableProducer(task.EnvFixedProjectIdTask, ""))
	if err != nil {
		return err
	}

	err = inspectionServer.AddTaskDefinition(env_test.MockedEnvironmentVariableProducer(task.EnvQuotaProjectIdTask, ""))
	if err != nil {
		return err
	}

	return nil
}

func TestInspectionTasksAreResolvable(t *testing.T) {
	inspection_test.ConformanceEveryInspectionTasksAreResolvable(t, "gcp", []inspection.PrepareInspectionServerFunc{
		common.PrepareInspectionServer,
		testPrepareInspectionServer,
	})
}

func TestConformanceTestForInspectionTypes(t *testing.T) {
	inspection_test.ConformanceTestForInspectionTypes(t, []inspection.PrepareInspectionServerFunc{
		common.PrepareInspectionServer,
		testPrepareInspectionServer,
	})
}
