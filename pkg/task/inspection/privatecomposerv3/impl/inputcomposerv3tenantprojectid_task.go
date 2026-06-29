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

package privatecomposerv3_impl

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/formtask"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
	privatecomposerv3_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposerv3/contract"
)

var InputComposerV3TenantProjectIdTask = formtask.NewTextFormTaskBuilder(
	privatecomposerv3_contract.InputComposerV3TenantProjectIdTaskID,
	googlecloudcommon_contract.PriorityForResourceIdentifierGroup+4900,
	"Managed Airflow 3 Tenant Project ID",
).
	WithDescription("Type the tenant project ID for the Managed Airflow 3 environment. You can find the tenant ID from the tenant project section in Google Admin. The project ID must end with '-tp'.").
	WithValidatingTiming(inspectionmetadata.Blur).
	WithDependencies([]taskid.UntypedTaskReference{
		googlecloudcommon_contract.InputProjectIdTaskID.Ref(),
		privatecommon_contract.JustificationFormTaskID.Ref(),
	}).
	WithDefaultValueFunc(func(ctx context.Context, previousValues []string) (string, error) {
		if len(previousValues) > 0 {
			return previousValues[0], nil
		}
		return "", nil
	}).
	WithValidator(func(ctx context.Context, value string) (string, error) {
		justification := coretask.GetTaskResult(ctx, privatecommon_contract.JustificationFormTaskID.Ref())
		composerProjectID := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputProjectIdTaskID.Ref())
		iamTokenInjector, err := khictx.GetValue(ctx, privatecommon_contract.APIClientIAMTokenInjectorOptionContextKey)
		if err != nil {
			return "IAMToken injector isn't set. Managed Airflow 3 parsers won't work unless you open KHI from Google Admin", nil
		}
		if value == "" || !strings.HasSuffix(value, "-tp") {
			link, linkErr := api.ToGoogleAdminLink(googlecloud.Project(composerProjectID), justification)
			if linkErr != nil {
				return linkErr.Error(), nil
			}
			if value == "" {
				return fmt.Sprintf("Tenant project ID is not provided. Please visit %s and open the project card and tenant project tool to find the project ID.", link), nil
			}
			return fmt.Sprintf("the tenant project id must end with '-tp', but got '%s'. Please visit %s and open the project card and tenant project tool to find the right project ID.", value, link), nil
		}
		if err == nil && iamTokenInjector != nil {
			if !iamTokenInjector.HasTokenFor(googlecloud.Project(value)) {
				link, linkErr := api.ToGoogleAdminLink(googlecloud.Project(value), justification)
				if linkErr != nil {
					return linkErr.Error(), nil
				}
				return fmt.Sprintf("IAMToken for project '%s' is not set. Please visit %s and open KHI from the project inspection card on the tenant project to set it", value, link), nil
			}
		}
		return "", nil
	}).
	Build()
