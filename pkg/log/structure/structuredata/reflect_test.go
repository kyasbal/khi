package structuredata

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil"
	corev1 "k8s.io/api/core/v1"
)

func TestReadPodManifest(t *testing.T) {
	testutil.InitTestIO()
	podYaml := testutil.MustReadText("test/k8s/sample_pod.yaml")
	sd, err := DataFromYaml(podYaml)
	if err != nil {
		t.Fatal(err)
	}
	var pod corev1.Pod
	err = ReadReflectK8sManifest(sd, &pod)
	if err != nil {
		t.Fatal(err)
	}
	if pod.UID != "7899f560-3d56-4831-a381-2691c28ea3e5" {
		t.Errorf("parsed pod UID is not matching with the expected value")
	}
}
