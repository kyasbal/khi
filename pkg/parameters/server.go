package parameters

import "flag"

var Server *ServerParameters = &ServerParameters{}

type ServerParameters struct {
	// ViewerMode limits the KHI feature to query logs with the backend. When it is true, KHI is only serve the frontend to open KHI files.
	ViewerMode *bool
	// Port is the port number where KHI server listens.
	Port *int
	// Host is the host address where KHI server serves.
	Host *string
	// BasePath specifies the base address of API endpoints. This path always ends with `/`.
	BasePath *string
	// FrontendResourceBasePath is another base address only for frontend assets. If this value is not set, this uses the BasePath value by default.
	FrontendResourceBasePath *string
	// FrontendAssetFolder is the root folder of the assets used in frontend including index.html.
	FrontendAssetFolder *string
}

// PostProcess implements ParameterStore.
func (s *ServerParameters) PostProcess() error {
	ensureEndsWithSlash(s.BasePath)
	if *s.FrontendResourceBasePath == "" {
		*s.FrontendResourceBasePath = *s.BasePath
	}
	ensureEndsWithSlash(s.FrontendResourceBasePath)
	return nil
}

// Prepare implements ParameterStore.
func (s *ServerParameters) Prepare() error {
	s.ViewerMode = flag.Bool("viewer-mode", false, "Limits the KHI feature to query logs with the backend. When it is true, KHI is only serve the frontend to open KHI files.")
	s.Port = flag.Int("port", 8080, "The port number where KHI server listens.")
	s.Host = flag.String("host", "localhost", "The host address where KHI server serves.")
	s.BasePath = flag.String("base-path", "/", "The base address of API endpoints.")
	s.FrontendResourceBasePath = flag.String("frontend-resource-base-path", "", "Another base address only for frontend assets. If this value is not set, this uses `--base-path` value by default.")
	s.FrontendAssetFolder = flag.String("frontend-asset-folder", "./web", "The root folder of the assets used in frontend including index.html.")
	return nil
}

func ensureEndsWithSlash(v *string) {
	if *v != "" && (*v)[len(*v)-1] != '/' {
		*v += "/"
	}
}

var _ ParameterStore = (*ServerParameters)(nil)
