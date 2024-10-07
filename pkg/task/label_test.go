package task

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewLabelSet(t *testing.T) {
	labelSet := NewLabelSet()

	if diff := cmp.Diff(labelSet, &LabelSet{
		rawLabels: map[string]interface{}{},
	}, cmp.AllowUnexported(LabelSet{})); diff != "" {
		t.Errorf("Generated KHITaskLabelSet is mismatched, (-actual,+expected);\n%s", diff)
	}
}
