package ioconfig

import (
	"os"
	"path/filepath"
	"testing"

	env_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/inspection/env"
	task_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/task"
)

func TestTestIOConfigCanFindTheRoot(t *testing.T) {
	vs, err := task_test.RunTaskGraph(TestIOConfig, 0, map[string]any{})
	if err != nil {
		t.Errorf("unxepected error %v", err)
	}
	ioConfig, err := GetIOConfigFromTaskVariable(vs)
	if err != nil {
		t.Errorf("unxepected error %v", err)
	}
	stat, err := os.Stat(ioConfig.ApplicationRoot)
	if err != nil {
		t.Errorf("unxepected error %v", err)
	}
	if !stat.IsDir() {
		t.Errorf("the result application root must be a directory")
	}
}

func TestProductionIOConfigConvertPathToAbs(t *testing.T) {
	vs, err := task_test.RunTaskGraph(ProductionIOConfig, 0, map[string]any{}, env_test.MockedEnvironmentVariableProducer(EnvDataDestinationTask, "./data"), env_test.MockedEnvironmentVariableProducer(EnvTemporaryFolderTask, "/tmp"))
	if err != nil {
		t.Errorf("unxepected error %v", err)
	}
	ioConfig, err := GetIOConfigFromTaskVariable(vs)
	if err != nil {
		t.Errorf("unxepected error %v", err)
	}
	if !filepath.IsAbs(ioConfig.ApplicationRoot) {
		t.Errorf("the given application folder must be abs path")
	}
	if !filepath.IsAbs(ioConfig.DataDestination) {
		t.Errorf("the given data destination folder must be abs path")
	}
	stat, err := os.Stat(ioConfig.ApplicationRoot)
	if err != nil {
		t.Errorf("unxepected error %v", err)
	}
	if !stat.IsDir() {
		t.Errorf("the result application root must be a directory")
	}
}
