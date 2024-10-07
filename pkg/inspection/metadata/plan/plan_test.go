package plan

import (
	"testing"

	metadata_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/metadata"
)

func TestConformance(t *testing.T) {
	metadata_test.ConformanceMetadataTypeTest(t, &InspectionPlan{})
}
