package privatecommon_impl

import (
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
)

func Register(inspectionServer coreinspection.InspectionTaskRegistry) error {
	return coretask.RegisterTasks(inspectionServer,
		FilenameHeaderMetadataGeneratorTask,
		JustificationFormTask,
	)
}
