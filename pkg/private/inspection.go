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

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	defaultinit "github.com/GoogleCloudPlatform/khi/pkg/core/init/default"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
)

// InitializerIDPrivateInspection configures InspectionTaskServer with private run context options.
const InitializerIDPrivateInspection coreinit.InitializerID = "khi.private/inspection"

// PrivateInspectionInitializer configures InspectionTaskServer with private injectors.
var PrivateInspectionInitializer = &coreinit.Initializer{
	ID: InitializerIDPrivateInspection,
	Dependencies: []coreinit.InitializerID{
		defaultinit.InitializerIDInspectionTaskServer,
		InitializerIDPrivateIAMToken,
	},
	Before: []coreinit.InitializerID{
		defaultinit.InitializerIDServerRunner,
		defaultinit.InitializerIDJobRunner,
	},
	Init: func(ctx *coreinit.InitContext) error {
		if privateparameters.Private.InspectionMode != nil && *privateparameters.Private.InspectionMode {
			taskServer := coreinit.MustGet(ctx, defaultinit.InspectionTaskServerKey)
			injector, ok := coreinit.Get(ctx, IAMTokenInjectorKey)
			if ok && injector != nil {
				taskServer.AddRunContextOption(coreinspection.RunContextOptionArrayElementFromValue[googlecloud.CallOptionInjectorOption](googlecloudcommon_contract.APICallOptionsInjectorContextKey, injector))
				taskServer.AddRunContextOption(func(taskCtx context.Context, mode inspectioncore_contract.InspectionTaskModeType) (context.Context, error) {
					return khictx.WithValue(taskCtx, privatecommon_contract.APIClientIAMTokenInjectorOptionContextKey, injector), nil
				})
			}
		}
		return nil
	},
}

func init() {
	coreinit.RegisterInitializer(PrivateInspectionInitializer)
}
