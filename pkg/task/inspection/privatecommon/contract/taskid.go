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

package privatecommon_contract

import (
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
)

const PrivateCommonTaskIDPrefix = "private.khi.google.com/"

// JustificationFormTaskID is a task for a non editable input form to show the justification.
var JustificationFormTaskID = taskid.NewDefaultImplementationID[string](PrivateCommonTaskIDPrefix + "private/justification")

// FileNameHeaderMetadataGeneratorTask is a task for generating the default file name of downloaded file in the header metadata.
var FileNameHeaderMetadataGeneratorTask = taskid.NewDefaultImplementationID[struct{}](PrivateCommonTaskIDPrefix + "private/header-metadata-filename")

var APIClientCallOptionsInjectorTaskOverrideID = taskid.NewImplementationID(googlecloudcommon_contract.APIClientCallOptionsInjectorTaskID.Ref(), "private")
