package common

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/ioconfig"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
)

func PrepareInspectionServer(rootServer *inspection.InspectionTaskServer) error {
	err := rootServer.AddTaskDefinition(task.InspectionTimeProducer)
	if err != nil {
		return err
	}

	err = rootServer.AddTaskDefinition(ioconfig.ProductionIOConfig)
	if err != nil {
		return err
	}

	err = rootServer.AddTaskDefinition(task.BuilderGeneratorTask)
	if err != nil {
		return err
	}

	err = rootServer.AddTaskDefinition(task.ReaderFactoryGeneratorTask)
	if err != nil {
		return err
	}

	return nil
}
