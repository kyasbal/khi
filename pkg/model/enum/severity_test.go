package enum

import (
	"fmt"
	"testing"
)

func TestSeverityMetadataIsFilled(t *testing.T) {
	for i := 0; i <= int(severityUnusedEnd); i++ {
		if _, ok := Severities[Severity(i)]; !ok {
			t.Errorf("SeverityMetadata[%d] is not filled", i)
		}
	}
}

func TestSeverityMetadataIsValid(t *testing.T) {
	for i := 0; i <= int(severityUnusedEnd); i++ {
		if severity, ok := Severities[Severity(i)]; ok {
			t.Run(fmt.Sprintf("%d-%s", i, Severities[Severity(i)].Label), func(t *testing.T) {
				if severity.EnumKeyName == "" {
					t.Errorf("EnumKeyName in `%s(%d)` is empty", severity.Label, i)
				}
				if severity.LabelColor == "" {
					t.Errorf("LabelColor in `%s(%d)` is empty", severity.Label, i)
				}
				if severity.BackgroundColor == "" {
					t.Errorf("LabelBackgroundColor in `%s(%d)` is empty", severity.Label, i)
				}
			})
		}
	}
}
