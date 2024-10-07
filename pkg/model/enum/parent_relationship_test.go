package enum

import (
	"fmt"
	"strings"
	"testing"
)

func TestParentRelationshipMetadataIsFilled(t *testing.T) {
	for i := 0; i <= int(relationshipUnusedEnd); i++ {
		if _, ok := ParentRelationships[ParentRelationShip(i)]; !ok {
			t.Errorf("ParentRelationshipMetadata[%d] is not filled", i)
		}
	}
}

func TestParentRelationshipMetadataIsValid(t *testing.T) {
	for i := 0; i <= int(relationshipUnusedEnd); i++ {
		if relationship, ok := ParentRelationships[ParentRelationShip(i)]; ok {
			t.Run(fmt.Sprintf("%d-%s", i, relationship.EnumKeyName), func(t *testing.T) {
				if relationship.EnumKeyName == "" {
					t.Errorf("EnumKeyName in `%s(%d)` is empty", relationship.Label, i)
				}
				if relationship.Visible {
					if relationship.LabelColor == "" {
						t.Errorf("LabelColor in `%s(%d)` is empty", relationship.Label, i)
					}
					if relationship.LabelBackgroundColor == "" {
						t.Errorf("LabelBackgroundColor in `%s(%d)` is empty", relationship.Label, i)
					}
					if relationship.Hint == "" {
						t.Errorf("Hint in `%s(%d)` is empty", relationship.Label, i)
					}
					if strings.Contains(relationship.Label, " ") {
						t.Errorf("Label in `%s(%d)` contains space(label must be valid as css class)", relationship.Label, i)
					}
				}
			})
		}
	}
}
