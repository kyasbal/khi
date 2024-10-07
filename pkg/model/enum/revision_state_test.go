package enum

import (
	"fmt"
	"testing"
)

func TestRevisionStatesIsFilled(t *testing.T) {
	for i := 0; i <= int(revisionStateUnusedEnd); i++ {
		if _, ok := RevisionStates[RevisionState(i)]; !ok {
			t.Errorf("RevisionState[%d] is not filled", i)
		}
	}
}

func TestRevisionStatesIsValid(t *testing.T) {
	for i := 0; i <= int(revisionStateUnusedEnd); i++ {
		if state, ok := RevisionStates[RevisionState(i)]; ok {
			t.Run(fmt.Sprintf("%d-%s", i, state.EnumKeyName), func(t *testing.T) {
				if state.EnumKeyName == "" {
					t.Errorf("EnumKeyName in `%s(%d)` is empty", state.EnumKeyName, i)
				}
				if state.BackgroundColor == "" {
					t.Errorf("LabelBackgroundColor in `%s(%d)` is empty", state.EnumKeyName, i)
				}
				if state.CSSSelector == "" {
					t.Errorf("CSSSelector in `%s(%d)` is empty", state.EnumKeyName, i)
				}
				if state.Label == "" {
					t.Errorf("Label in `%s(%d)` is empty", state.EnumKeyName, i)
				}
			})
		}
	}
}
