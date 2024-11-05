package parameters

import "flag"

var Common *CommonParameters = &CommonParameters{}

type CommonParameters struct {
	// DataDestinationFolder is the folder path where the final khi file to be stored for serving.
	DataDestinationFolder *string
	// TemporaryFolder is the folder path where be used as a working directory to generate the final khi file.
	TemporaryFolder *string
}

// PostProcess implements ParameterStore.
func (c *CommonParameters) PostProcess() error {
	return nil
}

// Prepare implements ParameterStore.
func (c *CommonParameters) Prepare() error {
	c.DataDestinationFolder = flag.String("data-destination-folder", "./data", "The folder path where the final khi file to be stored for serving.")
	c.TemporaryFolder = flag.String("temporary-folder", "/tmp", "The folder path where be used as a working directory to generate the final khi file.")
	return nil
}

var _ ParameterStore = (*CommonParameters)(nil)
