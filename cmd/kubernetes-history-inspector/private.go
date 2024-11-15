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
	"log/slog"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/lifecycle"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
	privateLifecycle "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/lifecycle"
	privateParameters "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/parameters"
	privateIndex "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/server/index"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/task"
)

func init() {
	slog.Info("You are using internal build of Kubernetes History Inspector")

	lifecycle.Default.AddHandler(privateLifecycle.NewAnalyticsLifecycleHandler())
	lifecycle.Default.AddHandler(privateLifecycle.NewIAMTokenSetupLifecycleHandler())

	parameters.AddStore(privateParameters.Private)

	taskSetRegistrer = append(taskSetRegistrer, task.PrepareInspectionServer)

	privateIndex.RegisterAll()
}
