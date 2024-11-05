package gcp

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/common"
	inspection_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/inspection"
)

func testPrepareInspectionServer(inspectionServer *inspection.InspectionTaskServer) error {
	err := commonPreparation(inspectionServer)
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
