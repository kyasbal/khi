package ioconfig

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var IOConfigTaskName = task.KHISystemPrefix + "inspection/ioconfig"

type IOConfig struct {
	// The project root folder
	ApplicationRoot string
	// The folder to save khi files
	DataDestination string
	// TemporaryFolder working folder
	TemporaryFolder string
}

var ProductionIOConfig = task.NewCachedProcessor(IOConfigTaskName, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
	dataDestinationFolder := "./data"
	if parameters.Common.DataDestinationFolder != nil {
		dataDestinationFolder = *parameters.Common.DataDestinationFolder
	}
	temporaryFolder := "/tmp"
	if parameters.Common.TemporaryFolder != nil {
		temporaryFolder = *parameters.Common.TemporaryFolder
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(dataDestinationFolder) {
		dataDestinationFolder = filepath.Join(dir, dataDestinationFolder)
	}
	return &IOConfig{
		ApplicationRoot: dir,
		DataDestination: dataDestinationFolder,
		TemporaryFolder: temporaryFolder,
	}, nil
})

var TestIOConfig = task.NewCachedProcessor(IOConfigTaskName, []string{}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".root")); err == nil {
			break
		}
		pathsSegments := strings.Split(dir, "/")
		dir = "/" + filepath.Join(pathsSegments[:len(pathsSegments)-1]...)
	}
	return &IOConfig{
		ApplicationRoot: dir + "/",
		DataDestination: "/tmp/",
		TemporaryFolder: "/tmp/",
	}, nil
})

func GetIOConfigFromTaskVariable(v *task.VariableSet) (*IOConfig, error) {
	return task.GetTypedVariableFromTaskVariable[*IOConfig](v, IOConfigTaskName, nil)
}
