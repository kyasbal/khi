package ioconfig

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/env"
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

var EnvDataDestinationTaskId = task.KHISystemPrefix + "inspection/data-location"

var EnvDataDestinationTask = env.EnvironmentVariableProducer(EnvDataDestinationTaskId, "DATA_DESTINATION_FOLDER", "./data")

var EnvTemporaryFolderTaskId = task.KHISystemPrefix + "inspection/tmp-location"

var EnvTemporaryFolderTask = env.EnvironmentVariableProducer(EnvTemporaryFolderTaskId, "TMPORARY_FOLDER", "/tmp/")

var ProductionIOConfig = task.NewCachedProcessor(IOConfigTaskName, []string{EnvDataDestinationTaskId, EnvTemporaryFolderTaskId}, func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
	dataFolder, err := env.GetEnvironmentVariableFromTaskVariables(ctx, EnvDataDestinationTaskId, v)
	if err != nil {
		return nil, err
	}
	tmpFolder, err := env.GetEnvironmentVariableFromTaskVariables(ctx, EnvTemporaryFolderTaskId, v)
	if err != nil {
		return nil, err
	}
	dataFolderPath := dataFolder.Value
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(dataFolderPath) {
		dataFolderPath = filepath.Join(dir, dataFolderPath)
	}
	slog.InfoContext(ctx, fmt.Sprintf("Application root: %s , Data destination path: %s, Temporary folder: %s", dir, dataFolderPath, tmpFolder.Value))
	return &IOConfig{
		ApplicationRoot: dir,
		DataDestination: dataFolderPath,
		TemporaryFolder: tmpFolder.Value,
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
