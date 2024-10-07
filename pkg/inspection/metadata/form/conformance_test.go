package form

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	metadata_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/metadata"
)

func newFormFieldsForConformanceTest() metadata.Metadata {
	forms := (&FormFieldSetMetadataFactory{}).Instanciate().(*FormFieldSet)
	forms.SetField(&FormField{})
	return forms
}

func TestProgressConformance(t *testing.T) {
	metadata_test.ConformanceMetadataTypeTest(t, newFormFieldsForConformanceTest())
}
