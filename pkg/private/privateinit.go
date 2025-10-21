package private

import (
	"fmt"
	"log/slog"

	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	"github.com/GoogleCloudPlatform/khi/pkg/lifecycle"
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	privatelifecycle "github.com/GoogleCloudPlatform/khi/pkg/private/lifecycle"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	"github.com/GoogleCloudPlatform/khi/pkg/private/server/index"
	"github.com/GoogleCloudPlatform/khi/pkg/server"
)

func init() {
	coreinit.RegisterInitExtension(2000, &privateInitExtension{})

	lifecycle.Default.AddHandler(privatelifecycle.NewAnalyticsLifecycleHandler())
}

type privateInitExtension struct {
}

// BeforeAll implements coreinit.InitExtension.
func (p *privateInitExtension) BeforeAll() error {
	slog.Info("You are using internal build of Kubernetes History Inspector")

	writer, err := errorreport.NewCloudErrorReportWriter("kubernetes-history-inspector", "AIzaSyDs5n1loDhJzlhMlNqkVCxvsLTGeA3uoc8") // 2nd argument is API key restricted only for error reporting. It's not sensitive value and this initialization happened before reading arguments thus this value is hard coded.
	if err != nil {
		slog.Warn(fmt.Sprintf("KHI fails to initialize Cloud Error Reporting feature with the following error. Please report this error message to khi-dev@google.com\n%s", err.Error()))
	} else {
		errorreport.DefaultErrorReporter = errorreport.NewReporter(writer)
	}
	return nil
}

// ConfigureParameterStore implements coreinit.InitExtension.
func (p *privateInitExtension) ConfigureParameterStore() error {
	parameters.AddStore(privateparameters.Private)
	return nil
}

// AfterParsingParameters implements coreinit.InitExtension.
func (p *privateInitExtension) AfterParsingParameters() error {
	if privateparameters.Private.GALabels != nil {
		metadata := privateparameters.Private.GetMapOfGALabels()
		for key, value := range metadata {
			errorreport.DefaultErrorReporter.SetMetadataEntry(key, value)
		}
	}
	return nil
}

// ConfigureInspectionTaskServer implements coreinit.InitExtension.
func (p *privateInitExtension) ConfigureInspectionTaskServer(taskServer *coreinspection.InspectionTaskServer) error {
	return nil
}

// ConfigureKHIWebServerFactory implements coreinit.InitExtension.
func (p *privateInitExtension) ConfigureKHIWebServerFactory(serverFactory *server.ServerFactory) error {
	index.RegisterAll()
	return nil
}

// BeforeTerminate implements coreinit.InitExtension.
func (p *privateInitExtension) BeforeTerminate() error {
	return nil
}

var _ coreinit.InitExtension = (*privateInitExtension)(nil)
