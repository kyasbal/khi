package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var ForceUpdateGolden = false

// VerifyWithGolden test if the given body is matching with the golden file saved for each tests under /test/golden/
func VerifyWithGolden(t *testing.T, verificationTargetName string, body string) {
	t.Helper()
	_, updateGolden := os.LookupEnv("UPDATE_GOLDEN")
	goldenName := t.Name() + "-" + verificationTargetName
	goldenPath := fmt.Sprintf("test/golden/%s", goldenName)
	if updateGolden || ForceUpdateGolden {
		file, err := os.OpenFile(goldenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Golden file updated: %s", goldenPath)
	} else {
		golden := MustReadText(goldenPath)
		if diff := cmp.Diff(golden, body); diff != "" {
			t.Errorf("input is not matching with the golden (-want,+got):\n%s", diff)
		}
	}
}
