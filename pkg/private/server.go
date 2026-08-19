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
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	defaultinit "github.com/GoogleCloudPlatform/khi/pkg/core/init/default"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	privateserver "github.com/GoogleCloudPlatform/khi/pkg/private/server"
	"github.com/GoogleCloudPlatform/khi/pkg/private/server/index"
)

// InitializerIDPrivateServer configures private web server routes and tag generators.
const InitializerIDPrivateServer coreinit.InitializerID = "khi.private/server"

// PrivateServerInitializer registers index tag generators and mounts private API routes.
var PrivateServerInitializer = &coreinit.Initializer{
	ID: InitializerIDPrivateServer,
	Dependencies: []coreinit.InitializerID{
		defaultinit.InitializerIDGinServer,
		InitializerIDPrivateIAMToken,
	},
	Before: []coreinit.InitializerID{
		defaultinit.InitializerIDServerRunner,
	},
	Init: func(ctx *coreinit.InitContext) error {
		jobParams := coreinit.MustGet(ctx, defaultinit.JobParametersKey)
		if *jobParams.JobMode {
			return nil
		}
		index.RegisterAll()
		if privateparameters.Private.InspectionMode != nil && *privateparameters.Private.InspectionMode {
			router := coreinit.MustGet(ctx, defaultinit.GinRouterKey)
			injector, ok := coreinit.Get(ctx, IAMTokenInjectorKey)
			if ok && injector != nil {
				privateserver.ConfigureRoute(router, injector)
			}
		}
		return nil
	},
}

func init() {
	coreinit.RegisterInitializer(PrivateServerInitializer)
}
