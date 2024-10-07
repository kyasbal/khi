package testlog

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTestLogWith(t *testing.T) {
	tl := New(BaseYaml(`foo: bar`))
	expectedTl1 := `foo: bar1
`
	expectedTl2 := `foo: bar2
`
	// With returns a new instance of TestLog and each instances are independent
	tl1 := tl.With(StringField("foo", "bar1"))
	tl2 := tl1.With(StringField("foo", "bar2"))

	tl1Yaml := tl1.MustBuildYamlString()
	tl2Yaml := tl2.MustBuildYamlString()
	if diff := cmp.Diff(tl1Yaml, expectedTl1); diff != "" {
		t.Errorf("Yaml string mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(tl2Yaml, expectedTl2); diff != "" {
		t.Errorf("Yaml string mismatch (-want +got):\n%s", diff)
	}
}
