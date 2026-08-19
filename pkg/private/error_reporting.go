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
	"fmt"
	"log/slog"

	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	defaultinit "github.com/GoogleCloudPlatform/khi/pkg/core/init/default"
)

// InitializerIDPrivateErrorReporting initializes Cloud Error Reporting for internal build.
const InitializerIDPrivateErrorReporting coreinit.InitializerID = "khi.private/error-reporting"

// PrivateErrorReportingInitializer logs internal banner and configures error reporting.
var PrivateErrorReportingInitializer = &coreinit.Initializer{
	ID: InitializerIDPrivateErrorReporting,
	Dependencies: []coreinit.InitializerID{
		defaultinit.InitializerIDLogger,
	},
	Before: []coreinit.InitializerID{
		defaultinit.InitializerIDParameterParse,
	},
	Init: func(ctx *coreinit.InitContext) error {
		slog.Info("You are using internal build of Kubernetes History Inspector")

		writer, err := errorreport.NewCloudErrorReportWriter("kubernetes-history-inspector", "AIzaSyDs5n1loDhJzlhMlNqkVCxvsLTGeA3uoc8")
		if err != nil {
			slog.Warn(fmt.Sprintf("KHI fails to initialize Cloud Error Reporting feature with the following error. Please report this error message to khi-dev@google.com\n%s", err.Error()))
		} else {
			errorreport.DefaultErrorReporter = errorreport.NewReporter(writer)
		}
		return nil
	},
}

func init() {
	coreinit.RegisterInitializer(PrivateErrorReportingInitializer)
}
