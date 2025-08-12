package private_taskid

import (
	gcp_task "github.com/GoogleCloudPlatform/khi/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/khi/pkg/task/core/contract/taskid"
)

// JustificationFormTaskID is a task for a non editable input form to show the justification.
var JustificationFormTaskID = taskid.NewDefaultImplementationID[string](gcp_task.GCPPrefix + "private/justification")

// FileNameHeaderMetadataGeneratorTask is a task for generating the default file name of downloaded file in the header metadata.
var FileNameHeaderMetadataGeneratorTask = taskid.NewDefaultImplementationID[struct{}](gcp_task.GCPPrefix + "private/header-metadata-filename")
