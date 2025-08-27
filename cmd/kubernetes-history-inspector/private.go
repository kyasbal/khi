// Copyright 2024 Google LLC
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

package main

// private.go
// This file is only included only for our internal build.

import (
	"fmt"
	"log/slog"

	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	"github.com/GoogleCloudPlatform/khi/pkg/lifecycle"
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	privateLifecycle "github.com/GoogleCloudPlatform/khi/pkg/private/lifecycle"
	privateParameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	privateIndex "github.com/GoogleCloudPlatform/khi/pkg/private/server/index"
)

func init() {
	slog.Info("You are using internal build of Kubernetes History Inspector")

	writer, err := errorreport.NewCloudErrorReportWriter("kubernetes-history-inspector", "AIzaSyDs5n1loDhJzlhMlNqkVCxvsLTGeA3uoc8") // 2nd argument is API key restricted only for error reporting. It's not sensitive value and this initialization happened before reading arguments thus this value is hard coded.
	if err != nil {
		slog.Warn(fmt.Sprintf("KHI fails to initialize Cloud Error Reporting feature with the following error. Please report this error message to khi-dev@google.com\n%s", err.Error()))
	} else {
		errorreport.DefaultErrorReporter = errorreport.NewReporter(writer)
	}

	lifecycle.Default.AddHandler(privateLifecycle.NewAnalyticsLifecycleHandler())
	lifecycle.Default.AddHandler(privateLifecycle.NewIAMTokenSetupLifecycleHandler())
	lifecycle.Default.AddHandler(privateLifecycle.NewErrorReportLifecycleHandler())

	parameters.AddStore(privateParameters.Private)

	privateIndex.RegisterAll()
}
