package private_taskid

import (
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	gcp_task "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task"
)

// JustificationFormTaskID is a task for a non editable input form to show the justification.
var JustificationFormTaskID = taskid.NewDefaultImplementationID[string](gcp_task.GCPPrefix + "private/justification")

// FileNameHeaderMetadataGeneratorTask is a task for generating the default file name of downloaded file in the header metadata.
var FileNameHeaderMetadataGeneratorTask = taskid.NewDefaultImplementationID[struct{}](gcp_task.GCPPrefix + "private/header-metadata-filename")
