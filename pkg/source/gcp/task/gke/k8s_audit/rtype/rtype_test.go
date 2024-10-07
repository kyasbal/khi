package rtype

import (
	"fmt"
	"testing"
)

func TestTypesAreFilled(t *testing.T) {
	for i := 1; i <= RTypeUnusedEnd; i++ {
		t.Run(fmt.Sprintf("check-%d-filled", i), func(t *testing.T) {
			for _, value := range Types {
				if value == i {
					return
				}
			}
			t.Errorf("type(%d) is not included in the Types", i)
		})
	}
}
