package common

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
	inspection_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/inspection"
)

func TestInspectionTasksAreResolvable(t *testing.T) {
	inspection_test.ConformanceEveryInspectionTasksAreResolvable(t, "common", []inspection.PrepareInspectionServerFunc{PrepareInspectionServer})
}
