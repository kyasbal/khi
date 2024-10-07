package enum

import (
	"fmt"
	"testing"
)

func TestRevisionVerbIsFilled(t *testing.T) {
	for i := 0; i <= int(revisionVerbUnusedEnd); i++ {
		if _, ok := RevisionVerbs[RevisionVerb(i)]; !ok {
			t.Errorf("RevisionVerb[%d] is not filled", i)
		}
	}
}

func TestRevisionVerbIsValid(t *testing.T) {
	for i := 0; i <= int(revisionVerbUnusedEnd); i++ {
		if verb, ok := RevisionVerbs[RevisionVerb(i)]; ok {
			t.Run(fmt.Sprintf("%d-%s", i, verb.EnumKeyName), func(t *testing.T) {
				if verb.Label == "" {
					t.Errorf("Label in %s(%d) is empty", verb.EnumKeyName, i)
				}
				if verb.LabelBackgroundColor == "" {
					t.Errorf("LabelBackgroundColor in %s(%d) is empty", verb.EnumKeyName, i)
				}
				if verb.CSSSelector == "" {
					t.Errorf("CSSSelector in %s(%d) is empty", verb.EnumKeyName, i)
				}
			})
		}
	}
}
