// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package private

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
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
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
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
		fixedProjectID := *parameters.Auth.FixedProjectID
		if fixedProjectID == "" {
			panic("fixed project ID is not set. Default IAM token must be used with fixed project ID. b/485682103")
		}
		p.iamTokenInjector = iamtoken.NewInjector()
		p.iamTokenInjector.SetTokenFor(googlecloud.Project(fixedProjectID), iamToken)
	}
	return nil
}

// ConfigureInspectionTaskServer implements coreinit.InitExtension.
func (p *privateInitExtension) ConfigureInspectionTaskServer(taskServer *coreinspection.InspectionTaskServer) error {
	if privateparameters.Private.InspectionMode != nil && *privateparameters.Private.InspectionMode {
		taskServer.AddRunContextOption(coreinspection.RunContextOptionArrayElementFromValue[googlecloud.CallOptionInjectorOption](googlecloudcommon_contract.APICallOptionsInjectorContextKey, p.iamTokenInjector))
		taskServer.AddRunContextOption(func(ctx context.Context, mode inspectioncore_contract.InspectionTaskModeType) (context.Context, error) {
			return khictx.WithValue(ctx, privatecommon_contract.APIClientIAMTokenInjectorOptionContextKey, p.iamTokenInjector), nil
		})
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
