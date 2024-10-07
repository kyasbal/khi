package progress

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	metadata_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/metadata"
)

func newProgressforConformanceTest() metadata.Metadata {
	progress := (&ProgressMetadataFactory{}).Instanciate().(*Progress)
	progress.GetTaskProgress("foo")
	progress.GetTaskProgress("bar")
	return progress
}

func TestProgressConformance(t *testing.T) {
	metadata_test.ConformanceMetadataTypeTest(t, newProgressforConformanceTest())
}
