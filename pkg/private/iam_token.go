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
	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	defaultinit "github.com/GoogleCloudPlatform/khi/pkg/core/init/default"
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
)

var (
	// IAMTokenInjectorKey stores the IAMTokenCallOptionInjectorOption instance.
	IAMTokenInjectorKey = typedmap.NewTypedKey[*iamtoken.IAMTokenCallOptionInjectorOption]("khi.google.com/init/private/iam-token-injector")
)

// InitializerIDPrivateIAMToken initializes IAM token injector from parsed parameters.
const InitializerIDPrivateIAMToken coreinit.InitializerID = "khi.private/iam-token"

// PrivateIAMTokenInitializer initializes the IAM token injector from parsed flags.
var PrivateIAMTokenInitializer = &coreinit.Initializer{
	ID: InitializerIDPrivateIAMToken,
	Dependencies: []coreinit.InitializerID{
		defaultinit.InitializerIDParameterParse,
	},
	Init: func(ctx *coreinit.InitContext) error {
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
			injector := iamtoken.NewInjector()
			injector.SetTokenFor(googlecloud.Project(fixedProjectID), iamToken)
			coreinit.Set(ctx, IAMTokenInjectorKey, injector)
		}
		return nil
	},
}

func init() {
	coreinit.RegisterInitializer(PrivateIAMTokenInitializer)
}
