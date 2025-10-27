package private

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	"github.com/GoogleCloudPlatform/khi/pkg/lifecycle"
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	privatelifecycle "github.com/GoogleCloudPlatform/khi/pkg/private/lifecycle"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	privateserver "github.com/GoogleCloudPlatform/khi/pkg/private/server"
	"github.com/GoogleCloudPlatform/khi/pkg/private/server/index"
	"github.com/GoogleCloudPlatform/khi/pkg/server"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
)

func init() {
	coreinit.RegisterInitExtension(2000, &privateInitExtension{})

	lifecycle.Default.AddHandler(privatelifecycle.NewAnalyticsLifecycleHandler())
}

type privateInitExtension struct {
	iamTokenInjector *iamtoken.IAMTokenCallOptionInjectorOption
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

	if privateparameters.Private.InspectionMode != nil && *privateparameters.Private.InspectionMode {
		iamToken := *privateparameters.Private.IAMToken
		p.iamTokenInjector = iamtoken.New(iamToken)
	}
	return nil
}

// ConfigureInspectionTaskServer implements coreinit.InitExtension.
func (p *privateInitExtension) ConfigureInspectionTaskServer(taskServer *coreinspection.InspectionTaskServer) error {
	if privateparameters.Private.InspectionMode != nil && *privateparameters.Private.InspectionMode {
		taskServer.AddRunContextOption(coreinspection.RunContextOptionArrayElementFromValue[googlecloud.CallOptionInjectorOption](googlecloudcommon_contract.APICallOptionsInjectorContextKey, p.iamTokenInjector))
	}
	return nil
}

// ConfigureKHIWebServerFactory implements coreinit.InitExtension.
func (p *privateInitExtension) ConfigureKHIWebServerFactory(serverFactory *server.ServerFactory) error {
	index.RegisterAll()

	if privateparameters.Private.InspectionMode != nil && *privateparameters.Private.InspectionMode {
		basePath := strings.TrimSuffix(*parameters.Server.BasePath, "/")
		serverFactory.AddOptions(privateserver.NewPrivateServerOption(basePath, p.iamTokenInjector))
	}
	return nil
}

// BeforeTerminate implements coreinit.InitExtension.
func (p *privateInitExtension) BeforeTerminate() error {
	return nil
}

var _ coreinit.InitExtension = (*privateInitExtension)(nil)
