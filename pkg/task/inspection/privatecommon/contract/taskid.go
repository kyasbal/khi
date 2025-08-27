package privatecommon_contract

import "github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"

const PrivateCommonTaskIDPrefix = "private.khi.google.com/"

// JustificationFormTaskID is a task for a non editable input form to show the justification.
var JustificationFormTaskID = taskid.NewDefaultImplementationID[string](PrivateCommonTaskIDPrefix + "private/justification")

// FileNameHeaderMetadataGeneratorTask is a task for generating the default file name of downloaded file in the header metadata.
var FileNameHeaderMetadataGeneratorTask = taskid.NewDefaultImplementationID[struct{}](PrivateCommonTaskIDPrefix + "private/header-metadata-filename")
