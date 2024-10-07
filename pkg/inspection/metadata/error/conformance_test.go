package error

import (
	"testing"

	metadata_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/metadata"
)

func TestProgressConformance(t *testing.T) {
	metadata_test.ConformanceMetadataTypeTest(t, &ErrorMessageSet{
		[]*ErrorMessage{
			NewNotFoundErrorMessage("foo-bar"),
		},
	})
}
