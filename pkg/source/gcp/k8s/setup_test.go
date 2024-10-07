package k8s

import (
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil"
)

func TestMain(m *testing.M) {
	testutil.InitTestIO()
	code := m.Run()
	os.Exit(code)
}
