package enum

import (
	"fmt"
	"testing"
)

func TestLogTypeMetadataIsFilled(t *testing.T) {
	for i := 0; i <= int(logTypeUnusedEnd); i++ {
		if _, ok := LogTypes[LogType(i)]; !ok {
			t.Errorf("LogTypeMetadata[%d] is not filled", i)
		}
	}
}

func TestLogTypeMetadataIsValid(t *testing.T) {
	for i := 0; i <= int(logTypeUnusedEnd); i++ {
		if logType, ok := LogTypes[LogType(i)]; ok {
			t.Run(fmt.Sprintf("%d-%s", i, logType.EnumKeyName), func(t *testing.T) {
				if logType.EnumKeyName == "" {
					t.Errorf("EnumKeyName in `%s(%d)` is empty", logType.Label, i)
				}
				if logType.LabelBackgroundColor == "" {
					t.Errorf("LabelBackgroundColor in `%s(%d)` is empty", logType.Label, i)
				}
				if logType.Label == "" {
					t.Errorf("Label in `%s(%d)` is empty", logType.Label, i)
				}
			})
		}
	}
}
