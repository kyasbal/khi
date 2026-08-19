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
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
)

// InitializerIDPrivateParameters registers private parameter stores.
const InitializerIDPrivateParameters coreinit.InitializerID = "khi.private/parameters"

// PrivateParametersInitializer registers private parameters into the global parameters store.
var PrivateParametersInitializer = &coreinit.Initializer{
	ID: InitializerIDPrivateParameters,
	Before: []coreinit.InitializerID{
		defaultinit.InitializerIDParameterParse,
	},
	Init: func(ctx *coreinit.InitContext) error {
		parameters.AddStore(privateparameters.Private)
		return nil
	},
}

func init() {
	coreinit.RegisterInitializer(PrivateParametersInitializer)
}
