package parameters

import (
	"errors"
	"flag"
)

var Job *JobParameters = &JobParameters{}

type JobParameters struct {
	// JobMode
	// If this flag is set, KHI run as job mode and doesn't serve as a web server.
	JobMode *bool
	// InspectionType is the inspection type used in the job mode. Inspection types are defined in `inspection-type.go` in multiple directory. Check them for more details.
	InspectionType *string
	// InspectionFeatures is comma separated feature list to query.
	InspectionFeatures *string
	// InspectionValues is the JSON represented parameters.
	InspectionValues *string
	// ExportDestination is the destination file path of KHI file written after the query.
	ExportDestination *string
}

// PostProcess implements ParameterStore.
func (j *JobParameters) PostProcess() error {
	if *j.JobMode && (*j.InspectionType == "" || *j.InspectionFeatures == "" || *j.InspectionValues == "" || *j.ExportDestination == "") {
		return errors.New("`--job-inspection-type`, `--job-inspection-features`, `--job-inspection-values` and `--job-export-destination` are required when `--job-mode` is set")
	}
	return nil
}

// Prepare implements ParameterStore.
func (j *JobParameters) Prepare() error {
	j.JobMode = flag.Bool("job-mode", false, "If this flag is set, KHI run as job mode and doesn't serve as a web server.")
	j.InspectionType = flag.String("job-inspection-type", "", "(Job mode only)The inspection type used in the job mode. Inspection types are defined in `inspection-type.go` in multiple directory. Check them for more details.")
	j.InspectionFeatures = flag.String("job-inspection-features", "", "(Job mode only)Comma separated feature list to query.")
	j.InspectionValues = flag.String("job-inspection-values", "", "(Job mode only)The JSON represented parameters.")
	j.ExportDestination = flag.String("job-export-destination", "", "(Job mode only)The destination file path of KHI file written after the query.")
	return nil
}

var _ ParameterStore = (*JobParameters)(nil)
