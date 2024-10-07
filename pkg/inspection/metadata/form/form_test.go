package form

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func fieldWithIdAndPriorityForTest(id string, priority int) *FormField {
	return &FormField{
		Id:       id,
		Priority: priority,
	}
}

func TestFormFieldSetShouldSortOnAddingNewField(t *testing.T) {
	fsActual := (&FormFieldSetMetadataFactory{}).Instanciate().(*FormFieldSet)
	fsActual.SetField(fieldWithIdAndPriorityForTest("foo", 1))
	fsActual.SetField(fieldWithIdAndPriorityForTest("bar", 3))
	fsActual.SetField(fieldWithIdAndPriorityForTest("qux", 2))

	fsExpected := &FormFieldSet{
		fields: []*FormField{
			fieldWithIdAndPriorityForTest("bar", 3),
			fieldWithIdAndPriorityForTest("qux", 2),
			fieldWithIdAndPriorityForTest("foo", 1),
		},
	}

	if diff := cmp.Diff(fsActual, fsExpected, cmp.AllowUnexported(FormFieldSet{})); diff != "" {
		t.Errorf("FieldSet has fields in unexpected shape\n%v", diff)
	}
}
