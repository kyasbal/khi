package inspectiondata

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestFileSystemResultRepository(t *testing.T) {
	testutil.InitTestIO()
	repo := NewFileSystemInspectionResultRepository("/tmp/test.json")
	t.Run("ReadInspectionResult can read inspection result written with WriteInspectionResult", func(t *testing.T) {
		testInspectionData := []byte{
			0x01, 0x02, 0x03, 0x04, 0x05,
		}

		writer, writeErr := repo.GetWriter()
		if writeErr != nil {
			t.Errorf("writeErr: want nil, got %s", writeErr)
		}
		writer.Write(testInspectionData)
		repo.Close()
		received, readErr := repo.GetReader()
		var readTarget = make([]byte, 5)
		_, err := received.Read(readTarget)
		if err != nil {
			t.Errorf("unexpected errir %s", err)
		}
		repo.Close()

		if readErr != nil {
			t.Errorf("readErr: want nil, got %s", readErr)
		}
		if diff := cmp.Diff(testInspectionData, readTarget); diff != "" {
			t.Errorf("+testInspectionData, -received,%s", diff)
		}
	})
}
