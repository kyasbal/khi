// Copyright 2025 Google LLC
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

package parser

import (
	"context"

	"github.com/GoogleCloudPlatform/khi/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/khi/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/khi/pkg/source/oss/constant"
	"github.com/GoogleCloudPlatform/khi/pkg/source/oss/form"
	"github.com/GoogleCloudPlatform/khi/pkg/task"
)

// TODO: This file is a place holder for adding OSS type KHI parser.
var OSSPlaceHolderParser = inspection_task.NewInspectionProcessor(
	constant.OSSTaskPrefix+"placeholder",
	[]string{
		form.AuditLogFilesForm.ID().String(),
		form.TestTextForm.ID().String(),
	},
	func(ctx context.Context, taskMode int, v *task.VariableSet, progress *progress.TaskProgress) (any, error) {
		return nil, nil
	},
	inspection_task.FeatureTaskLabel("placeholder", "test", false),
	inspection_task.InspectionTypeLabel(constant.OSSInspectionTypeID),
)
